package graphBetaPlayLostModeSoundManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *PlayLostModeSoundManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data PlayLostModeSoundManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	if len(data.ManagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ManagedDevices {
			id := device.DeviceID.ValueString()
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"Sound will only be played once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	if len(data.ComanagedDevices) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, device := range data.ComanagedDevices {
			id := device.DeviceID.ValueString()
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_devices"),
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"Sound will only be played once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	for _, managedDevice := range data.ManagedDevices {
		managedID := managedDevice.DeviceID.ValueString()
		for _, comanagedDevice := range data.ComanagedDevices {
			comanagedID := comanagedDevice.DeviceID.ValueString()
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_devices"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
						"A device should only be in one list. Sound play will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating play lost mode sound action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string
	var unsupervisedManagedDevices []string
	var unsupervisedComanagedDevices []string
	var notInLostModeManagedDevices []string
	var notInLostModeComanagedDevices []string

	for _, managedDevice := range data.ManagedDevices {
		deviceID := managedDevice.DeviceID.ValueString()
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentManagedDevices = append(nonExistentManagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("managed_devices"),
					"Error Validating Managed Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility - lost mode sound only works on iOS/iPadOS devices
			if device.GetDeviceType() != nil {
				deviceType := device.GetDeviceType().String()
				if deviceType != "iPad" && deviceType != "iPhone" && deviceType != "iPod" {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (Type: %s)", deviceID, deviceType))
					continue
				}
			} else {
				unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (Unknown device type)", deviceID))
				continue
			}

			// Check if device is supervised (required for lost mode)
			if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
				unsupervisedManagedDevices = append(unsupervisedManagedDevices, deviceID)
			}

			// Check if device is in lost mode
			if device.GetLostModeState() != nil {
				lostModeState := device.GetLostModeState().String()
				if lostModeState == "disabled" {
					notInLostModeManagedDevices = append(notInLostModeManagedDevices, fmt.Sprintf("%s (state: %s)", deviceID, lostModeState))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	for _, comanagedDevice := range data.ComanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()
		device, err := a.client.
			DeviceManagement().
			ComanagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentComanagedDevices = append(nonExistentComanagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("comanaged_devices"),
					"Error Validating Co-Managed Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility - lost mode sound only works on iOS/iPadOS devices
			if device.GetDeviceType() != nil {
				deviceType := device.GetDeviceType().String()
				if deviceType != "iPad" && deviceType != "iPhone" && deviceType != "iPod" {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (Type: %s)", deviceID, deviceType))
					continue
				}
			} else {
				unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (Unknown device type)", deviceID))
				continue
			}

			// Check if device is supervised (required for lost mode)
			if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
				unsupervisedComanagedDevices = append(unsupervisedComanagedDevices, deviceID)
			}

			// Check if device is in lost mode
			if device.GetLostModeState() != nil {
				lostModeState := device.GetLostModeState().String()
				if lostModeState == "disabled" {
					notInLostModeComanagedDevices = append(notInLostModeComanagedDevices, fmt.Sprintf("%s (state: %s)", deviceID, lostModeState))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully", deviceID))
		}
	}

	if len(nonExistentManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Non-Existent Managed Devices",
			fmt.Sprintf("The following managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentManagedDevices, ", ")),
		)
	}

	if len(nonExistentComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Non-Existent Co-Managed Devices",
			fmt.Sprintf("The following co-managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing co-managed devices.",
				strings.Join(nonExistentComanagedDevices, ", ")),
		)
	}

	if len(unsupportedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Unsupported Managed Devices for Lost Mode Sound",
			fmt.Sprintf("Lost mode sound is only supported on iOS and iPadOS devices. The following managed devices are not supported: %s. "+
				"Please remove non-iOS/iPadOS devices from the configuration.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Unsupported Co-Managed Devices for Lost Mode Sound",
			fmt.Sprintf("Lost mode sound is only supported on iOS and iPadOS devices. The following co-managed devices are not supported: %s. "+
				"Please remove non-iOS/iPadOS devices from the configuration.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	if len(unsupervisedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Unsupervised Managed Devices",
			fmt.Sprintf("Lost mode sound requires supervised iOS/iPadOS devices. The following managed devices are not supervised: %s. "+
				"Please ensure devices are enrolled via DEP/ABM or manually supervised.",
				strings.Join(unsupervisedManagedDevices, ", ")),
		)
	}

	if len(unsupervisedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Unsupervised Co-Managed Devices",
			fmt.Sprintf("Lost mode sound requires supervised iOS/iPadOS devices. The following co-managed devices are not supervised: %s. "+
				"Please ensure devices are enrolled via DEP/ABM or manually supervised.",
				strings.Join(unsupervisedComanagedDevices, ", ")),
		)
	}

	if len(notInLostModeManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Managed Devices Not in Lost Mode",
			fmt.Sprintf("The following managed devices are not in lost mode: %s. "+
				"Devices must be in lost mode before playing the lost mode sound. "+
				"Please enable lost mode on these devices first.",
				strings.Join(notInLostModeManagedDevices, ", ")),
		)
	}

	if len(notInLostModeComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Co-Managed Devices Not in Lost Mode",
			fmt.Sprintf("The following co-managed devices are not in lost mode: %s. "+
				"Devices must be in lost mode before playing the lost mode sound. "+
				"Please enable lost mode on these devices first.",
				strings.Join(notInLostModeComanagedDevices, ", ")),
		)
	}
}

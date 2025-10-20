package graphBetaDisableLostModeManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *DisableLostModeManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data DisableLostModeManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	// Get managed device IDs
	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
	}

	// Get co-managed device IDs
	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one device list is provided
	if len(managedDeviceIDs) == 0 && len(comanagedDeviceIDs) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_device_ids' or 'comanaged_device_ids' must be provided with at least one device ID.",
		)
		return
	}

	if len(managedDeviceIDs) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, id := range managedDeviceIDs {
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_device_ids"),
				"Duplicate Managed Device IDs Found",
				fmt.Sprintf("The following managed device IDs are duplicated in the configuration: %s. "+
					"Lost mode will only be disabled once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	if len(comanagedDeviceIDs) > 0 {
		seen := make(map[string]bool)
		var duplicates []string
		for _, id := range comanagedDeviceIDs {
			if seen[id] {
				duplicates = append(duplicates, id)
			}
			seen[id] = true
		}

		if len(duplicates) > 0 {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("comanaged_device_ids"),
				"Duplicate Co-Managed Device IDs Found",
				fmt.Sprintf("The following co-managed device IDs are duplicated in the configuration: %s. "+
					"Lost mode will only be disabled once per device, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	for _, managedID := range managedDeviceIDs {
		for _, comanagedID := range comanagedDeviceIDs {
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_device_ids"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_device_ids and comanaged_device_ids. "+
						"A device should only be in one list. Lost mode disable will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating disable lost mode action for %d managed and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string
	var unsupervisedManagedDevices []string
	var unsupervisedComanagedDevices []string
	var notInLostModeManagedDevices []string
	var notInLostModeComanagedDevices []string

	// Validate managed devices
	for _, deviceID := range managedDeviceIDs {
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
					path.Root("managed_device_ids"),
					"Error Validating Managed Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility - lost mode works on iOS/iPadOS and ChromeOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				supportedOS := map[string]bool{
					"ios":      true,
					"ipados":   true,
					"chromeos": true,
				}
				if !supportedOS[os] {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
					continue
				}
			} else {
				unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}

			// Check if device is supervised (required for lost mode)
			if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
				unsupervisedManagedDevices = append(unsupervisedManagedDevices, deviceID)
			}

			// Check if device is actually in lost mode
			if device.GetLostModeState() != nil {
				lostModeState := device.GetLostModeState().String()
				if lostModeState == "disabled" {
					notInLostModeManagedDevices = append(notInLostModeManagedDevices, fmt.Sprintf("%s (state: %s)", deviceID, lostModeState))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	// Validate co-managed devices
	for _, deviceID := range comanagedDeviceIDs {
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
					path.Root("comanaged_device_ids"),
					"Error Validating Co-Managed Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check platform compatibility - lost mode works on iOS/iPadOS and ChromeOS
			if device.GetOperatingSystem() != nil {
				os := strings.ToLower(*device.GetOperatingSystem())
				supportedOS := map[string]bool{
					"ios":      true,
					"ipados":   true,
					"chromeos": true,
				}
				if !supportedOS[os] {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, *device.GetOperatingSystem()))
					continue
				}
			} else {
				unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (Unknown OS)", deviceID))
				continue
			}

			// Check if device is supervised (required for lost mode)
			if device.GetIsSupervised() == nil || !*device.GetIsSupervised() {
				unsupervisedComanagedDevices = append(unsupervisedComanagedDevices, deviceID)
			}

			// Check if device is actually in lost mode
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
			path.Root("managed_device_ids"),
			"Non-Existent Managed Devices",
			fmt.Sprintf("The following managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentManagedDevices, ", ")),
		)
	}

	if len(nonExistentComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Non-Existent Co-Managed Devices",
			fmt.Sprintf("The following co-managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing co-managed devices.",
				strings.Join(nonExistentComanagedDevices, ", ")),
		)
	}

	if len(unsupportedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_device_ids"),
			"Unsupported Managed Devices for Lost Mode",
			fmt.Sprintf("Lost mode is only supported on iOS, iPadOS, and ChromeOS devices. The following managed devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Unsupported Co-Managed Devices for Lost Mode",
			fmt.Sprintf("Lost mode is only supported on iOS, iPadOS, and ChromeOS devices. The following co-managed devices are not supported: %s. "+
				"Please remove unsupported devices from the configuration.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	if len(unsupervisedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_device_ids"),
			"Unsupervised Managed Devices",
			fmt.Sprintf("Lost mode requires supervised iOS/iPadOS devices. The following managed devices are not supervised: %s. "+
				"Please ensure devices are enrolled via DEP/ABM or manually supervised.",
				strings.Join(unsupervisedManagedDevices, ", ")),
		)
	}

	if len(unsupervisedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_device_ids"),
			"Unsupervised Co-Managed Devices",
			fmt.Sprintf("Lost mode requires supervised iOS/iPadOS devices. The following co-managed devices are not supervised: %s. "+
				"Please ensure devices are enrolled via DEP/ABM or manually supervised.",
				strings.Join(unsupervisedComanagedDevices, ", ")),
		)
	}

	if len(notInLostModeManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_device_ids"),
			"Managed Devices Not in Lost Mode",
			fmt.Sprintf("The following managed devices do not appear to be in lost mode: %s. "+
				"Disabling lost mode on these devices may have no effect. "+
				"Please verify the current lost mode state of these devices before proceeding.",
				strings.Join(notInLostModeManagedDevices, ", ")),
		)
	}

	if len(notInLostModeComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("comanaged_device_ids"),
			"Co-Managed Devices Not in Lost Mode",
			fmt.Sprintf("The following co-managed devices do not appear to be in lost mode: %s. "+
				"Disabling lost mode on these devices may have no effect. "+
				"Please verify the current lost mode state of these devices before proceeding.",
				strings.Join(notInLostModeComanagedDevices, ", ")),
		)
	}
}

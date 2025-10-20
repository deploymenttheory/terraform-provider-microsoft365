package graphBetaActivateDeviceEsimManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (a *ActivateDeviceEsimManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ActivateDeviceEsimManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Basic validation (always run)
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	// No additional configuration validation needed for simplified model

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
					"eSIM activation will only be attempted once per device, but you should remove duplicates from your configuration.",
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
					"eSIM activation will only be attempted once per device, but you should remove duplicates from your configuration.",
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
						"A device should only be in one list. eSIM activation will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating activate device eSIM action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	// Only perform online validation if enabled and client is configured
	validateExists := true
	if !data.ValidateDeviceExists.IsNull() && !data.ValidateDeviceExists.IsUnknown() {
		validateExists = data.ValidateDeviceExists.ValueBool()
	}

	if !validateExists {
		tflog.Debug(ctx, "Device existence validation disabled, skipping online validation")
		return
	}

	if a.client == nil {
		tflog.Debug(ctx, "Client not configured, skipping device existence validation")
		return
	}

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedDeviceTypeManagedDevices []string
	var unsupportedDeviceTypeComanagedDevices []string

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
		} else {
			// Validate deviceType supports eSIM activation (iOS/iPadOS only)
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				supportsESIM := deviceType == models.IPHONE_DEVICETYPE ||
					deviceType == models.IPAD_DEVICETYPE

				if !supportsESIM {
					unsupportedDeviceTypeManagedDevices = append(unsupportedDeviceTypeManagedDevices,
						fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Managed device %s (deviceType: %s) supports eSIM", deviceID, deviceType.String()))
				}
			} else {
				unsupportedDeviceTypeManagedDevices = append(unsupportedDeviceTypeManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
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
		} else {
			// Validate deviceType supports eSIM activation (iOS/iPadOS only)
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				supportsESIM := deviceType == models.IPHONE_DEVICETYPE ||
					deviceType == models.IPAD_DEVICETYPE

				if !supportsESIM {
					unsupportedDeviceTypeComanagedDevices = append(unsupportedDeviceTypeComanagedDevices,
						fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s (deviceType: %s) supports eSIM", deviceID, deviceType.String()))
				}
			} else {
				unsupportedDeviceTypeComanagedDevices = append(unsupportedDeviceTypeComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
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

	if len(unsupportedDeviceTypeManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_devices"),
			"Unsupported Device Type for eSIM on Managed Devices",
			fmt.Sprintf("The following managed devices do not support eSIM activation: %s. "+
				"eSIM activation is only supported on iPhone (8) and iPad (9) devices.",
				strings.Join(unsupportedDeviceTypeManagedDevices, ", ")),
		)
	}

	if len(unsupportedDeviceTypeComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("comanaged_devices"),
			"Unsupported Device Type for eSIM on Co-Managed Devices",
			fmt.Sprintf("The following co-managed devices do not support eSIM activation: %s. "+
				"eSIM activation is only supported on iPhone (8) and iPad (9) devices.",
				strings.Join(unsupportedDeviceTypeComanagedDevices, ", ")),
		)
	}
}

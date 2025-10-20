package graphBetaCleanWindowsManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (a *CleanWindowsManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data CleanWindowsManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Basic validation - at least one device must be specified
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	// Check for duplicate managed device IDs
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
					"Each device will only be cleaned once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for duplicate co-managed device IDs
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
					"Each device will only be cleaned once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for devices in both lists
	for _, managedDevice := range data.ManagedDevices {
		managedID := managedDevice.DeviceID.ValueString()
		for _, comanagedDevice := range data.ComanagedDevices {
			comanagedID := comanagedDevice.DeviceID.ValueString()
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_devices"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
						"A device should only be in one list. Clean operation will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating clean Windows device action for %d managed and %d co-managed device(s)",
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
	var nonWindowsManagedDevices []string
	var unsupportedVersionManagedDevices []string

	// Validate managed devices
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
			// Validate deviceType supports Windows clean operation
			// Valid types: desktop, windowsRT, windows10x, cloudPC
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
					deviceType == models.WINDOWSRT_DEVICETYPE ||
					deviceType == models.WINDOWS10X_DEVICETYPE ||
					deviceType == models.CLOUDPC_DEVICETYPE

				if !isWindowsDevice {
					nonWindowsManagedDevices = append(nonWindowsManagedDevices,
						fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				} else {
					// For desktop and windowsRT, verify OS version is Windows 10+
					if (deviceType == models.DESKTOP_DEVICETYPE || deviceType == models.WINDOWSRT_DEVICETYPE) &&
						device.GetOsVersion() != nil {
						osVersion := *device.GetOsVersion()
						if !strings.HasPrefix(osVersion, "10.") {
							unsupportedVersionManagedDevices = append(unsupportedVersionManagedDevices,
								fmt.Sprintf("%s (deviceType: %s, osVersion: %s)", deviceID, deviceType.String(), osVersion))
						}
					}
				}
			} else {
				nonWindowsManagedDevices = append(nonWindowsManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	var nonExistentComanagedDevices []string
	var nonWindowsComanagedDevices []string
	var unsupportedVersionComanagedDevices []string

	// Validate co-managed devices
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
			// Validate deviceType supports Windows clean operation
			// Valid types: desktop, windowsRT, windows10x, cloudPC
			if device.GetDeviceType() != nil {
				deviceType := *device.GetDeviceType()

				isWindowsDevice := deviceType == models.DESKTOP_DEVICETYPE ||
					deviceType == models.WINDOWSRT_DEVICETYPE ||
					deviceType == models.WINDOWS10X_DEVICETYPE ||
					deviceType == models.CLOUDPC_DEVICETYPE

				if !isWindowsDevice {
					nonWindowsComanagedDevices = append(nonWindowsComanagedDevices,
						fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				} else {
					// For desktop and windowsRT, verify OS version is Windows 10+
					if (deviceType == models.DESKTOP_DEVICETYPE || deviceType == models.WINDOWSRT_DEVICETYPE) &&
						device.GetOsVersion() != nil {
						osVersion := *device.GetOsVersion()
						if !strings.HasPrefix(osVersion, "10.") {
							unsupportedVersionComanagedDevices = append(unsupportedVersionComanagedDevices,
								fmt.Sprintf("%s (deviceType: %s, osVersion: %s)", deviceID, deviceType.String(), osVersion))
						}
					}
				}
			} else {
				nonWindowsComanagedDevices = append(nonWindowsComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
			}
			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully", deviceID))
		}
	}

	// Report validation errors
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

	if len(nonWindowsManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Non-Windows Managed Devices",
			fmt.Sprintf("The following managed devices are not Windows devices: %s. "+
				"Clean operation only supports Windows 10 and Windows 11.",
				strings.Join(nonWindowsManagedDevices, ", ")),
		)
	}

	if len(nonWindowsComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Non-Windows Co-Managed Devices",
			fmt.Sprintf("The following co-managed devices are not Windows devices: %s. "+
				"Clean operation only supports Windows 10 and Windows 11.",
				strings.Join(nonWindowsComanagedDevices, ", ")),
		)
	}

	if len(unsupportedVersionManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_devices"),
			"Potentially Unsupported Windows Versions",
			fmt.Sprintf("The following managed devices may be running unsupported Windows versions: %s. "+
				"Clean operation is designed for Windows 10 and Windows 11.",
				strings.Join(unsupportedVersionManagedDevices, ", ")),
		)
	}

	if len(unsupportedVersionComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("comanaged_devices"),
			"Potentially Unsupported Windows Versions",
			fmt.Sprintf("The following co-managed devices may be running unsupported Windows versions: %s. "+
				"Clean operation is designed for Windows 10 and Windows 11.",
				strings.Join(unsupportedVersionComanagedDevices, ", ")),
		)
	}
}

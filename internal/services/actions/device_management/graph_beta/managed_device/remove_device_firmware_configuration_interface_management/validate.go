package graphBetaRemoveDeviceFirmwareConfigurationInterfaceManagementManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RemoveDeviceFirmwareConfigurationInterfaceManagementManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert framework lists to Go slices
	var managedDeviceIDs []string
	var comanagedDeviceIDs []string

	if !data.ManagedDeviceIDs.IsNull() && !data.ManagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ManagedDeviceIDs.ElementsAs(ctx, &managedDeviceIDs, false)...)
	}

	if !data.ComanagedDeviceIDs.IsNull() && !data.ComanagedDeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.ComanagedDeviceIDs.ElementsAs(ctx, &comanagedDeviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

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
					"DFCI removal will only be performed once per device, but you should remove duplicates from your configuration.",
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
					"DFCI removal will only be performed once per device, but you should remove duplicates from your configuration.",
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
						"A device should only be in one list. DFCI removal will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating remove DFCI management action for %d managed and %d co-managed device(s)",
		len(managedDeviceIDs), len(comanagedDeviceIDs)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string

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
			// Check if device is Windows (DFCI is Windows-specific)
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

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
			// Check if device is Windows
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
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
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_device_ids"),
			"Non-Windows Devices Detected",
			fmt.Sprintf("The following managed devices are not Windows devices: %s. "+
				"DFCI (Device Firmware Configuration Interface) is only supported on Windows devices with compatible firmware. "+
				"These devices will be skipped.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("comanaged_device_ids"),
			"Non-Windows Devices Detected",
			fmt.Sprintf("The following co-managed devices are not Windows devices: %s. "+
				"DFCI (Device Firmware Configuration Interface) is only supported on Windows devices with compatible firmware. "+
				"These devices will be skipped.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	// Add important warning about DFCI removal
	if len(managedDeviceIDs)+len(comanagedDeviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"DFCI Management Removal",
			fmt.Sprintf("This action will remove DFCI (Device Firmware Configuration Interface) management from %d device(s).\n\n"+
				"Important considerations:\n"+
				"- Devices must have DFCI-capable firmware (Surface devices, select OEM devices)\n"+
				"- After removal, UEFI/BIOS settings can no longer be managed remotely via Intune\n"+
				"- Standard MDM management continues - only firmware-level control is removed\n"+
				"- This action cannot be easily reversed\n"+
				"- Physical access may be required to re-enable DFCI on the device\n"+
				"- Ensure you understand the implications before proceeding",
				len(managedDeviceIDs)+len(comanagedDeviceIDs)),
		)
	}
}

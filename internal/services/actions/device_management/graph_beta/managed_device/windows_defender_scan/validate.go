package graphBetaWindowsDefenderScan

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *WindowsDefenderScanAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data WindowsDefenderScanActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one device list is provided
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
					"Each device will only be scanned once, but you should remove duplicates from your configuration.",
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
					"Each device will only be scanned once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	managedIDs := make(map[string]bool)
	for _, device := range data.ManagedDevices {
		managedIDs[device.DeviceID.ValueString()] = true
	}

	for _, device := range data.ComanagedDevices {
		id := device.DeviceID.ValueString()
		if managedIDs[id] {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("managed_devices"),
				"Device ID in Both Lists",
				fmt.Sprintf("Device ID %s appears in both managed_devices and comanaged_devices. "+
					"A device should only be in one list. The scan will be attempted for both endpoints, "+
					"but one may fail if the device is not actually of that type.",
					id),
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating Windows Defender scan action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentManagedDevices []string
	var nonWindowsManagedDevices []string
	var nonExistentComanagedDevices []string
	var nonWindowsComanagedDevices []string

	// Validate managed devices exist and are Windows
	for _, device := range data.ManagedDevices {
		deviceID := device.DeviceID.ValueString()
		managedDevice, err := a.client.
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
		} else if managedDevice != nil {
			// Check that device is Windows
			if managedDevice.GetOperatingSystem() != nil {
				os := strings.ToLower(*managedDevice.GetOperatingSystem())
				if os != "windows" {
					nonWindowsManagedDevices = append(nonWindowsManagedDevices,
						fmt.Sprintf("%s (OS: %s)", deviceID, *managedDevice.GetOperatingSystem()))
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully", deviceID))
		}
	}

	// Validate co-managed devices exist and are Windows
	for _, device := range data.ComanagedDevices {
		deviceID := device.DeviceID.ValueString()
		comanagedDevice, err := a.client.
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
		} else if comanagedDevice != nil {
			// Check that device is Windows
			if comanagedDevice.GetOperatingSystem() != nil {
				os := strings.ToLower(*comanagedDevice.GetOperatingSystem())
				if os != "windows" {
					nonWindowsComanagedDevices = append(nonWindowsComanagedDevices,
						fmt.Sprintf("%s (OS: %s)", deviceID, *comanagedDevice.GetOperatingSystem()))
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

	if len(nonWindowsManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Non-Windows Devices",
			fmt.Sprintf("The Windows Defender scan action only works on Windows devices. "+
				"The following managed devices are not Windows devices: %s. "+
				"Please remove non-Windows devices from the managed_devices list.",
				strings.Join(nonWindowsManagedDevices, ", ")),
		)
	}

	if len(nonWindowsComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Non-Windows Co-Managed Devices",
			fmt.Sprintf("The Windows Defender scan action only works on Windows devices. "+
				"The following co-managed devices are not Windows devices: %s. "+
				"Please remove non-Windows devices from the comanaged_devices list.",
				strings.Join(nonWindowsComanagedDevices, ", ")),
		)
	}

	// Count scan types for informational purposes
	quickScanCount := 0
	fullScanCount := 0
	for _, device := range data.ManagedDevices {
		if device.QuickScan.ValueBool() {
			quickScanCount++
		} else {
			fullScanCount++
		}
	}
	for _, device := range data.ComanagedDevices {
		if device.QuickScan.ValueBool() {
			quickScanCount++
		} else {
			fullScanCount++
		}
	}

	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)

	// Informational message about scan behavior
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_devices"),
		"Windows Defender Scan Information",
		fmt.Sprintf("This action will initiate Windows Defender antivirus scans on %d device(s): "+
			"%d quick scan(s) and %d full scan(s).\n\n"+
			"**Scan Duration:**\n"+
			"- Quick scans: 5-15 minutes (scans common threat locations)\n"+
			"- Full scans: 30+ minutes to hours (comprehensive scan of all files)\n\n"+
			"**Important:**\n"+
			"- Devices must be online to receive the scan command\n"+
			"- Full scans can significantly impact device performance\n"+
			"- Scans run in the background and users can continue working\n"+
			"- Results are reported to Microsoft Intune admin center\n"+
			"- Threats found will be quarantined automatically\n"+
			"- Consider device usage patterns when scheduling full scans",
			totalDevices, quickScanCount, fullScanCount),
	)

	// Warning about full scans
	if fullScanCount > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_devices"),
			"Full Scan Performance Impact",
			fmt.Sprintf("%d device(s) are configured for full scan. Full scans are comprehensive but can take 30+ minutes to several hours "+
				"and may impact device performance during scanning. Users may experience slower system response. "+
				"Consider using quick scans for routine checks and reserve full scans for security incidents or off-hours.",
				fullScanCount),
		)
	}
}

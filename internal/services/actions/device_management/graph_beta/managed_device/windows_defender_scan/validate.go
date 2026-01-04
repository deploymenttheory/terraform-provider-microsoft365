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

	// Check for duplicate device IDs in managed devices
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

	// Check for duplicate device IDs in co-managed devices
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

	// Check for devices appearing in both lists
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

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(data.ManagedDevices),
		"comanaged_count": len(data.ComanagedDevices),
		"total_devices":   totalDevices,
		"quick_scan":      quickScanCount,
		"full_scan":       fullScanCount,
	})
}

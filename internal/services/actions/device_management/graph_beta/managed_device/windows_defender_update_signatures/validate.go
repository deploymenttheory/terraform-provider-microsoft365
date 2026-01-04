package graphBetaWindowsDefenderUpdateSignatures

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *WindowsDefenderUpdateSignaturesAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data WindowsDefenderUpdateSignaturesActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

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

	// Validate that at least one device list is provided
	if len(managedDeviceIDs) == 0 && len(comanagedDeviceIDs) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_device_ids' or 'comanaged_device_ids' must be provided with at least one device ID.",
		)
		return
	}

	// Check for duplicate device IDs in managed devices
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
					"Each device will only have signatures updated once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for duplicate device IDs in co-managed devices
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
					"Each device will only have signatures updated once, but you should remove duplicates from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for devices appearing in both lists
	for _, managedID := range managedDeviceIDs {
		for _, comanagedID := range comanagedDeviceIDs {
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_device_ids"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_device_ids and comanaged_device_ids. "+
						"A device should only be in one list. The signature update will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)

	// Informational message about signature update behavior
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"Windows Defender Signature Update Information",
		fmt.Sprintf("This action will force %d Windows device(s) to immediately update their Windows Defender antivirus signatures. "+
			"Devices must be online with internet connectivity to download the latest threat definitions. "+
			"The update process takes 1-5 minutes and runs in the background with minimal performance impact. "+
			"No device reboot is required. Updated signatures provide protection against the latest threats and malware. "+
			"This is useful before running antivirus scans or in response to new threat intelligence.",
			totalDevices),
	)

	// Best practice recommendation for large batches
	if totalDevices > 100 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_device_ids"),
			"Large Batch Update",
			fmt.Sprintf("You are updating signatures on %d devices. While this is supported, consider: "+
				"1) Staggering updates to reduce network load "+
				"2) Monitoring for failed updates "+
				"3) Ensuring adequate bandwidth "+
				"4) Allowing time for all devices to complete. "+
				"Large batches may take longer to process and report status.",
				totalDevices),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(managedDeviceIDs),
		"comanaged_count": len(comanagedDeviceIDs),
		"total_devices":   totalDevices,
	})
}

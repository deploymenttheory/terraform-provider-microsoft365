package graphBetaSyncManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *SyncManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data SyncManagedDeviceActionModel

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
					"Each device will only be synced once, but you should remove duplicates from your configuration.",
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
					"Each device will only be synced once, but you should remove duplicates from your configuration.",
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
						"A device should only be in one list. The sync will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)

	// Informational message about sync behavior
	resp.Diagnostics.AddAttributeWarning(
		path.Root("managed_device_ids"),
		"Device Sync Information",
		fmt.Sprintf("This action will force %d device(s) to immediately check in with Intune. "+
			"Devices must be online to receive the sync command. Offline devices will sync when they come back online. "+
			"The sync process applies the latest policies, configurations, and updates. "+
			"Multiple syncs in a short period may queue and delay each other. "+
			"Normal check-in interval is every 8 hours; this action forces immediate sync (within 1-5 minutes for online devices).",
			totalDevices),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(managedDeviceIDs),
		"comanaged_count": len(comanagedDeviceIDs),
		"total_devices":   totalDevices,
	})
}

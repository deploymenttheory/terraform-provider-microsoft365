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

	// Check for devices appearing in both lists
	if len(managedDeviceIDs) > 0 && len(comanagedDeviceIDs) > 0 {
		for _, managedID := range managedDeviceIDs {
			for _, comanagedID := range comanagedDeviceIDs {
				if managedID == comanagedID {
					resp.Diagnostics.AddWarning(
						"Device ID in Both Lists",
						fmt.Sprintf("Device ID %s appears in both managed_device_ids and comanaged_device_ids. "+
							"A device should only be in one list. Lost mode disable will be attempted for both endpoints, "+
							"but one may fail if the device is not actually of that type.",
							managedID),
					)
				}
			}
		}
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_devices":   len(managedDeviceIDs),
		"comanaged_devices": len(comanagedDeviceIDs),
	})
}

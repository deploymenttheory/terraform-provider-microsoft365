package graphBetaMoveDevicesToOUManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *MoveDevicesToOUManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data MoveDevicesToOUManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate organizational_unit_path is provided
	if data.OrganizationalUnitPath.IsNull() || data.OrganizationalUnitPath.IsUnknown() || data.OrganizationalUnitPath.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organizational_unit_path"),
			"Organizational Unit Path Required",
			"The organizational_unit_path attribute must be provided with a valid Active Directory OU distinguished name.",
		)
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

	// Validate at least one device list is provided
	if len(managedDeviceIDs) == 0 && len(comanagedDeviceIDs) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_device_ids' or 'comanaged_device_ids' must be provided with at least one device ID.",
		)
		return
	}

	// Check for duplicate managed device IDs
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
				fmt.Sprintf("The following managed device IDs are duplicated in managed_device_ids: %s. "+
					"Duplicates will be ignored, but you should remove them from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for duplicate co-managed device IDs
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
				fmt.Sprintf("The following co-managed device IDs are duplicated in comanaged_device_ids: %s. "+
					"Duplicates will be ignored, but you should remove them from your configuration.",
					strings.Join(duplicates, ", ")),
			)
		}
	}

	// Check for devices in both lists
	for _, managedID := range managedDeviceIDs {
		for _, comanagedID := range comanagedDeviceIDs {
			if managedID == comanagedID {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("managed_device_ids"),
					"Device ID in Both Lists",
					fmt.Sprintf("Device ID %s appears in both managed_device_ids and comanaged_device_ids. "+
						"A device should only be in one list. The move will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	// Add informational note about the operation
	ouPath := data.OrganizationalUnitPath.ValueString()
	if len(managedDeviceIDs)+len(comanagedDeviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"Active Directory OU Move Operation",
			fmt.Sprintf("This action will move %d device(s) to the Active Directory OU: %s\n\n"+
				"Important notes:\n"+
				"- The OU must exist in your on-premises Active Directory\n"+
				"- Azure AD Connect must have permissions to move computer objects to this OU\n"+
				"- Changes will reflect after the next Azure AD Connect sync cycle\n"+
				"- Only hybrid Azure AD joined Windows devices can be moved\n"+
				"- Ensure the OU path is correct as this cannot be easily reverted",
				len(managedDeviceIDs)+len(comanagedDeviceIDs), ouPath),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(managedDeviceIDs),
		"comanaged_count": len(comanagedDeviceIDs),
		"ou_path":         ouPath,
	})
}

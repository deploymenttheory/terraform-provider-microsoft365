package graphBetaRotateLocalAdminPasswordManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RotateLocalAdminPasswordManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RotateLocalAdminPasswordManagedDeviceActionModel

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
					"Password rotation will only be performed once per device, but you should remove duplicates from your configuration.",
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
					"Password rotation will only be performed once per device, but you should remove duplicates from your configuration.",
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
						"A device should only be in one list. Password rotation will be attempted for both endpoints, "+
						"but one may fail if the device is not actually of that type.",
						managedID),
				)
			}
		}
	}

	totalDevices := len(managedDeviceIDs) + len(comanagedDeviceIDs)

	// Add important warning about the operation
	resp.Diagnostics.AddWarning(
		"Local Administrator Password Rotation",
		fmt.Sprintf("This action will rotate the local administrator password for %d device(s).\n\n"+
			"Important notes:\n"+
			"- Devices must have Windows LAPS policy configured and enabled\n"+
			"- New password will be automatically generated (complex, random)\n"+
			"- Password stored securely in Azure AD or Intune\n"+
			"- Authorized administrators can retrieve the new password\n"+
			"- Does not affect device operation or require restart\n"+
			"- Critical for security compliance and privileged access management\n"+
			"- Previous password will no longer work for local admin access",
			totalDevices),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(managedDeviceIDs),
		"comanaged_count": len(comanagedDeviceIDs),
		"total_devices":   totalDevices,
	})
}

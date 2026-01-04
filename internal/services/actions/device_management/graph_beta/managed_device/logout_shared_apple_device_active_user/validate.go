package graphBetaLogoutSharedAppleDeviceActiveUser

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *LogoutSharedAppleDeviceActiveUserAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data LogoutSharedAppleDeviceActiveUserActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert framework list to Go slice
	var deviceIDs []string
	if !data.DeviceIDs.IsNull() && !data.DeviceIDs.IsUnknown() {
		resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one device is provided
	if len(deviceIDs) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one device ID must be provided in 'device_ids'.",
		)
		return
	}

	// Check for duplicate device IDs
	seen := make(map[string]bool)
	var duplicates []string
	for _, id := range deviceIDs {
		if seen[id] {
			duplicates = append(duplicates, id)
		}
		seen[id] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs Found",
			fmt.Sprintf("The following device IDs are duplicated in device_ids: %s. "+
				"Logout will only be performed once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	// General warning about Shared iPad mode requirement
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"Shared iPad Mode Requirement",
		fmt.Sprintf("This action only works on iPads configured in Shared iPad mode. "+
			"Regular (non-shared) iPads will not be affected by this action, even if they meet other requirements. "+
			"Ensure the %d device(s) in this action are actually configured in Shared iPad mode. "+
			"The action will fail gracefully if devices are not in Shared iPad mode or if no user is currently logged in.",
			len(deviceIDs)),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"device_count": len(deviceIDs),
	})
}

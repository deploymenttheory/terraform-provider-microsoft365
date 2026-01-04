package graphBetaPauseConfigurationRefreshManagedDevice

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *PauseConfigurationRefreshManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data PauseConfigurationRefreshManagedDeviceActionModel

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

	// Add informational warning about the pause operation
	totalDevices := len(data.ManagedDevices) + len(data.ComanagedDevices)
	if totalDevices > 0 {
		resp.Diagnostics.AddWarning(
			"Configuration Refresh Pause",
			fmt.Sprintf("This action will pause configuration refresh for %d device(s).\n\n"+
				"Important notes:\n"+
				"- Devices will not receive new policy updates during the pause period\n"+
				"- Existing applied policies remain in effect\n"+
				"- Configuration refresh automatically resumes after the pause period\n"+
				"- Users can still manually sync from Company Portal\n"+
				"- Critical security updates may still be applied\n"+
				"- Use this feature during maintenance windows or troubleshooting only",
				totalDevices),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"managed_count":   len(data.ManagedDevices),
		"comanaged_count": len(data.ComanagedDevices),
		"total_devices":   totalDevices,
	})
}

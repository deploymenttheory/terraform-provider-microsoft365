package graphBetaRecoverManagedDevicePasscode

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RecoverManagedDevicePasscodeAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RecoverManagedDevicePasscodeActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
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
			fmt.Sprintf("The following device IDs are duplicated in the configuration: %s. "+
				"Passcode recovery will only be attempted once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"total_devices": len(deviceIDs),
	})
}

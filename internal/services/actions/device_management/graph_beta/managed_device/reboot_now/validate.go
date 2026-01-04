package graphBetaRebootNowManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RebootNowManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RebootNowManagedDeviceActionModel

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
				"Reboot command will only be sent once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	// General warning about user impact
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"User Impact Warning",
		fmt.Sprintf("Rebooting %d device(s) will immediately restart them when online. "+
			"Users may lose unsaved work and active sessions will be terminated. "+
			"Consider scheduling this action during maintenance windows or notifying users in advance.",
			len(deviceIDs)),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"total_devices": len(deviceIDs),
	})
}

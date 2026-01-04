package graphBetaShutdownManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ShutdownManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ShutdownManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert framework list to Go slice
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
				"Shutdown command will only be sent once per device, but you should remove duplicates from your configuration.",
				strings.Join(duplicates, ", ")),
		)
	}

	// Critical warning about shutdown requiring manual power-on
	resp.Diagnostics.AddAttributeWarning(
		path.Root("device_ids"),
		"Critical: Manual Power-On Required",
		fmt.Sprintf("Shutting down %d device(s) will POWER THEM OFF COMPLETELY. "+
			"Physical access will be required to power devices back on. "+
			"Users may lose unsaved work and will be unable to access their devices until manually powered on. "+
			"Consider using reboot action instead if devices need to come back online automatically. "+
			"Ensure you have legitimate business reason and proper authorization for this disruptive action.",
			len(deviceIDs)),
	)

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"device_count": len(deviceIDs),
	})
}

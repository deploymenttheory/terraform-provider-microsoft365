package graphBetaResetManagedDevicePasscode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *ResetManagedDevicePasscodeAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data ResetManagedDevicePasscodeActionModel

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
	deviceIDMap := make(map[string]bool)
	var duplicates []string

	for _, deviceID := range deviceIDs {
		if deviceIDMap[deviceID] {
			duplicates = append(duplicates, deviceID)
		}
		deviceIDMap[deviceID] = true
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("device_ids"),
			"Duplicate Device IDs",
			fmt.Sprintf("The following device IDs are duplicated in the list: %v. "+
				"Each device passcode will only be reset once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	// General warning about device connectivity and passcode retrieval
	if len(deviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"Device Connectivity Required",
			fmt.Sprintf("Resetting passcodes for %d device(s). "+
				"Devices must be online and connected to receive the reset passcode command. "+
				"If a device is offline, the command will be queued and executed when the device next checks in. "+
				"After successful reset, you must retrieve the new temporary passcodes from the Intune portal "+
				"and communicate them securely to the device users.",
				len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"total_devices": len(deviceIDs),
	})
}

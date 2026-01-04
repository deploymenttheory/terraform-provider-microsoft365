package graphBetaRemoteLockManagedDevice

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *RemoteLockManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RemoteLockManagedDeviceActionModel

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
				"Each device will only receive the lock command once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	// General warning about device action
	if len(deviceIDs) > 0 {
		resp.Diagnostics.AddWarning(
			"Device Lock Warning",
			fmt.Sprintf("You are about to remotely lock %d device(s). "+
				"Devices will lock immediately when they receive the command. "+
				"Users will need to enter their passcode to unlock. "+
				"Ensure you have authorization to lock these devices. "+
				"If devices are offline, the lock command will be queued and executed when they next check in with Intune.",
				len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"total_devices": len(deviceIDs),
	})
}

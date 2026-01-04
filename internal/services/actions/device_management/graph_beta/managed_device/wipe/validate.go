package graphBetaWipeManagedDevice

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *WipeManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data WipeManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.DeviceIDs.IsNull() || data.DeviceIDs.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Missing Required Configuration",
			"device_ids must be specified and contain at least one device ID.",
		)
		return
	}

	var deviceIDs []string
	resp.Diagnostics.Append(data.DeviceIDs.ElementsAs(ctx, &deviceIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(deviceIDs) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("device_ids"),
			"Invalid Configuration",
			"device_ids must contain at least one device ID.",
		)
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
				"Each device will only be wiped once, but you should remove duplicates from your configuration.",
				duplicates),
		)
	}

	// Validate conflicting options
	if !data.KeepUserData.IsNull() && !data.KeepUserData.IsUnknown() && data.KeepUserData.ValueBool() {
		if !data.ObliterationBehavior.IsNull() && !data.ObliterationBehavior.IsUnknown() {
			obliterationBehaviorValue := data.ObliterationBehavior.ValueString()
			if obliterationBehaviorValue != "doNotObliterate" {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("obliteration_behavior"),
					"Conflicting Configuration",
					"When keep_user_data is true, obliteration_behavior should typically be set to 'doNotObliterate'. "+
						"Using obliteration with keep_user_data=true may result in unexpected behavior.",
				)
			}
		}
	}

	// Validate macOS unlock code format
	if !data.MacOsUnlockCode.IsNull() && !data.MacOsUnlockCode.IsUnknown() {
		macOsUnlockCode := data.MacOsUnlockCode.ValueString()
		if len(macOsUnlockCode) != 6 {
			resp.Diagnostics.AddAttributeError(
				path.Root("macos_unlock_code"),
				"Invalid macOS Unlock Code",
				"The macOS unlock code must be exactly 6 digits.",
			)
		}
	}

	// Final warning about data loss
	if !data.KeepUserData.IsNull() && !data.KeepUserData.IsUnknown() && !data.KeepUserData.ValueBool() {
		resp.Diagnostics.AddWarning(
			"Data Loss Warning",
			fmt.Sprintf("You are about to wipe %d device(s) with keep_user_data=false. "+
				"This will permanently delete ALL data on these devices. "+
				"This action cannot be undone. Please ensure you have reviewed the device list carefully.", len(deviceIDs)),
		)
	}

	tflog.Debug(ctx, "Static validation completed", map[string]any{
		"device_count": len(deviceIDs),
	})
}

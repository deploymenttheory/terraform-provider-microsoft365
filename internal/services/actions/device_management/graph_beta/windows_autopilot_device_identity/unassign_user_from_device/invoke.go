package graphBetaUnassignUserFromDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Invoke performs the action to unassign a user from an Autopilot device.
func (a *UnassignUserFromDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data UnassignUserFromDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	// Read action config data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.WindowsAutopilotDeviceIdentityID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing action %s, unassigning user from device ID: %s", ActionName, deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Unassigning user from Autopilot device...",
	})

	// Execute the unassign user from device operation
	// Note: This API call does not require a request body
	err := a.client.DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceID).
		UnassignUserFromDevice().
		Post(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "UnassignUserFromDevice", a.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully unassigned user from device %s", deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully unassigned user from device %s", deviceID),
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

package graphBetaAssignUserToDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Invoke performs the action to assign a user to an Autopilot device.
func (a *AssignUserToDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data AssignUserToDeviceActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.WindowsAutopilotDeviceIdentityID.ValueString()
	userPrincipalName := data.UserPrincipalName.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing action %s, assigning user %s to device ID: %s", ActionName, userPrincipalName, deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Assigning user to Autopilot device...",
	})

	requestBody, err := constructRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing request",
			fmt.Sprintf("Could not construct request for assign user to device: %s", err.Error()),
		)
		return
	}

	err = a.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceID).
		AssignUserToDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "AssignUserToDevice", a.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned user %s to device %s", userPrincipalName, deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Successfully assigned user %s to device %s", userPrincipalName, deviceID),
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

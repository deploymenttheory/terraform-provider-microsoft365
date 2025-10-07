package graphBetaAllowNextEnrollment

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Invoke performs the action to allow next enrollment for an Autopilot device.
func (a *AllowNextEnrollmentAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data AllowNextEnrollmentActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.WindowsAutopilotDeviceIdentityID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing action %s, allowing next autopilot enrollment for device ID: %s", ActionName, deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Allowing next autopilot enrollment for device %s...", deviceID),
	})

	err := a.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceID).
		AllowNextEnrollment().
		Post(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Action", a.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully allowed next autopilot enrollment for device %s", deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("successfully allowed next autopilot enrollment for device %s", deviceID),
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

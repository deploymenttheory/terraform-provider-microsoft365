package graphBetaUpdateDeviceProperties

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Invoke performs the action to update properties on an Autopilot device.
func (a *UpdateDevicePropertiesAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data UpdateDevicePropertiesActionModel

	tflog.Debug(ctx, fmt.Sprintf("Starting %s", ActionName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := data.WindowsAutopilotDeviceIdentityID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Performing action %s, updating device properties for device ID: %s", ActionName, deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Updating device properties for %s...", deviceID),
	})

	requestBody, err := constructRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing request",
			fmt.Sprintf("Could not construct request for update device properties: %s", err.Error()),
		)
		return
	}

	err = a.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceID).
		UpdateDeviceProperties().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Action", a.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated device properties for device %s", deviceID))

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Device properties updated successfully for %s", deviceID),
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished %s", ActionName))
}

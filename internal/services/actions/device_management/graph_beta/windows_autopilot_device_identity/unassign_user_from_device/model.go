package graphBetaUnassignUserFromDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UnassignUserFromDeviceActionModel struct {
	WindowsAutopilotDeviceIdentityID types.String   `tfsdk:"windows_autopilot_device_identity_id"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}
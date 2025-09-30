package graphBetaAssignUserToDevice

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AssignUserToDeviceActionModel struct {
	WindowsAutopilotDeviceIdentityID types.String   `tfsdk:"windows_autopilot_device_identity_id"`
	UserPrincipalName                types.String   `tfsdk:"user_principal_name"`
	AddressableUserName              types.String   `tfsdk:"addressable_user_name"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}
// REF: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-assignusertodevice?view=graph-rest-beta
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

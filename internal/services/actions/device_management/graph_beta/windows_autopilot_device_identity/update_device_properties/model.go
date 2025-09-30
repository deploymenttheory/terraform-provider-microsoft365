// REF: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-updatedeviceproperties?view=graph-rest-beta
package graphBetaUpdateDeviceProperties

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpdateDevicePropertiesActionModel struct {
	WindowsAutopilotDeviceIdentityID types.String   `tfsdk:"windows_autopilot_device_identity_id"`
	UserPrincipalName                types.String   `tfsdk:"user_principal_name"`
	AddressableUserName              types.String   `tfsdk:"addressable_user_name"`
	GroupTag                         types.String   `tfsdk:"group_tag"`
	DisplayName                      types.String   `tfsdk:"display_name"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}

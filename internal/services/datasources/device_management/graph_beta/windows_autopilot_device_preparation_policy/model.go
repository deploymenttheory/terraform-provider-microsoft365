// https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsAutopilotDevicePreparationPolicyDataSourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.Set      `tfsdk:"role_scope_tag_ids"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

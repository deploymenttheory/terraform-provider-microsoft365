// REF: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta&tabs=go
package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchDeploymentStateResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	DeploymentId   types.String   `tfsdk:"deployment_id"`
	RequestedValue types.String   `tfsdk:"requested_value"`
	EffectiveValue types.String   `tfsdk:"effective_value"`
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
}

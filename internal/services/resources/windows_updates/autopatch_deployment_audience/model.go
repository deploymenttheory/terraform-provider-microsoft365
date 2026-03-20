package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchDeploymentAudienceResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	MemberType types.String   `tfsdk:"member_type"`
	Members    types.Set      `tfsdk:"members"`
	Exclusions types.Set      `tfsdk:"exclusions"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

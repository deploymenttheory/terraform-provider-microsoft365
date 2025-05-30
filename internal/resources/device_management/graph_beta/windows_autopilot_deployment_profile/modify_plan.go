package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modification for Windows Autopilot Deployment Profile resource
func (r *WindowsAutopilotDeploymentProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no plan modifications are needed
	// This can be extended in the future if specific plan modifications are required
}

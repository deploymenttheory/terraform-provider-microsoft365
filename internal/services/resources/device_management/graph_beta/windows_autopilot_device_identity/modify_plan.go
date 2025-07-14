package graphBetaWindowsAutopilotDeviceIdentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modification for Windows Autopilot Device Identity resource
func (r *WindowsAutopilotDeviceIdentityResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no plan modifications are needed
	// This can be extended in the future if specific plan modifications are required
}

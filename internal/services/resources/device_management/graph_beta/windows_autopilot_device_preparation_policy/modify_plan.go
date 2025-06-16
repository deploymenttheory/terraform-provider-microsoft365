package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan allows modification of the plan before it's applied.
func (r *WindowsAutopilotDevicePreparationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Use default values for optional fields that are null
	if req.Plan.Raw.IsNull() {
		return
	}

	// Handle plan modifications if needed for specific fields
	var plan WindowsAutopilotDevicePreparationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No specific plan modifications needed for this resource
}

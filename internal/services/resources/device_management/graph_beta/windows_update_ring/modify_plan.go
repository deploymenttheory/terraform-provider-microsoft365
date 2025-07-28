package graphBetaWindowsUpdateRing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan modifies the planned state of the resource.
func (r *WindowsUpdateRingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Add any plan modifications here
}

package graphBetaMacOSVppApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan modifies the planned resource state before it's applied.
func (r *MacOSVppAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No plan modifications needed for this resource
}

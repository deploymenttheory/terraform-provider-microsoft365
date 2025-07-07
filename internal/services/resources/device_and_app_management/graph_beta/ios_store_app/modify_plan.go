package graphBetaIOSStoreApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modification for the resource.
func (r *IOSStoreAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No custom plan modifications needed
}

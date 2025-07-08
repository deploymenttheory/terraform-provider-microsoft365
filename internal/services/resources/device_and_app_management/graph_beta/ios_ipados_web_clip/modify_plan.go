package graphBetaIOSiPadOSWebClip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan is called when the provider has an opportunity to modify the plan.
func (r *IOSiPadOSWebClipResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No custom plan modifications needed at this time
}

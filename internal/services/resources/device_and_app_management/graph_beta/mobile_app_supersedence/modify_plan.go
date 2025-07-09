package graphBetaMobileAppSupersedence

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// planModifiers returns the plan modifiers for the provider
func (r *MobileAppSupersedenceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No custom plan modifications needed at this time
}

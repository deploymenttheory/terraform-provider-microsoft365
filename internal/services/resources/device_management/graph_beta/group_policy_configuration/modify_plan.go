package graphBetaGroupPolicyConfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles the ModifyPlan operation for Group Policy Configuration resources.
func (r *GroupPolicyConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently, no custom plan modification logic is needed for this resource
	// This method satisfies the ResourceWithModifyPlan interface
}

package graphBetaGroupPolicyMultiTextValue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles the ModifyPlan operation for Group Policy Multi-Text Value resources.
func (r *GroupPolicyMultiTextValueResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently, no custom plan modification logic is needed for this resource
	// This method satisfies the ResourceWithModifyPlan interface
}

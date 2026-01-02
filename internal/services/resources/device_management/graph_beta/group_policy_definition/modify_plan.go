package graphBetaGroupPolicyDefinition

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan satisfies the ResourceWithModifyPlan interface
func (r *GroupPolicyDefinitionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No custom plan modifications needed
	// All ID resolution and validation happens during CRUD operations
}

package graphBetaAppControlForBusinessManagedInstaller

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modifications for the resource.
func (r *AppControlForBusinessManagedInstallerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No special plan modifications needed for this resource
	// The default behavior is sufficient
}
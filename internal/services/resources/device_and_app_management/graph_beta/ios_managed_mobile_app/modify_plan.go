package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modifications for the iOS Managed Mobile App resource.
func (r *IOSManagedMobileAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no custom plan modification logic is required for this resource.
	// This method is implemented to satisfy the resource.ResourceWithModifyPlan interface.
}
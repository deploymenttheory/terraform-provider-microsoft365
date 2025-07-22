package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modifications for the Android Managed Mobile App resource.
func (r *AndroidManagedMobileAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no custom plan modification logic is required for this resource.
	// This method is implemented to satisfy the resource.ResourceWithModifyPlan interface.
}
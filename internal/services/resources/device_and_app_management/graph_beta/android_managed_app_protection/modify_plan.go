package graphBetaDeviceAndAppManagementAndroidManagedAppProtection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modifications for the Android Managed App Protection resource.
func (r *AndroidManagedAppProtectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no custom plan modification logic is required for this resource.
	// This method is implemented to satisfy the resource.ResourceWithModifyPlan interface.
}

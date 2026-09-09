package graphBetaDeviceAndAppManagementIosManagedAppProtection

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modifications for the iOS Managed App Protection resource.
func (r *IosManagedAppProtectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently no custom plan modification logic is required for this resource.
}
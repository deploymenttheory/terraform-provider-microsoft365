package graphBetaDeviceManagementIntuneBrandingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for MacOS PKG apps diff suppression to avoid state inconsistency errors
func (r *IntuneBrandingProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	tflog.Debug(ctx, "Modify Plan Place holder")
}

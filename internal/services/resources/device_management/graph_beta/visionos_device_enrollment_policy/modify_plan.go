package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan implements resource.ResourceWithModifyPlan. No custom plan modification logic is
// currently required beyond what's declared per-attribute in the schema.
func (r *VisionOSDeviceEnrollmentPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "ModifyPlan placeholder")
}

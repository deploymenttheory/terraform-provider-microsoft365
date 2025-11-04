package graphBetaGroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for groups
func (r *GroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Don't validate on destroy
	if req.Plan.Raw.IsNull() {
		return
	}

	tflog.Debug(ctx, "Starting plan modification for group resource")

	var plan GroupResourceModel

	// Get the planned configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform cross-attribute validation
	resp.Diagnostics.Append(ValidateGroupConfiguration(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Finished plan modification for group resource")
}

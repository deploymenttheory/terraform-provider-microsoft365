package graphBetaMacOSPKGApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for diff suppression
// ModifyPlan handles plan modification for diff suppression
func (r *MacOSPKGAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state MacOSPKGAppResourceModel
	var plan MacOSPKGAppResourceModel

	// Get current state and plan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If plan has a MacOSPkgApp but included_apps is null, and state has included_apps,
	// copy included_apps from state to plan
	if plan.MacOSPkgApp != nil && plan.MacOSPkgApp.IncludedApps == nil &&
		state.MacOSPkgApp != nil && state.MacOSPkgApp.IncludedApps != nil {

		tflog.Debug(ctx, "Setting included_apps in plan from state to avoid inconsistency")
		plan.MacOSPkgApp.IncludedApps = state.MacOSPkgApp.IncludedApps

		// Set the modified plan
		resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	}
}

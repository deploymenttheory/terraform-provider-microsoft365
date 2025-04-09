package graphBetaMacOSPKGApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for MacOS PKG apps diff suppression to avoid state inconsistency errors
func (r *MacOSPKGAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If we don't have a plan or we're creating a new resource, nothing to do here
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state MacOSPKGAppResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle included_apps that might be inferred during creation
	if plan.MacOSPkgApp.IncludedApps.IsUnknown() && !state.MacOSPkgApp.IncludedApps.IsNull() {
		tflog.Debug(ctx, "Propagating inferred included_apps from state into plan")
		plan.MacOSPkgApp.IncludedApps = state.MacOSPkgApp.IncludedApps
	}

	// Set the modified plan
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

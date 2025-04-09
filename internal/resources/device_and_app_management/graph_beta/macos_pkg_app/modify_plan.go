package graphBetaMacOSPKGApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for MacOS PKG apps to avoid state inconsistency errors
func (r *MacOSPKGAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// We only need to check for this inconsistency during updates (not creation)
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

	tflog.Debug(ctx, "Checking for values that may have been inferred from PKG analysis")

	// Handle included_apps that might be inferred during creation
	if plan.MacOSPkgApp != nil && state.MacOSPkgApp != nil {
		if len(plan.MacOSPkgApp.IncludedApps) == 0 && len(state.MacOSPkgApp.IncludedApps) > 0 {
			tflog.Debug(ctx, "Setting included_apps in plan to match state since it was inferred during creation",
				map[string]interface{}{
					"stateIncludedAppsCount": len(state.MacOSPkgApp.IncludedApps),
				})

			// Copy the included apps from state to plan
			plan.MacOSPkgApp.IncludedApps = make([]MacOSIncludedAppResourceModel, len(state.MacOSPkgApp.IncludedApps))
			copy(plan.MacOSPkgApp.IncludedApps, state.MacOSPkgApp.IncludedApps)
		}
	}

	// Set the modified plan
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

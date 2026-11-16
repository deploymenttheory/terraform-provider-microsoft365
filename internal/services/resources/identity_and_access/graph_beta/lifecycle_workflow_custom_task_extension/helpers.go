package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringAttributes and objectAttributes are the writable properties sent to Graph.
var (
	stringAttributes = []string{"display_name", "description"}
	objectAttributes = []string{"endpoint_configuration", "authentication_configuration", "client_configuration", "callback_configuration"}
)

// hasGraphChanges compares plan and state and reports whether any property sent to Graph changed.
// An update that only changes timeouts skips the PATCH, so the consistency predicate never waits
// for a last_modified_date_time that will not change.
func hasGraphChanges(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	for _, name := range stringAttributes {
		var p, s types.String
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		if diags.HasError() {
			return false, diags
		}
		if changed(p, s) {
			return true, diags
		}
	}

	for _, name := range objectAttributes {
		var p, s types.Object
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		if diags.HasError() {
			return false, diags
		}
		if changed(p, s) {
			return true, diags
		}
	}

	return false, diags
}

func changed(planValue, stateValue attr.Value) bool {
	return planValue.IsUnknown() || !planValue.Equal(stateValue)
}

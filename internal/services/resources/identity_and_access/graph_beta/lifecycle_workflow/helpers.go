package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// versionedAttributes can only be changed with createNewVersion; PATCH rejects them.
var versionedAttributes = []string{"category", "execution_conditions", "tasks"}

// basicAttributes are the only properties Graph accepts on PATCH.
var basicAttributes = []string{"display_name", "description", "is_enabled", "is_scheduling_enabled"}

// updateKind describes which Graph call an update needs.
type updateKind int

const (
	updateNone updateKind = iota
	updateBasic
	updateNewVersion
)

// planUpdateKind compares plan and state and returns whether the update needs a new workflow
// version, a PATCH of the basic properties, or no Graph write at all (for example when only
// timeouts changed).
func planUpdateKind(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (updateKind, diag.Diagnostics) {
	var diags diag.Diagnostics

	for _, name := range versionedAttributes {
		planValue, stateValue, d := attributePair(ctx, plan, state, name)
		diags.Append(d...)
		if diags.HasError() {
			return updateNone, diags
		}
		if !planValue.Equal(stateValue) {
			return updateNewVersion, diags
		}
	}

	for _, name := range basicAttributes {
		planValue, stateValue, d := attributePair(ctx, plan, state, name)
		diags.Append(d...)
		if diags.HasError() {
			return updateNone, diags
		}
		if !planValue.Equal(stateValue) {
			return updateBasic, diags
		}
	}

	return updateNone, diags
}

func attributePair(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, name string) (attr.Value, attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planValue, stateValue attr.Value

	switch name {
	case "execution_conditions":
		var p, s types.Object
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	case "tasks":
		var p, s types.List
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	case "is_enabled", "is_scheduling_enabled":
		var p, s types.Bool
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	default:
		var p, s types.String
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	}

	return planValue, stateValue, diags
}

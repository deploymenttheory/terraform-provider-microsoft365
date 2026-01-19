package utilityDeploymentScheduler

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setStateDefaults sets default values for optional state fields
func setStateDefaults(state *DeploymentSchedulerDataSourceModel) {
	if state.RequireAllConditions.IsNull() || state.RequireAllConditions.IsUnknown() {
		state.RequireAllConditions = types.BoolValue(true)
	}

	if state.ManualOverride.IsNull() || state.ManualOverride.IsUnknown() {
		state.ManualOverride = types.BoolValue(false)
	}
}

// setReleasedScopeIDs sets the released scope IDs based on whether conditions are met
func setReleasedScopeIDs(state *DeploymentSchedulerDataSourceModel, conditionsMet bool) {
	if conditionsMet {
		// Release the scope ID(s)
		if !state.ScopeId.IsNull() && !state.ScopeId.IsUnknown() {
			// Singular scope_id provided
			state.ReleasedScopeId = state.ScopeId
			state.ReleasedScopeIds = types.ListNull(types.StringType)
		} else if !state.ScopeIds.IsNull() && !state.ScopeIds.IsUnknown() {
			// Multiple scope_ids provided
			state.ReleasedScopeIds = state.ScopeIds
			state.ReleasedScopeId = types.StringNull()
		}
	} else {
		// Don't release - return null
		state.ReleasedScopeId = types.StringNull()
		state.ReleasedScopeIds = types.ListNull(types.StringType)
	}
}

// setStateID sets the resource ID
func setStateID(state *DeploymentSchedulerDataSourceModel) {
	state.Id = types.StringValue(fmt.Sprintf("deployment-scheduler-%s", state.Name.ValueString()))
}

// buildConditionsDetail builds the conditions_detail object for state
func buildConditionsDetail(
	ctx context.Context,
	resp *datasource.ReadResponse,
	timeConditionExists bool,
	delayStartTimeBy int64,
	deploymentStartTimeStr string,
	currentTimeStr string,
	hoursElapsed float64,
	timeConditionMet bool,
) types.Object {
	if !timeConditionExists {
		// No time condition - return null detail
		return types.ObjectNull(
			map[string]attr.Type{
				"time_condition_detail": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"required":              types.BoolType,
						"delay_start_time_by":   types.Int64Type,
						"deployment_start_time": types.StringType,
						"current_time":          types.StringType,
						"hours_elapsed":         types.Float64Type,
						"condition_met":         types.BoolType,
					},
				},
			},
		)
	}

	timeDetailAttrs := map[string]attr.Value{
		"required":              types.BoolValue(true),
		"delay_start_time_by":   types.Int64Value(delayStartTimeBy),
		"deployment_start_time": types.StringValue(deploymentStartTimeStr),
		"current_time":          types.StringValue(currentTimeStr),
		"hours_elapsed":         types.Float64Value(hoursElapsed),
		"condition_met":         types.BoolValue(timeConditionMet),
	}

	timeDetailObj, diags := types.ObjectValue(
		map[string]attr.Type{
			"required":              types.BoolType,
			"delay_start_time_by":   types.Int64Type,
			"deployment_start_time": types.StringType,
			"current_time":          types.StringType,
			"hours_elapsed":         types.Float64Type,
			"condition_met":         types.BoolType,
		},
		timeDetailAttrs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return types.ObjectNull(
			map[string]attr.Type{
				"time_condition_detail": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"required":              types.BoolType,
						"delay_start_time_by":   types.Int64Type,
						"deployment_start_time": types.StringType,
						"current_time":          types.StringType,
						"hours_elapsed":         types.Float64Type,
						"condition_met":         types.BoolType,
					},
				},
			},
		)
	}

	conditionsDetailAttrs := map[string]attr.Value{
		"time_condition_detail": timeDetailObj,
	}

	conditionsDetailObj, diags := types.ObjectValue(
		map[string]attr.Type{
			"time_condition_detail": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"required":              types.BoolType,
					"delay_start_time_by":   types.Int64Type,
					"deployment_start_time": types.StringType,
					"current_time":          types.StringType,
					"hours_elapsed":         types.Float64Type,
					"condition_met":         types.BoolType,
				},
			},
		},
		conditionsDetailAttrs,
	)
	resp.Diagnostics.Append(diags...)

	return conditionsDetailObj
}

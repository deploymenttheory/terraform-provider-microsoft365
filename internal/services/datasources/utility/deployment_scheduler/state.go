package utilityDeploymentScheduler

import (
	"fmt"

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

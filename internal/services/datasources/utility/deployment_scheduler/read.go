package utilityDeploymentScheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Deployment Scheduler data source.
// This function orchestrates the evaluation of deployment conditions.
func (d *deploymentSchedulerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DeploymentSchedulerDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	currentTime := time.Now().UTC()
	currentTimeStr := currentTime.Format(time.RFC3339)

	deploymentStartTime, deploymentStartTimeStr, err := parseDeploymentStartTime(ctx, &state, currentTime, currentTimeStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			fmt.Sprintf("Failed to parse deployment start time: %s", err.Error()),
		)
		return
	}
	state.DeploymentStartTime = types.StringValue(deploymentStartTimeStr)

	setStateDefaults(&state)

	// Check manual override first - bypasses all other conditions
	if state.ManualOverride.ValueBool() {
		handleManualOverride(ctx, &state, resp)
		return
	}

	// Evaluate all conditions
	timeConditionMet, timeConditionExists, delayStartTimeBy, hoursElapsed, _, err := evaluateTimeCondition(ctx, &state, currentTime, deploymentStartTime)
	if err != nil {
		resp.Diagnostics.AddError(
			"Configuration Error",
			fmt.Sprintf("Failed to evaluate time condition: %s", err.Error()),
		)
		return
	}

	inclusionWindowMet, _, err := evaluateInclusionWindows(ctx, currentTime, state.InclusionTimeWindows)
	if err != nil {
		resp.Diagnostics.AddError(
			"Configuration Error",
			fmt.Sprintf("Failed to evaluate inclusion windows: %s", err.Error()),
		)
		return
	}

	exclusionWindowActive, _, err := evaluateExclusionWindows(ctx, currentTime, state.ExclusionTimeWindows)
	if err != nil {
		resp.Diagnostics.AddError(
			"Configuration Error",
			fmt.Sprintf("Failed to evaluate exclusion windows: %s", err.Error()),
		)
		return
	}

	dependencyGateMet, _, err := evaluateDependencyGate(ctx, &state, currentTime, deploymentStartTime)
	if err != nil {
		resp.Diagnostics.AddError(
			"Configuration Error",
			fmt.Sprintf("Failed to evaluate dependency gate: %s", err.Error()),
		)
		return
	}

	// Determine overall condition status
	allConditionsMet := timeConditionMet && inclusionWindowMet && !exclusionWindowActive && dependencyGateMet

	// Update state
	state.ConditionMet = types.BoolValue(allConditionsMet)
	setReleasedScopeIDs(&state, allConditionsMet)
	state.ConditionsDetail = buildConditionsDetail(
		ctx, resp, timeConditionExists, delayStartTimeBy,
		deploymentStartTimeStr, currentTimeStr, hoursElapsed, timeConditionMet,
	)
	setStateID(&state)

	// Log result
	if allConditionsMet {
		if !state.ReleasedScopeId.IsNull() {
			tflog.Info(ctx, "gate open, releasing scope_id", map[string]any{
				"scope_id": state.ReleasedScopeId.ValueString(),
			})
		} else if !state.ReleasedScopeIds.IsNull() {
			tflog.Info(ctx, "gate open, releasing scope_ids", map[string]any{
				"count": len(state.ReleasedScopeIds.Elements()),
			})
		}
	} else {
		tflog.Info(ctx, "gate closed, returning null scope values")
	}

	// Save state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, "finished datasource read", map[string]any{
		"datasource": DataSourceName,
	})
}

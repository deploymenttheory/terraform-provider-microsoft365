package utilityDeploymentScheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// handleManualOverride processes manual override and sets state accordingly
func handleManualOverride(ctx context.Context, state *DeploymentSchedulerDataSourceModel, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Manual override enabled - forcing gate open")

	// Release scope ID(s) immediately
	setReleasedScopeIDs(state, true)
	if !state.ReleasedScopeId.IsNull() {
		tflog.Info(ctx, fmt.Sprintf("Manual override: releasing scope_id: %s", state.ReleasedScopeId.ValueString()))
	} else if !state.ReleasedScopeIds.IsNull() {
		tflog.Info(ctx, fmt.Sprintf("Manual override: releasing %d scope IDs", len(state.ReleasedScopeIds.Elements())))
	}

	state.ConditionMet = types.BoolValue(true)

	// Build structured status message
	msgBuilder := newStatusMessageBuilder(true)
	msgBuilder.setManualOverride()
	state.StatusMessage = types.StringValue(msgBuilder.build())

	setStateID(state)

	// Set null conditions_detail since we're bypassing evaluation
	state.ConditionsDetail = types.ObjectNull(
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

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// parseDeploymentStartTime parses or initializes the deployment start time
func parseDeploymentStartTime(ctx context.Context, state *DeploymentSchedulerDataSourceModel, currentTime time.Time, currentTimeStr string) (time.Time, string, error) {
	if !state.DeploymentStartTime.IsNull() && !state.DeploymentStartTime.IsUnknown() {
		deploymentStartTimeStr := state.DeploymentStartTime.ValueString()
		deploymentStartTime, err := time.Parse(time.RFC3339, deploymentStartTimeStr)
		if err != nil {
			return time.Time{}, "", fmt.Errorf("could not parse deployment_start_time '%s'. Must be in RFC3339 format (e.g., '2024-01-15T00:00:00Z'): %s", deploymentStartTimeStr, err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Using provided deployment start time: %s", deploymentStartTimeStr))
		return deploymentStartTime, deploymentStartTimeStr, nil
	}

	// No deployment start time provided - use current time
	tflog.Debug(ctx, fmt.Sprintf("No deployment_start_time provided, using current time: %s", currentTimeStr))
	return currentTime, currentTimeStr, nil
}

// evaluateTimeCondition evaluates the time-based condition and returns whether it's met, detail string, and metadata
func evaluateTimeCondition(ctx context.Context, state *DeploymentSchedulerDataSourceModel, currentTime time.Time, deploymentStartTime time.Time) (met bool, exists bool, delayHours int64, elapsed float64, detail string, err error) {
	if state.TimeCondition.IsNull() || state.TimeCondition.IsUnknown() {
		// No time condition = immediate release
		tflog.Debug(ctx, "No time condition specified, immediate release")
		return true, false, 0, 0, "", nil
	}

	var timeCondition TimeConditionModel
	diags := state.TimeCondition.As(ctx, &timeCondition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, false, 0, 0, "", fmt.Errorf("failed to parse time_condition: %v", diags)
	}

	delayStartTimeBy := timeCondition.DelayStartTimeBy.ValueInt64()
	if delayStartTimeBy < 0 {
		return false, false, 0, 0, "", fmt.Errorf("delay_start_time_by must be >= 0, got %d", delayStartTimeBy)
	}

	// Calculate hours elapsed since deployment start time
	duration := currentTime.Sub(deploymentStartTime)
	hoursElapsed := duration.Hours()

	// Check if delay has elapsed
	delayMet := hoursElapsed >= float64(delayStartTimeBy)

	// Check absolute_earliest constraint
	absoluteEarliestMet := true
	if !timeCondition.AbsoluteEarliest.IsNull() && !timeCondition.AbsoluteEarliest.IsUnknown() {
		absoluteEarliestTime, err := time.Parse(time.RFC3339, timeCondition.AbsoluteEarliest.ValueString())
		if err != nil {
			return false, false, 0, 0, "", fmt.Errorf("could not parse absolute_earliest '%s'. Must be in RFC3339 format: %s", timeCondition.AbsoluteEarliest.ValueString(), err)
		}
		absoluteEarliestMet = currentTime.After(absoluteEarliestTime) || currentTime.Equal(absoluteEarliestTime)
	}

	// Check absolute_latest constraint (gate permanently closes after this time)
	absoluteLatestExceeded := false
	if !timeCondition.AbsoluteLatest.IsNull() && !timeCondition.AbsoluteLatest.IsUnknown() {
		absoluteLatestTime, err := time.Parse(time.RFC3339, timeCondition.AbsoluteLatest.ValueString())
		if err != nil {
			return false, false, 0, 0, "", fmt.Errorf("could not parse absolute_latest '%s'. Must be in RFC3339 format: %s", timeCondition.AbsoluteLatest.ValueString(), err)
		}
		absoluteLatestExceeded = currentTime.After(absoluteLatestTime)
	}

	// Check max_open_duration_hours (gate closes after being open for this long)
	maxDurationExceeded := false
	if !timeCondition.MaxOpenDurationHours.IsNull() && !timeCondition.MaxOpenDurationHours.IsUnknown() {
		maxOpenDurationHours := timeCondition.MaxOpenDurationHours.ValueInt64()
		if maxOpenDurationHours > 0 {
			hoursOpen := hoursElapsed - float64(delayStartTimeBy)
			if hoursOpen > float64(maxOpenDurationHours) {
				maxDurationExceeded = true
			}
		}
	}

	// Time condition is met if: delay elapsed AND earliest met AND latest not exceeded AND max duration not exceeded
	timeConditionMet := delayMet && absoluteEarliestMet && !absoluteLatestExceeded && !maxDurationExceeded

	// Build detail message
	var timeDetail string
	if !timeConditionMet {
		if !delayMet {
			hoursRemaining := float64(delayStartTimeBy) - hoursElapsed
			timeDetail = fmt.Sprintf("Delay not elapsed (%.1fh / %.0fh required, %.1fh remaining)", hoursElapsed, float64(delayStartTimeBy), hoursRemaining)
		} else if !absoluteEarliestMet {
			timeDetail = "Before absolute_earliest time"
		} else if absoluteLatestExceeded {
			timeDetail = "After absolute_latest time (gate closed)"
		} else if maxDurationExceeded {
			timeDetail = "Max open duration exceeded (gate closed)"
		}
	} else {
		timeDetail = fmt.Sprintf("Delay elapsed (%.1fh / %.0fh required)", hoursElapsed, float64(delayStartTimeBy))
	}

	tflog.Debug(ctx, fmt.Sprintf("Time condition evaluation: %s, condition_met=%t", timeDetail, timeConditionMet))

	return timeConditionMet, true, delayStartTimeBy, hoursElapsed, timeDetail, nil
}

// evaluateDependencyGate evaluates whether the dependency gate condition is met
func evaluateDependencyGate(ctx context.Context, state *DeploymentSchedulerDataSourceModel, currentTime time.Time, deploymentStartTime time.Time) (met bool, message string, err error) {
	if state.DependsOnScheduler.IsNull() || state.DependsOnScheduler.IsUnknown() {
		return true, "", nil
	}

	var dependencyModel DependsOnSchedulerModel
	diags := state.DependsOnScheduler.As(ctx, &dependencyModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse depends_on_scheduler: %v", diags)
	}

	prerequisiteDelayStartTimeBy := dependencyModel.PrerequisiteDelayStartTimeBy.ValueInt64()
	minimumOpenHours := dependencyModel.MinimumOpenHours.ValueInt64()

	if minimumOpenHours < 0 {
		return false, "", fmt.Errorf("minimum_open_hours must be >= 0, got %d", minimumOpenHours)
	}

	// Calculate when the prerequisite scheduler opens (deployment start + prerequisite delay)
	prerequisiteOpenTime := float64(prerequisiteDelayStartTimeBy)

	// Calculate current elapsed time
	currentHoursElapsed := currentTime.Sub(deploymentStartTime).Hours()

	// Calculate how long the prerequisite has been open
	prerequisiteHoursOpen := currentHoursElapsed - prerequisiteOpenTime
	if prerequisiteHoursOpen < 0 {
		prerequisiteHoursOpen = 0
	}

	// Dependency gate is met if prerequisite has been open for minimum required hours
	dependencyGateMet := prerequisiteHoursOpen >= float64(minimumOpenHours)

	tflog.Debug(ctx, fmt.Sprintf(
		"Dependency gate evaluation: prerequisite_offset=%dh, current_elapsed=%.1fh, prerequisite_open=%.1fh, minimum_required=%dh, met=%t",
		prerequisiteDelayStartTimeBy, currentHoursElapsed, prerequisiteHoursOpen, minimumOpenHours, dependencyGateMet))

	var dependencyMessage string
	if dependencyGateMet {
		dependencyMessage = fmt.Sprintf("Prerequisite open for %.1fh / %.0fh required)", prerequisiteHoursOpen, float64(minimumOpenHours))
	} else {
		if prerequisiteHoursOpen > 0 {
			dependencyMessage = fmt.Sprintf("Prerequisite open for %.1fh / %.0fh required)", prerequisiteHoursOpen, float64(minimumOpenHours))
		} else {
			hoursUntilPrerequisiteOpens := prerequisiteOpenTime - currentHoursElapsed
			dependencyMessage = fmt.Sprintf("Prerequisite hasn't opened yet (%.1fh remaining)", hoursUntilPrerequisiteOpens)
		}
	}

	return dependencyGateMet, dependencyMessage, nil
}

package utilityDeploymentScheduler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// handleManualOverride processes manual override and sets state accordingly
func handleManualOverride(ctx context.Context, state *DeploymentSchedulerDataSourceModel, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "manual override enabled, forcing gate open")

	// Release scope ID(s) immediately
	setReleasedScopeIDs(state, true)
	if !state.ReleasedScopeId.IsNull() {
		tflog.Info(ctx, "manual override releasing scope_id", map[string]any{
			"scope_id": state.ReleasedScopeId.ValueString(),
		})
	} else if !state.ReleasedScopeIds.IsNull() {
		tflog.Info(ctx, "manual override releasing scope_ids", map[string]any{
			"count": len(state.ReleasedScopeIds.Elements()),
		})
	}

	state.ConditionMet = types.BoolValue(true)

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
			return time.Time{}, "", fmt.Errorf("failed to parse `deployment_start_time`: %w", err)
		}
		tflog.Debug(ctx, "using provided deployment start time", map[string]any{
			"deployment_start_time": deploymentStartTimeStr,
		})
		return deploymentStartTime, deploymentStartTimeStr, nil
	}

	tflog.Debug(ctx, "no deployment_start_time provided, using current time", map[string]any{
		"current_time": currentTimeStr,
	})
	return currentTime, currentTimeStr, nil
}

// evaluateTimeCondition evaluates the time-based condition and returns whether it's met, detail string, and metadata
func evaluateTimeCondition(ctx context.Context, state *DeploymentSchedulerDataSourceModel, currentTime time.Time, deploymentStartTime time.Time) (met bool, exists bool, delayHours int64, elapsed float64, detail string, err error) {
	if state.TimeCondition.IsNull() || state.TimeCondition.IsUnknown() {
		// No time condition = immediate release
		tflog.Debug(ctx, "no time_condition specified, immediate release")
		return true, false, 0, 0, "", nil
	}

	var timeCondition TimeConditionModel
	diags := state.TimeCondition.As(ctx, &timeCondition, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, false, 0, 0, "", fmt.Errorf("failed to parse `time_condition`: %v", diags)
	}

	// Schema validation ensures delay_start_time_by >= 0
	delayStartTimeBy := timeCondition.DelayStartTimeBy.ValueInt64()

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
			return false, false, 0, 0, "", fmt.Errorf("failed to parse `absolute_earliest`: %w", err)
		}
		absoluteEarliestMet = currentTime.After(absoluteEarliestTime) || currentTime.Equal(absoluteEarliestTime)
	}

	// Check absolute_latest constraint (gate permanently closes after this time)
	absoluteLatestExceeded := false
	if !timeCondition.AbsoluteLatest.IsNull() && !timeCondition.AbsoluteLatest.IsUnknown() {
		absoluteLatestTime, err := time.Parse(time.RFC3339, timeCondition.AbsoluteLatest.ValueString())
		if err != nil {
			return false, false, 0, 0, "", fmt.Errorf("failed to parse `absolute_latest`: %w", err)
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

	tflog.Debug(ctx, "time_condition evaluated", map[string]any{
		"condition_met":  timeConditionMet,
		"hours_elapsed":  hoursElapsed,
		"delay_required": delayStartTimeBy,
		"detail":         timeDetail,
	})

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
		return false, "", fmt.Errorf("failed to parse `depends_on_scheduler`: %v", diags)
	}

	// Schema validation ensures these values >= 0
	prerequisiteDelayStartTimeBy := dependencyModel.PrerequisiteDelayStartTimeBy.ValueInt64()
	minimumOpenHours := dependencyModel.MinimumOpenHours.ValueInt64()

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

	tflog.Debug(ctx, "depends_on_scheduler evaluated", map[string]any{
		"condition_met":                    dependencyGateMet,
		"prerequisite_delay_start_time_by": prerequisiteDelayStartTimeBy,
		"current_hours_elapsed":            currentHoursElapsed,
		"prerequisite_hours_open":          prerequisiteHoursOpen,
		"minimum_open_hours":               minimumOpenHours,
	})

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

// isWithinTimeWindow checks if the given time is within the specified window
func isWithinTimeWindow(ctx context.Context, checkTime time.Time, window TimeWindowModel) (bool, error) {
	// Check days of week if specified
	if !window.DaysOfWeek.IsNull() && !window.DaysOfWeek.IsUnknown() {
		var daysOfWeek []string
		diags := window.DaysOfWeek.ElementsAs(ctx, &daysOfWeek, false)
		if diags.HasError() {
			return false, fmt.Errorf("failed to parse `days_of_week`")
		}

		if len(daysOfWeek) > 0 {
			currentDay := checkTime.Weekday().String()
			dayMatched := false
			for _, day := range daysOfWeek {
				if strings.EqualFold(day, currentDay) {
					dayMatched = true
					break
				}
			}
			if !dayMatched {
				return false, nil
			}
		}
	}

	// Check absolute date range if specified
	if !window.DateStart.IsNull() && !window.DateStart.IsUnknown() {
		dateStart, err := time.Parse(time.RFC3339, window.DateStart.ValueString())
		if err != nil {
			return false, fmt.Errorf("failed to parse `date_start`: %w", err)
		}
		if checkTime.Before(dateStart) {
			return false, nil
		}
	}

	if !window.DateEnd.IsNull() && !window.DateEnd.IsUnknown() {
		dateEnd, err := time.Parse(time.RFC3339, window.DateEnd.ValueString())
		if err != nil {
			return false, fmt.Errorf("failed to parse `date_end`: %w", err)
		}
		if checkTime.After(dateEnd) {
			return false, nil
		}
	}

	// Check time of day range if specified
	hasTimeOfDayStart := !window.TimeOfDayStart.IsNull() && !window.TimeOfDayStart.IsUnknown()
	hasTimeOfDayEnd := !window.TimeOfDayEnd.IsNull() && !window.TimeOfDayEnd.IsUnknown()

	if hasTimeOfDayStart || hasTimeOfDayEnd {
		// Extract time of day from checkTime
		currentTimeOfDay := checkTime.Format("15:04:05")

		// Default to full day if not specified
		startTimeStr := "00:00:00"
		if hasTimeOfDayStart {
			startTimeStr = window.TimeOfDayStart.ValueString()
		}

		endTimeStr := "23:59:59"
		if hasTimeOfDayEnd {
			endTimeStr = window.TimeOfDayEnd.ValueString()
		}

		// Validate time format
		if _, err := time.Parse("15:04:05", startTimeStr); err != nil {
			return false, fmt.Errorf("failed to parse `time_of_day_start`: %w", err)
		}
		if _, err := time.Parse("15:04:05", endTimeStr); err != nil {
			return false, fmt.Errorf("failed to parse `time_of_day_end`: %w", err)
		}

		// Compare times as strings (works for HH:MM:SS format)
		if currentTimeOfDay < startTimeStr || currentTimeOfDay > endTimeStr {
			return false, nil
		}
	}

	// All checks passed
	return true, nil
}

// evaluateInclusionWindows checks if current time is within any inclusion window
func evaluateInclusionWindows(ctx context.Context, currentTime time.Time, inclusionTimeWindows types.Object) (bool, string, error) {
	if inclusionTimeWindows.IsNull() || inclusionTimeWindows.IsUnknown() {
		return true, "", nil
	}

	var inclusionModel InclusionTimeWindowsModel
	diags := inclusionTimeWindows.As(ctx, &inclusionModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse `inclusion_time_windows`")
	}

	var windows []TimeWindowModel
	diags = inclusionModel.Window.ElementsAs(ctx, &windows, false)
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse `inclusion_time_windows.window`")
	}

	if len(windows) == 0 {
		return true, "", nil
	}

	// Check if current time is within ANY window (OR logic)
	for i, window := range windows {
		isWithin, err := isWithinTimeWindow(ctx, currentTime, window)
		if err != nil {
			return false, "", fmt.Errorf("inclusion_time_windows[%d]: %w", i, err)
		}
		if isWithin {
			return true, fmt.Sprintf("within inclusion window %d", i+1), nil
		}
	}

	return false, "outside all inclusion windows", nil
}

// evaluateExclusionWindows checks if current time is within any exclusion window
func evaluateExclusionWindows(ctx context.Context, currentTime time.Time, exclusionTimeWindows types.Object) (bool, string, error) {
	if exclusionTimeWindows.IsNull() || exclusionTimeWindows.IsUnknown() {
		return false, "", nil
	}

	var exclusionModel ExclusionTimeWindowsModel
	diags := exclusionTimeWindows.As(ctx, &exclusionModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse `exclusion_time_windows`")
	}

	var windows []TimeWindowModel
	diags = exclusionModel.Window.ElementsAs(ctx, &windows, false)
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse `exclusion_time_windows.window`")
	}

	if len(windows) == 0 {
		return false, "", nil
	}

	// Check if current time is within ANY window (OR logic)
	for i, window := range windows {
		isWithin, err := isWithinTimeWindow(ctx, currentTime, window)
		if err != nil {
			return false, "", fmt.Errorf("exclusion_time_windows[%d]: %w", i, err)
		}
		if isWithin {
			return true, fmt.Sprintf("within exclusion window %d", i+1), nil
		}
	}

	return false, "outside all exclusion windows", nil
}

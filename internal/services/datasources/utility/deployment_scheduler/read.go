package utilityDeploymentScheduler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Deployment Scheduler data source.
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

	// Initialize time
	currentTime := time.Now().UTC()
	currentTimeStr := currentTime.Format(time.RFC3339)

	// Parse or initialize deployment start time
	deploymentStartTime, deploymentStartTimeStr, err := parseDeploymentStartTime(ctx, &state, currentTime, currentTimeStr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Deployment Start Time", err.Error())
		return
	}
	state.DeploymentStartTime = types.StringValue(deploymentStartTimeStr)

	// Set defaults for optional fields
	setStateDefaults(&state)

	// Handle manual override - bypasses all condition evaluation
	if state.ManualOverride.ValueBool() {
		handleManualOverride(ctx, &state, resp)
		return
	}

	// Evaluate time condition
	timeConditionMet, timeConditionExists, delayStartTimeBy, hoursElapsed, timeConditionDetail, err := evaluateTimeCondition(ctx, &state, currentTime, deploymentStartTime)
	if err != nil {
		resp.Diagnostics.AddError("Time Condition Evaluation Failed", err.Error())
		return
	}

	// Evaluate inclusion windows
	inclusionWindowMet, inclusionMessage, err := evaluateInclusionWindows(ctx, currentTime, state.InclusionTimeWindows)
	if err != nil {
		resp.Diagnostics.AddError("Inclusion Window Evaluation Failed", err.Error())
		return
	}

	// Evaluate exclusion windows
	exclusionWindowActive, exclusionMessage, err := evaluateExclusionWindows(ctx, currentTime, state.ExclusionTimeWindows)
	if err != nil {
		resp.Diagnostics.AddError("Exclusion Window Evaluation Failed", err.Error())
		return
	}

	// Evaluate dependency gate
	dependencyGateMet, dependencyMessage, err := evaluateDependencyGate(ctx, &state, currentTime, deploymentStartTime)
	if err != nil {
		resp.Diagnostics.AddError("Dependency Gate Evaluation Failed", err.Error())
		return
	}

	// Determine overall condition status
	allConditionsMet := timeConditionMet && inclusionWindowMet && !exclusionWindowActive && dependencyGateMet

	// Build structured status message
	msgBuilder := newStatusMessageBuilder(allConditionsMet)
	if timeConditionExists {
		msgBuilder.addTimeCondition(timeConditionMet, timeConditionDetail)
	}
	if inclusionMessage != "" {
		msgBuilder.addInclusionWindow(inclusionWindowMet, inclusionMessage)
	}
	if exclusionMessage != "" {
		msgBuilder.addExclusionWindow(exclusionWindowActive, exclusionMessage)
	}
	if dependencyMessage != "" {
		msgBuilder.setDependency(dependencyGateMet, dependencyMessage)
	}
	state.StatusMessage = types.StringValue(msgBuilder.build())

	// Set condition_met
	state.ConditionMet = types.BoolValue(allConditionsMet)

	// Set released_scope_id or released_scope_ids based on condition status
	setReleasedScopeIDs(&state, allConditionsMet)

	// Log the result
	if allConditionsMet {
		if !state.ReleasedScopeId.IsNull() {
			tflog.Info(ctx, fmt.Sprintf("Gate open: releasing scope_id: %s", state.ReleasedScopeId.ValueString()))
		} else if !state.ReleasedScopeIds.IsNull() {
			tflog.Info(ctx, fmt.Sprintf("Gate open: releasing %d scope IDs", len(state.ReleasedScopeIds.Elements())))
		}
	} else {
		tflog.Info(ctx, "Gate closed: returning null scope IDs")
	}

	// Build conditions_detail
	if timeConditionExists {
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
			return
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
		if resp.Diagnostics.HasError() {
			return
		}
		state.ConditionsDetail = conditionsDetailObj
	} else {
		// No time condition - set null detail
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
	}

	// Set ID
	setStateID(&state)

	// Save state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}

// isWithinTimeWindow checks if the given time is within the specified window
func isWithinTimeWindow(ctx context.Context, checkTime time.Time, window TimeWindowModel) (bool, error) {
	// Check days of week if specified
	if !window.DaysOfWeek.IsNull() && !window.DaysOfWeek.IsUnknown() {
		var daysOfWeek []string
		diags := window.DaysOfWeek.ElementsAs(ctx, &daysOfWeek, false)
		if diags.HasError() {
			return false, fmt.Errorf("failed to parse days_of_week")
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
			return false, fmt.Errorf("invalid date_start format: %w", err)
		}
		if checkTime.Before(dateStart) {
			return false, nil
		}
	}

	if !window.DateEnd.IsNull() && !window.DateEnd.IsUnknown() {
		dateEnd, err := time.Parse(time.RFC3339, window.DateEnd.ValueString())
		if err != nil {
			return false, fmt.Errorf("invalid date_end format: %w", err)
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
			return false, fmt.Errorf("invalid time_of_day_start format (must be HH:MM:SS): %w", err)
		}
		if _, err := time.Parse("15:04:05", endTimeStr); err != nil {
			return false, fmt.Errorf("invalid time_of_day_end format (must be HH:MM:SS): %w", err)
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
// Returns true if no windows defined, or if current time matches at least one window
func evaluateInclusionWindows(ctx context.Context, currentTime time.Time, inclusionTimeWindows types.Object) (bool, string, error) {
	if inclusionTimeWindows.IsNull() || inclusionTimeWindows.IsUnknown() {
		return true, "", nil // No inclusion windows = always allowed
	}

	var inclusionModel InclusionTimeWindowsModel
	diags := inclusionTimeWindows.As(ctx, &inclusionModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse inclusion_time_windows")
	}

	var windows []TimeWindowModel
	diags = inclusionModel.Window.ElementsAs(ctx, &windows, false)
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse inclusion time windows")
	}

	if len(windows) == 0 {
		return true, "", nil
	}

	// Check if current time is within ANY window (OR logic)
	for i, window := range windows {
		isWithin, err := isWithinTimeWindow(ctx, currentTime, window)
		if err != nil {
			return false, "", fmt.Errorf("window %d: %w", i, err)
		}
		if isWithin {
			return true, fmt.Sprintf("within inclusion window %d", i+1), nil
		}
	}

	return false, "outside all inclusion windows", nil
}

// evaluateExclusionWindows checks if current time is within any exclusion window
// Returns true if current time matches ANY exclusion window (deployment should be blocked)
func evaluateExclusionWindows(ctx context.Context, currentTime time.Time, exclusionTimeWindows types.Object) (bool, string, error) {
	if exclusionTimeWindows.IsNull() || exclusionTimeWindows.IsUnknown() {
		return false, "", nil // No exclusion windows = not blocked
	}

	var exclusionModel ExclusionTimeWindowsModel
	diags := exclusionTimeWindows.As(ctx, &exclusionModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse exclusion_time_windows")
	}

	var windows []TimeWindowModel
	diags = exclusionModel.Window.ElementsAs(ctx, &windows, false)
	if diags.HasError() {
		return false, "", fmt.Errorf("failed to parse exclusion time windows")
	}

	if len(windows) == 0 {
		return false, "", nil
	}

	// Check if current time is within ANY window (OR logic)
	for i, window := range windows {
		isWithin, err := isWithinTimeWindow(ctx, currentTime, window)
		if err != nil {
			return false, "", fmt.Errorf("window %d: %w", i, err)
		}
		if isWithin {
			return true, fmt.Sprintf("within exclusion window %d", i+1), nil
		}
	}

	return false, "outside all exclusion windows", nil
}

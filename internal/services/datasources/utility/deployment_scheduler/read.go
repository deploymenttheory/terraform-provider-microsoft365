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

	// Variable declarations
	var (
		// Time variables
		currentTime            = time.Now().UTC()
		currentTimeStr         = currentTime.Format(time.RFC3339)
		deploymentStartTime    time.Time
		deploymentStartTimeStr string

		// Time condition variables
		timeConditionMet    bool
		timeConditionExists bool
		delayStartTimeBy    int64
		hoursElapsed        float64
		timeConditionDetail string

		// Inclusion/Exclusion window variables
		inclusionWindowMet    = true
		inclusionMessage      string
		exclusionWindowActive bool
		exclusionMessage      string

		// Dependency variables
		dependencyGateMet = true
		dependencyMessage string

		// Message builder
		msgBuilder *StatusMessageBuilder
	)

	if !state.DeploymentStartTime.IsNull() && !state.DeploymentStartTime.IsUnknown() {
		// Use provided deployment start time
		deploymentStartTimeStr = state.DeploymentStartTime.ValueString()
		var err error
		deploymentStartTime, err = time.Parse(time.RFC3339, deploymentStartTimeStr)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Deployment Start Time",
				fmt.Sprintf("Could not parse deployment_start_time '%s'. Must be in RFC3339 format (e.g., '2024-01-15T00:00:00Z'): %s", deploymentStartTimeStr, err),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Using provided deployment start time: %s", deploymentStartTimeStr))
	} else {
		// No deployment start time provided - use current time on each evaluation
		// Note: For time-based conditions to work correctly across multiple applies,
		// deployment_start_time should be explicitly provided in configuration
		deploymentStartTime = currentTime
		deploymentStartTimeStr = currentTimeStr
		tflog.Debug(ctx, fmt.Sprintf("No deployment_start_time provided, using current time: %s", deploymentStartTimeStr))
	}
	state.DeploymentStartTime = types.StringValue(deploymentStartTimeStr)

	// Set defaults for optional fields
	setStateDefaults(&state)

	// Check manual override first - bypasses all other conditions
	if state.ManualOverride.ValueBool() {
		tflog.Info(ctx, "Manual override enabled - forcing gate open")

		// Release scope ID(s) immediately
		setReleasedScopeIDs(&state, true)
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

		setStateID(&state)

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

		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	// Evaluate time condition
	if !state.TimeCondition.IsNull() && !state.TimeCondition.IsUnknown() {
		timeConditionExists = true

		var timeCondition TimeConditionModel
		diags = state.TimeCondition.As(ctx, &timeCondition, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		delayStartTimeBy = timeCondition.DelayStartTimeBy.ValueInt64()
		if delayStartTimeBy < 0 {
			resp.Diagnostics.AddError(
				"Invalid Time Condition",
				fmt.Sprintf("delay_start_time_by must be >= 0, got %d", delayStartTimeBy),
			)
			return
		}

		// Calculate hours elapsed since deployment start time
		duration := currentTime.Sub(deploymentStartTime)
		hoursElapsed = duration.Hours()

		// Check if delay has elapsed
		delayMet := hoursElapsed >= float64(delayStartTimeBy)

		// Check absolute_earliest constraint
		absoluteEarliestMet := true
		if !timeCondition.AbsoluteEarliest.IsNull() && !timeCondition.AbsoluteEarliest.IsUnknown() {
			absoluteEarliestTime, err := time.Parse(time.RFC3339, timeCondition.AbsoluteEarliest.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid absolute_earliest Time",
					fmt.Sprintf("Could not parse absolute_earliest '%s'. Must be in RFC3339 format: %s", timeCondition.AbsoluteEarliest.ValueString(), err),
				)
				return
			}
			absoluteEarliestMet = currentTime.After(absoluteEarliestTime) || currentTime.Equal(absoluteEarliestTime)
		}

		// Check absolute_latest constraint (gate permanently closes after this time)
		absoluteLatestExceeded := false
		if !timeCondition.AbsoluteLatest.IsNull() && !timeCondition.AbsoluteLatest.IsUnknown() {
			absoluteLatestTime, err := time.Parse(time.RFC3339, timeCondition.AbsoluteLatest.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid absolute_latest Time",
					fmt.Sprintf("Could not parse absolute_latest '%s'. Must be in RFC3339 format: %s", timeCondition.AbsoluteLatest.ValueString(), err),
				)
				return
			}
			absoluteLatestExceeded = currentTime.After(absoluteLatestTime)
		}

		// Check max_open_duration_hours (gate auto-closes after being open for specified duration)
		maxDurationExceeded := false
		if !timeCondition.MaxOpenDurationHours.IsNull() && !timeCondition.MaxOpenDurationHours.IsUnknown() {
			maxDuration := timeCondition.MaxOpenDurationHours.ValueInt64()
			if maxDuration < 0 {
				resp.Diagnostics.AddError(
					"Invalid Time Condition",
					fmt.Sprintf("max_open_duration_hours must be >= 0, got %d", maxDuration),
				)
				return
			}
			if maxDuration > 0 {
				// Calculate when gate first opened (deployment_start_time + delay)
				gateOpenTime := deploymentStartTime.Add(time.Duration(delayStartTimeBy) * time.Hour)
				hoursOpen := currentTime.Sub(gateOpenTime).Hours()

				if hoursOpen > float64(maxDuration) {
					maxDurationExceeded = true
				}
			}
		}

		// Combine all time condition checks
		timeConditionMet = delayMet && absoluteEarliestMet && !absoluteLatestExceeded && !maxDurationExceeded

		tflog.Debug(ctx, fmt.Sprintf("Time condition evaluation: delay_hours=%d, hours_elapsed=%.2f, delay_met=%t, absolute_earliest_met=%t, absolute_latest_exceeded=%t, max_duration_exceeded=%t, overall_met=%t",
			delayStartTimeBy, hoursElapsed, delayMet, absoluteEarliestMet, absoluteLatestExceeded, maxDurationExceeded, timeConditionMet))

		// Build time condition detail message
		var timeDetail string
		if absoluteLatestExceeded {
			timeDetail = "Exceeded absolute_latest deadline - gate permanently closed"
		} else if maxDurationExceeded {
			gateOpenTime := deploymentStartTime.Add(time.Duration(delayStartTimeBy) * time.Hour)
			hoursOpen := currentTime.Sub(gateOpenTime).Hours()
			maxDuration := timeCondition.MaxOpenDurationHours.ValueInt64()
			timeDetail = fmt.Sprintf("Max open duration exceeded (open %.1fh / %.0fh max)", hoursOpen, float64(maxDuration))
		} else if !absoluteEarliestMet {
			absoluteEarliestTime, _ := time.Parse(time.RFC3339, timeCondition.AbsoluteEarliest.ValueString())
			timeDetail = fmt.Sprintf("Before absolute_earliest (%s)", absoluteEarliestTime.Format("2006-01-02 15:04 MST"))
		} else if !delayMet {
			timeDetail = fmt.Sprintf("Delay not elapsed (%.1fh / %.0fh required)", hoursElapsed, float64(delayStartTimeBy))
		} else {
			timeDetail = fmt.Sprintf("Delay elapsed (%.1fh / %.0fh required)", hoursElapsed, float64(delayStartTimeBy))
		}

		// Store for message builder (will be added after all conditions evaluated)
		timeConditionDetail = timeDetail
	} else {
		// No time condition = immediate release
		timeConditionMet = true
		timeConditionExists = false
		tflog.Debug(ctx, "No time condition specified, immediate release")
	}

	// Evaluate inclusion time windows
	if !state.InclusionTimeWindows.IsNull() && !state.InclusionTimeWindows.IsUnknown() {
		var err error
		inclusionWindowMet, inclusionMessage, err = evaluateInclusionWindows(ctx, currentTime, state.InclusionTimeWindows)
		if err != nil {
			resp.Diagnostics.AddError(
				"Inclusion Time Window Evaluation Failed",
				fmt.Sprintf("Could not evaluate inclusion_time_windows: %s", err),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Inclusion time window evaluation: %s, condition_met=%t", inclusionMessage, inclusionWindowMet))
	} else {
		inclusionWindowMet = true // No inclusion window = always allowed
	}

	// Evaluate exclusion time windows
	if !state.ExclusionTimeWindows.IsNull() && !state.ExclusionTimeWindows.IsUnknown() {
		var err error
		exclusionWindowActive, exclusionMessage, err = evaluateExclusionWindows(ctx, currentTime, state.ExclusionTimeWindows)
		if err != nil {
			resp.Diagnostics.AddError(
				"Exclusion Time Window Evaluation Failed",
				fmt.Sprintf("Could not evaluate exclusion_time_windows: %s", err),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Exclusion time window evaluation: %s, blocked=%t", exclusionMessage, exclusionWindowActive))
	} else {
		exclusionWindowActive = false // No exclusion window = not blocked
	}

	// Evaluate dependency gate
	if !state.DependsOnScheduler.IsNull() && !state.DependsOnScheduler.IsUnknown() {
		var dependencyModel DependsOnSchedulerModel
		diags = state.DependsOnScheduler.As(ctx, &dependencyModel, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		prerequisiteDelayStartTimeBy := dependencyModel.PrerequisiteDelayStartTimeBy.ValueInt64()
		minimumOpenHours := dependencyModel.MinimumOpenHours.ValueInt64()

		if minimumOpenHours < 0 {
			resp.Diagnostics.AddError(
				"Invalid Dependency Configuration",
				fmt.Sprintf("minimum_open_hours must be >= 0, got %d", minimumOpenHours),
			)
			return
		}

		// Calculate when prerequisite gate opened
		duration := currentTime.Sub(deploymentStartTime)
		currentHoursElapsed := duration.Hours()

		// Calculate when prerequisite would have opened
		prerequisiteOpenTime := float64(prerequisiteDelayStartTimeBy)

		// Calculate how long prerequisite has been open
		var prerequisiteHoursOpen float64
		if currentHoursElapsed >= prerequisiteOpenTime {
			prerequisiteHoursOpen = currentHoursElapsed - prerequisiteOpenTime
		} else {
			prerequisiteHoursOpen = 0
		}

		// Check if prerequisite has been open long enough
		dependencyGateMet = prerequisiteHoursOpen >= float64(minimumOpenHours)

		tflog.Debug(ctx, fmt.Sprintf("Dependency gate evaluation: prerequisite_delay=%dh, current_elapsed=%.2fh, prerequisite_open_for=%.2fh, minimum_required=%dh, met=%t",
			prerequisiteDelayStartTimeBy, currentHoursElapsed, prerequisiteHoursOpen, minimumOpenHours, dependencyGateMet))

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
	}

	// Determine overall condition status
	// All conditions must pass AND exclusion must not be active AND dependency must be met
	allConditionsMet := timeConditionMet && inclusionWindowMet && !exclusionWindowActive && dependencyGateMet

	// Build structured status message using message builder
	msgBuilder = newStatusMessageBuilder(allConditionsMet)

	// Add time condition if present
	if timeConditionExists {
		msgBuilder.addTimeCondition(timeConditionMet, timeConditionDetail)
	}

	// Add inclusion window if present
	if inclusionMessage != "" {
		msgBuilder.addInclusionWindow(inclusionWindowMet, inclusionMessage)
	}

	// Add exclusion window if present
	if exclusionMessage != "" {
		msgBuilder.addExclusionWindow(exclusionWindowActive, exclusionMessage)
	}

	// Add dependency if present
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

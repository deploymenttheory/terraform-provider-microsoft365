package utilityDeploymentScheduler

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentSchedulerDataSourceModel struct {
	Id                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	DeploymentStartTime  types.String   `tfsdk:"deployment_start_time"`
	ScopeId              types.String   `tfsdk:"scope_id"`
	ScopeIds             types.List     `tfsdk:"scope_ids"`
	TimeCondition        types.Object   `tfsdk:"time_condition"`
	InclusionTimeWindows types.Object   `tfsdk:"inclusion_time_windows"`
	ExclusionTimeWindows types.Object   `tfsdk:"exclusion_time_windows"`
	ManualOverride       types.Bool     `tfsdk:"manual_override"`
	DependsOnScheduler   types.Object   `tfsdk:"depends_on_scheduler"`
	RequireAllConditions types.Bool     `tfsdk:"require_all_conditions"`
	ReleasedScopeId      types.String   `tfsdk:"released_scope_id"`
	ReleasedScopeIds     types.List     `tfsdk:"released_scope_ids"`
	ConditionMet         types.Bool     `tfsdk:"condition_met"`
	ConditionsDetail     types.Object   `tfsdk:"conditions_detail"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type TimeConditionModel struct {
	DelayStartTimeBy     types.Int64  `tfsdk:"delay_start_time_by"`
	AbsoluteEarliest     types.String `tfsdk:"absolute_earliest"`
	AbsoluteLatest       types.String `tfsdk:"absolute_latest"`
	MaxOpenDurationHours types.Int64  `tfsdk:"max_open_duration_hours"`
}

type InclusionTimeWindowsModel struct {
	Window types.List `tfsdk:"window"`
}

type ExclusionTimeWindowsModel struct {
	Window types.List `tfsdk:"window"`
}

type TimeWindowModel struct {
	DaysOfWeek     types.List   `tfsdk:"days_of_week"`
	TimeOfDayStart types.String `tfsdk:"time_of_day_start"`
	TimeOfDayEnd   types.String `tfsdk:"time_of_day_end"`
	DateStart      types.String `tfsdk:"date_start"`
	DateEnd        types.String `tfsdk:"date_end"`
}

type DependsOnSchedulerModel struct {
	PrerequisiteDelayStartTimeBy types.Int64 `tfsdk:"prerequisite_delay_start_time_by"`
	MinimumOpenHours             types.Int64 `tfsdk:"minimum_open_hours"`
}

type ConditionsDetailModel struct {
	TimeConditionDetail types.Object `tfsdk:"time_condition_detail"`
}

type TimeConditionDetailModel struct {
	Required             types.Bool    `tfsdk:"required"`
	DelayStartTimeBy     types.Int64   `tfsdk:"delay_start_time_by"`
	DeploymentStartTime  types.String  `tfsdk:"deployment_start_time"`
	CurrentTime          types.String  `tfsdk:"current_time"`
	HoursElapsed         types.Float64 `tfsdk:"hours_elapsed"`
	ConditionMet         types.Bool    `tfsdk:"condition_met"`
	AbsoluteEarliestTime types.String  `tfsdk:"absolute_earliest_time"`
	AbsoluteLatestTime   types.String  `tfsdk:"absolute_latest_time"`
	MaxOpenDurationHours types.Int64   `tfsdk:"max_open_duration_hours"`
	EarliestTimeMet      types.Bool    `tfsdk:"earliest_time_met"`
	LatestTimeMet        types.Bool    `tfsdk:"latest_time_met"`
	MaxDurationMet       types.Bool    `tfsdk:"max_duration_met"`
}

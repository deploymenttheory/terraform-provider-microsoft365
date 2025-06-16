// Base resource REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-devicehealthscriptassignment-list?view=graph-rest-beta
package graphBetaWindowsRemediationScriptAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceHealthScriptAssignmentResourceModel struct {
	DeviceHealthScriptId types.String                  `tfsdk:"device_health_script_id"`
	ID                   types.String                  `tfsdk:"id"`
	Target               AssignmentTargetResourceModel `tfsdk:"target"`
	RunRemediationScript types.Bool                    `tfsdk:"run_remediation_script"`
	RunSchedule          *RunScheduleResourceModel     `tfsdk:"run_schedule"`
	Timeouts             timeouts.Value                `tfsdk:"timeouts"`
}

// Target models
type AssignmentTargetResourceModel struct {
	TargetType                                 types.String `tfsdk:"target_type"` // allDevices, allLicensedUsers, configurationManagerCollection, exclusionGroupAssignment, groupAssignment
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	GroupId                                    types.String `tfsdk:"group_id"`
	CollectionId                               types.String `tfsdk:"collection_id"`
}

// Schedule models
type RunScheduleResourceModel struct {
	Daily  *DailyScheduleResourceModel   `tfsdk:"daily"`
	Hourly *HourlyScheduleResourceModel  `tfsdk:"hourly"`
	Once   *RunOnceScheduleResourceModel `tfsdk:"once"`
}

type DailyScheduleResourceModel struct {
	Interval types.Int32  `tfsdk:"interval"`
	UseUtc   types.Bool   `tfsdk:"use_utc"`
	Time     types.String `tfsdk:"time"`
}

type HourlyScheduleResourceModel struct {
	Interval types.Int32 `tfsdk:"interval"`
	UseUtc   types.Bool  `tfsdk:"use_utc"`
}

type RunOnceScheduleResourceModel struct {
	Interval types.Int32  `tfsdk:"interval"`
	Date     types.String `tfsdk:"date"`
	UseUtc   types.Bool   `tfsdk:"use_utc"`
	Time     types.String `tfsdk:"time"`
}

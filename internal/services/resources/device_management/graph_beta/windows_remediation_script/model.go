// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicehealthscript?view=graph-rest-beta
package graphBetaWindowsRemediationScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceHealthScriptResourceModel defines the schema for a Device Health Script.
type DeviceHealthScriptResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	DisplayName               types.String   `tfsdk:"display_name"`
	Description               types.String   `tfsdk:"description"`
	Publisher                 types.String   `tfsdk:"publisher"`
	RunAs32Bit                types.Bool     `tfsdk:"run_as_32_bit"`
	RunAsAccount              types.String   `tfsdk:"run_as_account"`
	EnforceSignatureCheck     types.Bool     `tfsdk:"enforce_signature_check"`
	DetectionScriptContent    types.String   `tfsdk:"detection_script_content"`
	RemediationScriptContent  types.String   `tfsdk:"remediation_script_content"`
	RoleScopeTagIds           types.Set      `tfsdk:"role_scope_tag_ids"`
	Version                   types.String   `tfsdk:"version"`
	IsGlobalScript            types.Bool     `tfsdk:"is_global_script"`
	DeviceHealthScriptType    types.String   `tfsdk:"device_health_script_type"`
	CreatedDateTime           types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime      types.String   `tfsdk:"last_modified_date_time"`
	HighestAvailableVersion   types.String   `tfsdk:"highest_available_version"`
	DetectionScriptParameters types.List     `tfsdk:"detection_script_parameters"`
	Assignments               types.Set      `tfsdk:"assignments"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
}

// DeviceHealthScriptParameterModel defines a single script parameter.
type DeviceHealthScriptParameterModel struct {
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	IsRequired                       types.Bool   `tfsdk:"is_required"`
	ApplyDefaultValueWhenNotAssigned types.Bool   `tfsdk:"apply_default_value_when_not_assigned"`
}

// WindowsRemediationScriptAssignmentModel defines the schema for a Windows Remediation Script assignment.
type WindowsRemediationScriptAssignmentModel struct {
	// Target assignment fields - only one should be used at a time
	Type    types.String `tfsdk:"type"`     // "allDevicesAssignmentTarget", "allLicensedUsersAssignmentTarget", "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)

	// Assignment filter fields
	FilterId   types.String `tfsdk:"filter_id"`
	FilterType types.String `tfsdk:"filter_type"` // "include", "exclude", or "none"

	// Schedule configuration - only one should be used at a time
	DailySchedule   *WindowsRemediationScriptDailyScheduleModel   `tfsdk:"daily_schedule"`
	HourlySchedule  *WindowsRemediationScriptHourlyScheduleModel  `tfsdk:"hourly_schedule"`
	RunOnceSchedule *WindowsRemediationScriptRunOnceScheduleModel `tfsdk:"run_once_schedule"`
}

// WindowsRemediationScriptDailyScheduleModel defines the schema for a daily schedule.
type WindowsRemediationScriptDailyScheduleModel struct {
	Interval types.Int32  `tfsdk:"interval"` // Days between runs
	Time     types.String `tfsdk:"time"`     // Time of day in format "HH:MM:SS"
	UseUtc   types.Bool   `tfsdk:"use_utc"`  // Whether to use UTC time
}

// WindowsRemediationScriptHourlyScheduleModel defines the schema for an hourly schedule.
type WindowsRemediationScriptHourlyScheduleModel struct {
	Interval types.Int32 `tfsdk:"interval"` // Hours between runs
}

// WindowsRemediationScriptRunOnceScheduleModel defines the schema for a run once schedule.
type WindowsRemediationScriptRunOnceScheduleModel struct {
	Date   types.String `tfsdk:"date"`    // Date in format "YYYY-MM-DD"
	Time   types.String `tfsdk:"time"`    // Time of day in format "HH:MM:SS"
	UseUtc types.Bool   `tfsdk:"use_utc"` // Whether to use UTC time
}

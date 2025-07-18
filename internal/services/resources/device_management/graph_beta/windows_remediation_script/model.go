// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicehealthscript?view=graph-rest-beta
package graphBetaWindowsRemediationScript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceHealthScriptResourceModel defines the schema for a Device Health Script.
type DeviceHealthScriptResourceModel struct {
	ID                        types.String                                      `tfsdk:"id"`
	DisplayName               types.String                                      `tfsdk:"display_name"`
	Description               types.String                                      `tfsdk:"description"`
	Publisher                 types.String                                      `tfsdk:"publisher"`
	RunAs32Bit                types.Bool                                        `tfsdk:"run_as_32_bit"`
	RunAsAccount              types.String                                      `tfsdk:"run_as_account"`
	EnforceSignatureCheck     types.Bool                                        `tfsdk:"enforce_signature_check"`
	DetectionScriptContent    types.String                                      `tfsdk:"detection_script_content"`
	RemediationScriptContent  types.String                                      `tfsdk:"remediation_script_content"`
	RoleScopeTagIds           types.Set                                         `tfsdk:"role_scope_tag_ids"`
	Version                   types.String                                      `tfsdk:"version"`
	IsGlobalScript            types.Bool                                        `tfsdk:"is_global_script"`
	DeviceHealthScriptType    types.String                                      `tfsdk:"device_health_script_type"`
	CreatedDateTime           types.String                                      `tfsdk:"created_date_time"`
	LastModifiedDateTime      types.String                                      `tfsdk:"last_modified_date_time"`
	HighestAvailableVersion   types.String                                      `tfsdk:"highest_available_version"`
	DetectionScriptParameters types.List                                        `tfsdk:"detection_script_parameters"`
	Assignments               []WindowsRemediationScriptAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts                  timeouts.Value                                    `tfsdk:"timeouts"`
}

// DeviceHealthScriptParameterModel defines a single script parameter.
type DeviceHealthScriptParameterModel struct {
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	IsRequired                       types.Bool   `tfsdk:"is_required"`
	ApplyDefaultValueWhenNotAssigned types.Bool   `tfsdk:"apply_default_value_when_not_assigned"`
}

// WindowsRemediationScriptAssignmentResourceModel represents the assignment configuration
// This maps directly to the API structure where assignments are a flat list of DeviceHealthScriptAssignment objects
type WindowsRemediationScriptAssignmentResourceModel struct {
	// All Devices Assignment
	AllDevices      types.Bool   `tfsdk:"all_devices"`
	AllUsers        types.Bool   `tfsdk:"all_users"`
	Type            types.String `tfsdk:"type"`
	FilterId        types.String `tfsdk:"filter_id"`
	IncludeGroups   types.Set    `tfsdk:"include_groups"`
	ExcludeGroupIds types.Set    `tfsdk:"exclude_group_ids"`
}

// IncludeGroupResourceModel represents a GroupAssignmentTarget with its properties
// This directly maps to a GroupAssignmentTarget in the API
type IncludeGroupResourceModel struct {
	GroupId     types.String              `tfsdk:"group_id"`
	Type        types.String              `tfsdk:"type"`
	FilterId    types.String              `tfsdk:"filter_id"`
	RunSchedule *RunScheduleResourceModel `tfsdk:"run_schedule"`
}

// RunScheduleResourceModel represents different schedule types for assignments
type RunScheduleResourceModel struct {
	ScheduleType types.String `tfsdk:"schedule_type"` // "daily", "hourly", or "once"
	Interval     types.Int32  `tfsdk:"interval"`
	Time         types.String `tfsdk:"time"`    // For daily and once schedules
	Date         types.String `tfsdk:"date"`    // For once schedule
	UseUtc       types.Bool   `tfsdk:"use_utc"` // For daily and once schedules
}

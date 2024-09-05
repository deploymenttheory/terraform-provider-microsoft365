// https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementscript?view=graph-rest-beta

package graphbetadevicemanagementscript

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceManagementScriptResourceModel struct {
	ID                    types.String                                         `tfsdk:"id"`
	DisplayName           types.String                                         `tfsdk:"display_name"`
	Description           types.String                                         `tfsdk:"description"`
	ScriptContent         types.String                                         `tfsdk:"script_content"`
	CreatedDateTime       types.String                                         `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String                                         `tfsdk:"last_modified_date_time"`
	RunAsAccount          types.String                                         `tfsdk:"run_as_account"`
	EnforceSignatureCheck types.Bool                                           `tfsdk:"enforce_signature_check"`
	FileName              types.String                                         `tfsdk:"file_name"`
	RoleScopeTagIds       []types.String                                       `tfsdk:"role_scope_tag_ids"`
	RunAs32Bit            types.Bool                                           `tfsdk:"run_as_32_bit"`
	Assignments           []DeviceManagementScriptAssignmentResourceModel      `tfsdk:"assignments"`
	GroupAssignments      []DeviceManagementScriptGroupAssignmentResourceModel `tfsdk:"group_assignments"`
	Timeouts              timeouts.Value                                       `tfsdk:"timeouts"`
}

// https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicemanagementscriptassignment?view=graph-rest-beta
type DeviceManagementScriptAssignmentResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Target Target       `tfsdk:"target"`
}

type Target struct {
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	TargetType                                 types.String `tfsdk:"target_type"`
	EntraObjectId                              types.String `tfsdk:"entra_object_id"`
}

// https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicemanagementscriptgroupassignment?view=graph-rest-beta
type DeviceManagementScriptGroupAssignmentResourceModel struct {
	ID            types.String `tfsdk:"id"`
	TargetGroupId types.String `tfsdk:"target_group_id"`
}

// REF https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicemanagementscriptassignment?view=graph-rest-beta
package graphBetaMacosPlatformScriptAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MacosPlatformScriptAssignmentResourceModel struct {
	MacosPlatformScriptId types.String                  `tfsdk:"macos_platform_script_id"`
	ID                    types.String                  `tfsdk:"id"`
	Target                AssignmentTargetResourceModel `tfsdk:"target"`
	Timeouts              timeouts.Value                `tfsdk:"timeouts"`
}

type AssignmentTargetResourceModel struct {
	TargetType                                 types.String `tfsdk:"target_type"`
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	GroupId                                    types.String `tfsdk:"group_id"`
}

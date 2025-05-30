// Base resource REF: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeploymentprofileassignment-create?view=graph-rest-beta
package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsAutopilotDeploymentProfileAssignmentResourceModel struct {
	WindowsAutopilotDeploymentProfileId types.String                  `tfsdk:"windows_autopilot_deployment_profile_id"`
	ID                                  types.String                  `tfsdk:"id"`
	Target                              AssignmentTargetResourceModel `tfsdk:"target"`
	Source                              types.String                  `tfsdk:"source"`
	SourceId                            types.String                  `tfsdk:"source_id"`
	Timeouts                            timeouts.Value                `tfsdk:"timeouts"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta
type AssignmentTargetResourceModel struct {
	TargetType                                 types.String `tfsdk:"target_type"` // allDevices, allLicensedUsers, exclusionGroupAssignment, groupAssignment
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	GroupId                                    types.String `tfsdk:"group_id"`
}

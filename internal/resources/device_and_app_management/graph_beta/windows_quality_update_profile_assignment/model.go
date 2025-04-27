// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsqualityupdateprofileassignment?view=graph-rest-beta
package graphBetaWindowsQualityUpdateProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsQualityUpdateProfileAssignmentResourceModel defines the resource model for a Windows Quality Update Profile Assignment
type WindowsQualityUpdateProfileAssignmentResourceModel struct {
	ID                            types.String              `tfsdk:"id"`
	WindowsQualityUpdateProfileID types.String              `tfsdk:"windows_quality_update_profile_id"`
	Assignments                   []AssignmentResourceModel `tfsdk:"assignment"`
	Timeouts                      timeouts.Value            `tfsdk:"timeouts"`
}

// AssignmentResourceModel defines a single assignment block
type AssignmentResourceModel struct {
	Target   types.String `tfsdk:"target"` // "include" or "exclude"
	GroupIds types.Set    `tfsdk:"group_ids"`
}

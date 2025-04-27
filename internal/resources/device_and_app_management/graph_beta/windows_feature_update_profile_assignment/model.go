// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofileassignment?view=graph-rest-beta
package graphBetaWindowsFeatureUpdateProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsFeatureUpdateProfileAssignmentResourceModel defines the resource model for a Windows Driver Update Profile Assignment
type WindowsFeatureUpdateProfileAssignmentResourceModel struct {
	ID                            types.String              `tfsdk:"id"`
	WindowsFeatureUpdateProfileID types.String              `tfsdk:"windows_feature_update_profile_id"`
	Assignments                   []AssignmentResourceModel `tfsdk:"assignment"`
	Timeouts                      timeouts.Value            `tfsdk:"timeouts"`
}

// AssignmentResourceModel defines a single assignment block
type AssignmentResourceModel struct {
	Target   types.String `tfsdk:"target"` // "include" or "exclude"
	GroupIds types.Set    `tfsdk:"group_ids"`
}

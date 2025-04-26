// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofileassignment?view=graph-rest-beta
package graphBetaWindowsDriverUpdateProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsDriverUpdateProfileAssignmentResourceModel defines the resource model for a Windows Driver Update Profile Assignment
type WindowsDriverUpdateProfileAssignmentResourceModel struct {
	ID                           types.String              `tfsdk:"id"`
	WindowsDriverUpdateProfileID types.String              `tfsdk:"windows_driver_update_profile_id"`
	Assignments                  []AssignmentResourceModel `tfsdk:"assignment"`
	Timeouts                     timeouts.Value            `tfsdk:"timeouts"`
}

// AssignmentResourceModel defines a single assignment block
type AssignmentResourceModel struct {
	Target   types.String `tfsdk:"target"`    // "include" or "exclude"
	GroupIds types.Set    `tfsdk:"group_ids"` // Set of group IDs
}

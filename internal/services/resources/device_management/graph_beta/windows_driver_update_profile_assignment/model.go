// Base resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofileassignment?view=graph-rest-beta
package graphBetaWindowsDriverUpdateProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsDriverUpdateProfileAssignmentResourceModel struct {
	WindowsDriverUpdateProfileId types.String                  `tfsdk:"windows_driver_update_profile_id"`
	ID                           types.String                  `tfsdk:"id"`
	Target                       AssignmentTargetResourceModel `tfsdk:"target"`
	Timeouts                     timeouts.Value                `tfsdk:"timeouts"`
}

// Target models
type AssignmentTargetResourceModel struct {
	TargetType   types.String `tfsdk:"target_type"` // configurationManagerCollection, exclusionGroupAssignment, groupAssignment
	GroupId      types.String `tfsdk:"group_id"`
	CollectionId types.String `tfsdk:"collection_id"`
}

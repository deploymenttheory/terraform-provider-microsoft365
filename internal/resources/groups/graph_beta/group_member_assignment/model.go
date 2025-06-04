// REF: https://learn.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-post-members?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-delete-members?view=graph-rest-beta
package graphBetaGroupMemberAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupMemberAssignmentResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	GroupID           types.String   `tfsdk:"group_id"`
	MemberID          types.String   `tfsdk:"member_id"`
	MemberObjectType  types.String   `tfsdk:"member_object_type"`
	MemberType        types.String   `tfsdk:"member_type"`
	MemberDisplayName types.String   `tfsdk:"member_display_name"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/group-list-approleassignments?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-post-approleassignments?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-delete-approleassignments?view=graph-rest-beta
package graphBetaGroupAppRoleAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupAppRoleAssignmentResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	TargetGroupID        types.String   `tfsdk:"target_group_id"`
	ResourceObjectID     types.String   `tfsdk:"resource_object_id"`
	AppRoleID            types.String   `tfsdk:"app_role_id"`
	PrincipalDisplayName types.String   `tfsdk:"principal_display_name"`
	ResourceDisplayName  types.String   `tfsdk:"resource_display_name"`
	PrincipalType        types.String   `tfsdk:"principal_type"`
	CreationTimestamp    types.String   `tfsdk:"creation_timestamp"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

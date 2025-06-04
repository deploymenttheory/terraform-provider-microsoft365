// REF: https://learn.microsoft.com/en-us/graph/api/group-list-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-post-owners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/group-delete-owners?view=graph-rest-beta
package graphBetaGroupOwnerAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupOwnerAssignmentResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	GroupID          types.String   `tfsdk:"group_id"`
	OwnerID          types.String   `tfsdk:"owner_id"`
	OwnerObjectType  types.String   `tfsdk:"owner_object_type"`
	OwnerType        types.String   `tfsdk:"owner_type"`
	OwnerDisplayName types.String   `tfsdk:"owner_display_name"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

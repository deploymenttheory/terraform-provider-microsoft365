// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-list-manager?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-post-manager?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/agentuser-delete-manager?view=graph-rest-beta
package graphBetaUsersUserManager

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserManagerResourceModel describes the resource data model.
type UserManagerResourceModel struct {
	ID        types.String   `tfsdk:"id"`
	UserID    types.String   `tfsdk:"user_id"`
	ManagerID types.String   `tfsdk:"manager_id"`
	Timeouts  timeouts.Value `tfsdk:"timeouts"`
}

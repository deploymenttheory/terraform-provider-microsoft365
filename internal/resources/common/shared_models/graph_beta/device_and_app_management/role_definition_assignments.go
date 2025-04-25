// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-roleassignment?view=graph-rest-beta
package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type RoleAssignmentResourceModel struct {
	ID             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	ScopeMembers   types.Set    `tfsdk:"admin_group_users_group_ids"`
	ScopeType      types.String `tfsdk:"scope_type"`
	ResourceScopes types.Set    `tfsdk:"resource_scopes"`
}

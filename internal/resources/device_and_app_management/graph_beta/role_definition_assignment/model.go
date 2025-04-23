// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-roleassignment?view=graph-rest-beta
package graphBetaRoleDefinitionAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleDefinitionAssignmentResourceModel represents the Terraform resource model for role assignments
type RoleDefinitionAssignmentResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	RoleDefinitionID types.String   `tfsdk:"role_definition_id"`
	BuiltInRoleName  types.String   `tfsdk:"built_in_role_name"`
	DisplayName      types.String   `tfsdk:"display_name"`
	Description      types.String   `tfsdk:"description"`
	ScopeMembers     types.Set      `tfsdk:"scope_members"`
	ScopeType        types.String   `tfsdk:"scope_type"`
	ResourceScopes   types.Set      `tfsdk:"resource_scopes"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

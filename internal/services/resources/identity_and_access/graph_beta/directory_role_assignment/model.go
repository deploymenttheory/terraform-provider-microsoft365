// REF: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-beta
package graphBetaDirectoryRoleAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DirectoryRoleAssignmentResourceModel represents the schema for the directory role assignment resource.
type DirectoryRoleAssignmentResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	RoleDefinitionID types.String   `tfsdk:"role_definition_id"`
	PrincipalID      types.String   `tfsdk:"principal_id"`
	DirectoryScopeID types.String   `tfsdk:"directory_scope_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

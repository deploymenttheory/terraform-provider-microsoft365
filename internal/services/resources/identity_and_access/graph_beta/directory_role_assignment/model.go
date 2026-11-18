package graphBetaDirectoryRoleAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DirectoryRoleAssignmentResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	RoleDefinitionID types.String   `tfsdk:"role_definition_id"`
	PrincipalID      types.String   `tfsdk:"principal_id"`
	DirectoryScopeID types.String   `tfsdk:"directory_scope_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

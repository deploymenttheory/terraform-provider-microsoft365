// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-roledefinition?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-rolepermission?view=graph-rest-beta
package graphBetaRoleDefinition

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleDefinitionResourceModel struct {
	ID                      types.String                              `tfsdk:"id"`
	DisplayName             types.String                              `tfsdk:"display_name"`
	Description             types.String                              `tfsdk:"description"`
	IsBuiltIn               types.Bool                                `tfsdk:"is_built_in"`
	IsBuiltInRoleDefinition types.Bool                                `tfsdk:"is_built_in_role_definition"`
	RolePermissions         []RolePermissionResourceModel             `tfsdk:"role_permissions"`
	RoleScopeTagIds         types.Set                                 `tfsdk:"role_scope_tag_ids"`
	Assignments             *sharedmodels.RoleAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts                timeouts.Value                            `tfsdk:"timeouts"`
}

type RolePermissionResourceModel struct {
	AllowedResourceActions types.Set `tfsdk:"allowed_resource_actions"`
}

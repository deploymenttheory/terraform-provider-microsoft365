// REF: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroledefinition?view=graph-rest-beta
package graphBetaIdentityAndAccessRoleDefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleDefinitionResourceModel represents the schema for the custom Entra directory role definition resource
type RoleDefinitionResourceModel struct {
	ID              types.String               `tfsdk:"id"`
	DisplayName     types.String               `tfsdk:"display_name"`
	Description     types.String               `tfsdk:"description"`
	IsEnabled       types.Bool                 `tfsdk:"is_enabled"`
	TemplateId      types.String               `tfsdk:"template_id"`
	Version         types.String               `tfsdk:"version"`
	RolePermissions []RoleDefinitionPermission `tfsdk:"role_permissions"`
	ResourceScopes  types.Set                  `tfsdk:"resource_scopes"`
	IsBuiltIn       types.Bool                 `tfsdk:"is_built_in"`
	IsPrivileged    types.Bool                 `tfsdk:"is_privileged"`
	Timeouts        timeouts.Value             `tfsdk:"timeouts"`
}

// RoleDefinitionPermission represents one unifiedRolePermission entry of the role definition
type RoleDefinitionPermission struct {
	AllowedResourceActions types.Set `tfsdk:"allowed_resource_actions"`
}

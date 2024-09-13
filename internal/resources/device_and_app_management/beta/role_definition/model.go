// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-roledefinition?view=graph-rest-beta
package graphbetaroledefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleDefinitionResourceModel struct {
	ID                      types.String                  `tfsdk:"id" json:"id,omitempty"`
	DisplayName             types.String                  `tfsdk:"display_name" json:"displayName,omitempty"`
	Description             types.String                  `tfsdk:"description" json:"description,omitempty"`
	IsBuiltIn               types.Bool                    `tfsdk:"is_built_in" json:"isBuiltIn,omitempty"`
	IsBuiltInRoleDefinition types.Bool                    `tfsdk:"is_built_in_role_definition" json:"isBuiltInRoleDefinition,omitempty"`
	RolePermissions         []RolePermissionResourceModel `tfsdk:"role_permissions" json:"rolePermissions,omitempty"`
	RoleScopeTagIds         []types.String                `tfsdk:"role_scope_tag_ids" json:"roleScopeTagIds,omitempty"`
	Timeouts                timeouts.Value                `tfsdk:"timeouts" json:"-"` // Exclude from JSON
}

type RolePermissionResourceModel struct {
	Actions         []types.String                `tfsdk:"actions" json:"actions,omitempty"`
	ResourceActions []ResourceActionResourceModel `tfsdk:"resource_actions" json:"resourceActions,omitempty"`
}

type ResourceActionResourceModel struct {
	AllowedResourceActions    []types.String `tfsdk:"allowed_resource_actions" json:"allowedResourceActions,omitempty"`
	NotAllowedResourceActions []types.String `tfsdk:"not_allowed_resource_actions" json:"notAllowedResourceActions,omitempty"`
}

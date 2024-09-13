// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-roledefinition?view=graph-rest-beta
package graphbetaroledefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleDefinitionResourceModel struct {
	ID              types.String                  `tfsdk:"id"`
	DisplayName     types.String                  `tfsdk:"display_name"`
	Description     types.String                  `tfsdk:"description"`
	IsBuiltIn       types.Bool                    `tfsdk:"is_built_in"`
	RolePermissions []RolePermissionResourceModel `tfsdk:"role_permissions"`
	RoleScopeTagIds []types.String                `tfsdk:"role_scope_tag_ids"`
	Timeouts        timeouts.Value                `tfsdk:"timeouts"`
}

type RolePermissionResourceModel struct {
	ResourceActions []ResourceActionResourceModel `tfsdk:"resource_actions"`
}

type ResourceActionResourceModel struct {
	AllowedResourceActions    []types.String `tfsdk:"allowed_resource_actions"`
	NotAllowedResourceActions []types.String `tfsdk:"not_allowed_resource_actions"`
}

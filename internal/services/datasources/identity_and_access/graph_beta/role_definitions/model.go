// REF: https://learn.microsoft.com/en-us/graph/api/rbacapplication-list-roledefinitions?view=graph-rest-beta&tabs=http

package graphBetaRoleDefinitions

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleDefinitionsDataSourceModel struct {
	FilterType   types.String          `tfsdk:"filter_type"`
	FilterValue  types.String          `tfsdk:"filter_value"`
	ODataFilter  types.String          `tfsdk:"odata_filter"`
	ODataTop     types.Int32           `tfsdk:"odata_top"`
	ODataSkip    types.Int32           `tfsdk:"odata_skip"`
	ODataSelect  types.String          `tfsdk:"odata_select"`
	ODataOrderBy types.String          `tfsdk:"odata_orderby"`
	ODataCount   types.Bool            `tfsdk:"odata_count"`
	ODataSearch  types.String          `tfsdk:"odata_search"`
	ODataExpand  types.String          `tfsdk:"odata_expand"`
	Items        []RoleDefinitionModel `tfsdk:"items"`
	Timeouts     timeouts.Value        `tfsdk:"timeouts"`
}

type RoleDefinitionModel struct {
	ID              types.String          `tfsdk:"id"`
	Description     types.String          `tfsdk:"description"`
	DisplayName     types.String          `tfsdk:"display_name"`
	IsBuiltIn       types.Bool            `tfsdk:"is_built_in"`
	IsEnabled       types.Bool            `tfsdk:"is_enabled"`
	IsPrivileged    types.Bool            `tfsdk:"is_privileged"`
	ResourceScopes  []types.String        `tfsdk:"resource_scopes"`
	TemplateID      types.String          `tfsdk:"template_id"`
	Version         types.String          `tfsdk:"version"`
	RolePermissions []RolePermissionModel `tfsdk:"role_permissions"`
}

type RolePermissionModel struct {
	AllowedResourceActions []types.String `tfsdk:"allowed_resource_actions"`
	Condition              types.String   `tfsdk:"condition"`
}

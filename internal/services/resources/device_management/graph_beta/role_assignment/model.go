// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-rbac-deviceandappmanagementroleassignment?view=graph-rest-beta
package graphBetaRoleDefinitionAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleAssignmentResourceModel struct {
	ID               types.String                      `tfsdk:"id"`
	DisplayName      types.String                      `tfsdk:"display_name"`
	Description      types.String                      `tfsdk:"description"`
	RoleDefinitionId types.String                      `tfsdk:"role_definition_id"`
	Members          types.Set                         `tfsdk:"members"`
	ScopeConfig      []ScopeConfigurationResourceModel `tfsdk:"scope_configuration"`
	Timeouts         timeouts.Value                    `tfsdk:"timeouts"`
}

type ScopeConfigurationResourceModel struct {
	Type           types.String `tfsdk:"type"`
	ResourceScopes types.Set    `tfsdk:"resource_scopes"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/scopedrolemembership?view=graph-rest-beta
package graphBetaAdministrativeUnitRoleAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AdministrativeUnitRoleAssignmentResourceModel represents the schema for the scoped role membership resource
type AdministrativeUnitRoleAssignmentResourceModel struct {
	ID                          types.String   `tfsdk:"id"`
	AdministrativeUnitID        types.String   `tfsdk:"administrative_unit_id"`
	RoleID                      types.String   `tfsdk:"role_id"`
	RoleMemberID                types.String   `tfsdk:"role_member_id"`
	RoleMemberDisplayName       types.String   `tfsdk:"role_member_display_name"`
	RoleMemberUserPrincipalName types.String   `tfsdk:"role_member_user_principal_name"`
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
}

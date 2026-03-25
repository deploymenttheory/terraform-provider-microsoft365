// REF: https://learn.microsoft.com/en-us/graph/api/administrativeunit-list-members?view=graph-rest-beta
package graphBetaAdministrativeUnitMembership

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AdministrativeUnitMembershipResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	AdministrativeUnitID types.String   `tfsdk:"administrative_unit_id"`
	Members              types.Set      `tfsdk:"members"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/administrativeunit?view=graph-rest-beta
package graphBetaAdministrativeUnit

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AdministrativeUnitResourceModel represents the schema for the AdministrativeUnit resource
type AdministrativeUnitResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	Description                   types.String   `tfsdk:"description"`
	IsMemberManagementRestricted  types.Bool     `tfsdk:"is_member_management_restricted"`
	MembershipRule                types.String   `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String   `tfsdk:"membership_rule_processing_state"`
	MembershipType                types.String   `tfsdk:"membership_type"`
	Visibility                    types.String   `tfsdk:"visibility"`
	HardDelete                    types.Bool     `tfsdk:"hard_delete"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/authenticationstrengths-overview?view=graph-rest-beta
package graphBetaAuthenticationStrength

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthenticationStrengthResourceModel represents the schema for the Authentication Strength Policy resource
type AuthenticationStrengthResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	Description           types.String   `tfsdk:"description"`
	PolicyType            types.String   `tfsdk:"policy_type"`
	RequirementsSatisfied types.String   `tfsdk:"requirements_satisfied"`
	CreatedDateTime       types.String   `tfsdk:"created_date_time"`
	ModifiedDateTime      types.String   `tfsdk:"modified_date_time"`
	AllowedCombinations   types.Set      `tfsdk:"allowed_combinations"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

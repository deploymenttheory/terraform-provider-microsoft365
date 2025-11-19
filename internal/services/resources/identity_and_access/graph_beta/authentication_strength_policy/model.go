// REF: https://learn.microsoft.com/en-us/graph/api/resources/authenticationstrengthpolicy?view=graph-rest-beta
package graphBetaAuthenticationStrengthPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthenticationStrengthPolicyResourceModel represents the schema for the Authentication Strength Policy resource
type AuthenticationStrengthPolicyResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	DisplayName               types.String   `tfsdk:"display_name"`
	Description               types.String   `tfsdk:"description"`
	PolicyType                types.String   `tfsdk:"policy_type"`
	RequirementsSatisfied     types.String   `tfsdk:"requirements_satisfied"`
	CreatedDateTime           types.String   `tfsdk:"created_date_time"`
	ModifiedDateTime          types.String   `tfsdk:"modified_date_time"`
	AllowedCombinations       types.Set      `tfsdk:"allowed_combinations"`
	CombinationConfigurations types.List     `tfsdk:"combination_configurations"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`
}

// CombinationConfigurationModel represents a combination configuration (base type)
type CombinationConfigurationModel struct {
	ID                    types.String `tfsdk:"id"`
	ODataType             types.String `tfsdk:"odata_type"`
	AppliesToCombinations types.String `tfsdk:"applies_to_combinations"`

	// FIDO2 specific fields
	AllowedAAGUIDs types.Set `tfsdk:"allowed_aaguids"`

	// X.509 Certificate specific fields
	AllowedIssuerSkis types.Set `tfsdk:"allowed_issuer_skis"`
	AllowedIssuers    types.Set `tfsdk:"allowed_issuers"`
	AllowedPolicyOIDs types.Set `tfsdk:"allowed_policy_oids"`
}

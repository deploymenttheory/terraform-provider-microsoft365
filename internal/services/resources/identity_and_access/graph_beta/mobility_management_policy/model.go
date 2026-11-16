// REF: https://learn.microsoft.com/en-us/graph/api/resources/mobilitymanagementpolicy?view=graph-rest-beta
package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobilityManagementPolicyResourceModel represents the schema for the mobility management policy resource
type MobilityManagementPolicyResourceModel struct {
	ID                                        types.String   `tfsdk:"id"`
	PolicyType                                types.String   `tfsdk:"policy_type"`
	PolicyID                                  types.String   `tfsdk:"policy_id"`
	DisplayName                               types.String   `tfsdk:"display_name"`
	Description                               types.String   `tfsdk:"description"`
	IsValid                                   types.Bool     `tfsdk:"is_valid"`
	AppliesTo                                 types.String   `tfsdk:"applies_to"`
	IncludedGroups                            types.Set      `tfsdk:"included_groups"`
	ComplianceUrl                             types.String   `tfsdk:"compliance_url"`
	DiscoveryUrl                              types.String   `tfsdk:"discovery_url"`
	TermsOfUseUrl                             types.String   `tfsdk:"terms_of_use_url"`
	IsMdmEnrollmentDuringRegistrationDisabled types.Bool     `tfsdk:"is_mdm_enrollment_during_registration_disabled"`
	Timeouts                                  timeouts.Value `tfsdk:"timeouts"`
}

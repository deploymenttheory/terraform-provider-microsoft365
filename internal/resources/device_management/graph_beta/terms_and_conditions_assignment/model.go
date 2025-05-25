// resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-companyterms-termsandconditionsassignment?view=graph-rest-beta
package graphBetaTermsAndConditionsAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TermsAndConditionsAssignmentResourceModel struct {
	TermsAndConditionsId types.String                  `tfsdk:"terms_and_conditions_id"`
	ID                   types.String                  `tfsdk:"id"`
	Target               AssignmentTargetResourceModel `tfsdk:"target"`
	Timeouts             timeouts.Value                `tfsdk:"timeouts"`
}

// Target models
type AssignmentTargetResourceModel struct {
	TargetType   types.String `tfsdk:"target_type"` // allLicensedUsers, allDevices, groupAssignment, exclusionGroupAssignment, configurationManagerCollection
	GroupId      types.String `tfsdk:"group_id"`
	CollectionId types.String `tfsdk:"collection_id"`
}

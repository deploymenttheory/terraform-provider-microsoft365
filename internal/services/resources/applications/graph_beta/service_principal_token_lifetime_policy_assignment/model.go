// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-tokenlifetimepolicies?view=graph-rest-beta
package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServicePrincipalTokenLifetimePolicyAssignmentResourceModel represents the schema for the SP token lifetime policy assignment resource.
// The Microsoft Graph API does not return an ID for this assignment; ID is a provider-constructed composite of
// service_principal_id and token_lifetime_policy_id, used to track the assignment in Terraform state.
type ServicePrincipalTokenLifetimePolicyAssignmentResourceModel struct {
	ID                    types.String   `tfsdk:"id"` // synthetic: "{service_principal_id}/{token_lifetime_policy_id}"
	ServicePrincipalID    types.String   `tfsdk:"service_principal_id"`
	TokenLifetimePolicyID types.String   `tfsdk:"token_lifetime_policy_id"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

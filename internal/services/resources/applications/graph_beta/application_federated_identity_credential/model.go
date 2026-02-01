// REF: https://learn.microsoft.com/en-us/graph/api/application-list-federatedidentitycredentials?view=graph-rest-beta&tabs=http
// REF: https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-get?view=graph-rest-beta&tabs=http
package graphBetaApplicationFederatedIdentityCredential

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationFederatedIdentityCredentialResourceModel struct {
	ID                       types.String   `tfsdk:"id"`
	ApplicationID            types.String   `tfsdk:"application_id"`
	Name                     types.String   `tfsdk:"name"`
	Issuer                   types.String   `tfsdk:"issuer"`
	Subject                  types.String   `tfsdk:"subject"`
	Audiences                types.Set      `tfsdk:"audiences"`
	Description              types.String   `tfsdk:"description"`
	ClaimsMatchingExpression types.String   `tfsdk:"claims_matching_expression"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}

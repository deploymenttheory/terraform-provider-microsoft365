// REF: https://learn.microsoft.com/en-us/graph/api/resources/identitysecuritydefaultsenforcementpolicy?view=graph-rest-beta
package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecurityDefaultsEnforcementPolicyResourceModel represents the schema for the SecurityDefaultsEnforcementPolicy resource.
type SecurityDefaultsEnforcementPolicyResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	IsEnabled   types.Bool     `tfsdk:"is_enabled"`
	DisplayName types.String   `tfsdk:"display_name"`
	Description types.String   `tfsdk:"description"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

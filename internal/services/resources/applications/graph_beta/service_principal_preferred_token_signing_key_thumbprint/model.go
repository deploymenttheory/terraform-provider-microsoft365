// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta
package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel struct {
	Id                 types.String   `tfsdk:"id"`
	ServicePrincipalID types.String   `tfsdk:"service_principal_id"`
	Thumbprint         types.String   `tfsdk:"thumbprint"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

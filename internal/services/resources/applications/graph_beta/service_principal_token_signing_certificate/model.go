// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-addtokensigningcertificate?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/selfsignedcertificate?view=graph-rest-beta
package graphBetaServicePrincipalTokenSigningCertificate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServicePrincipalTokenSigningCertificateResourceModel struct {
	Id                 types.String   `tfsdk:"id"`
	ServicePrincipalID types.String   `tfsdk:"service_principal_id"`
	DisplayName        types.String   `tfsdk:"display_name"`
	EndDateTime        types.String   `tfsdk:"end_date_time"`
	KeyId              types.String   `tfsdk:"key_id"`
	StartDateTime      types.String   `tfsdk:"start_date_time"`
	Thumbprint         types.String   `tfsdk:"thumbprint"`
	Value              types.String   `tfsdk:"value"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

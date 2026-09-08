package graphBetaNetworkManagedTLSCertificate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkManagedTLSCertificateResourceModel represents a Microsoft-managed
// certificate authority used by Global Secure Access TLS inspection.
type NetworkManagedTLSCertificateResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	CommonName            types.String   `tfsdk:"common_name"`
	OrganizationName      types.String   `tfsdk:"organization_name"`
	ValidityMonths        types.Int32    `tfsdk:"validity_months"`
	Enabled               types.Bool     `tfsdk:"enabled"`
	Status                types.String   `tfsdk:"status"`
	CreatedDateTime       types.String   `tfsdk:"created_date_time"`
	ValidityStartDateTime types.String   `tfsdk:"validity_start_date_time"`
	ValidityEndDateTime   types.String   `tfsdk:"validity_end_date_time"`
	Certificate           types.String   `tfsdk:"certificate"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

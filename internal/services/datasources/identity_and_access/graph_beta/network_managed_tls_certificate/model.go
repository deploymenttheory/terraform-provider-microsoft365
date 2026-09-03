package graphBetaNetworkManagedTLSCertificate

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkManagedTLSCertificateDataSourceModel struct {
	ID                     types.String   `tfsdk:"id"`
	CertificateAuthorityID types.String   `tfsdk:"certificate_authority_id"`
	Name                   types.String   `tfsdk:"name"`
	CommonName             types.String   `tfsdk:"common_name"`
	OrganizationName       types.String   `tfsdk:"organization_name"`
	ValidityMonths         types.Int32    `tfsdk:"validity_months"`
	Status                 types.String   `tfsdk:"status"`
	CreatedDateTime        types.String   `tfsdk:"created_date_time"`
	ValidityStartDateTime  types.String   `tfsdk:"validity_start_date_time"`
	ValidityEndDateTime    types.String   `tfsdk:"validity_end_date_time"`
	Certificate            types.String   `tfsdk:"certificate"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}

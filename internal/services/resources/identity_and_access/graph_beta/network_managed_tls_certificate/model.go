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

// managedTLSCertificateRequest is the managedCertificateAuthorityCertificate create body, and the status-only update
// body. Fields left nil are omitted, as the Entra portal omits status on create.
type managedTLSCertificateRequest struct {
	Name             *string `json:"name,omitempty"`
	CommonName       *string `json:"commonName,omitempty"`
	OrganizationName *string `json:"organizationName,omitempty"`
	ValidityMonths   *int32  `json:"validityMonths,omitempty"`
	Status           *string `json:"status,omitempty"`
}

// managedTLSCertificateResponse holds the managedCertificateAuthorityCertificate fields this resource reads.
type managedTLSCertificateResponse struct {
	ID               *string                                `json:"id"`
	Name             *string                                `json:"name"`
	CommonName       *string                                `json:"commonName"`
	OrganizationName *string                                `json:"organizationName"`
	ValidityMonths   *int32                                 `json:"validityMonths"`
	Status           *string                                `json:"status"`
	CreatedDateTime  *string                                `json:"createdDateTime"`
	Certificate      *string                                `json:"certificate"`
	Validity         *managedTLSCertificateValidityResponse `json:"validity"`
}

// managedTLSCertificateValidityResponse is the certificate's validity window.
type managedTLSCertificateValidityResponse struct {
	StartDateTime *string `json:"startDateTime"`
	EndDateTime   *string `json:"endDateTime"`
}

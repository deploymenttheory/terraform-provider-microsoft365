package graphBetaNetworkManagedTLSCertificate

import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"

// MapRemoteStateToDataSource maps the certificate to the data source model, keeping the configured id.
func MapRemoteStateToDataSource(data *NetworkManagedTLSCertificateDataSourceModel, remote *managedTLSCertificateResponse) {
	if remote == nil {
		return
	}

	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.CommonName = convert.GraphToFrameworkString(remote.CommonName)
	data.OrganizationName = convert.GraphToFrameworkString(remote.OrganizationName)
	data.ValidityMonths = convert.GraphToFrameworkInt32(remote.ValidityMonths)
	data.Status = convert.GraphToFrameworkString(remote.Status)
	data.CreatedDateTime = convert.GraphToFrameworkString(remote.CreatedDateTime)
	data.Certificate = convert.GraphToFrameworkString(remote.Certificate)
	if remote.Validity != nil {
		data.ValidityStartDateTime = convert.GraphToFrameworkString(remote.Validity.StartDateTime)
		data.ValidityEndDateTime = convert.GraphToFrameworkString(remote.Validity.EndDateTime)
	}
}

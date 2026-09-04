package graphBetaNetworkManagedTLSCertificate

import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"

func mapRemoteStateToDataSource(data *NetworkManagedTLSCertificateDataSourceModel, remote *managedTLSCertificateResponse) {
	if remote == nil {
		return
	}

	data.Name = convert.GraphToFrameworkString(remote.name)
	data.CommonName = convert.GraphToFrameworkString(remote.commonName)
	data.OrganizationName = convert.GraphToFrameworkString(remote.organizationName)
	data.ValidityMonths = convert.GraphToFrameworkInt32(remote.validityMonths)
	data.Status = convert.GraphToFrameworkString(remote.status)
	data.CreatedDateTime = convert.GraphToFrameworkString(remote.createdDateTime)
	data.Certificate = convert.GraphToFrameworkString(remote.certificate)
	if remote.validity != nil {
		data.ValidityStartDateTime = convert.GraphToFrameworkString(remote.validity.startDateTime)
		data.ValidityEndDateTime = convert.GraphToFrameworkString(remote.validity.endDateTime)
	}
}

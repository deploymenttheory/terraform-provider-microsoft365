package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapRemoteStateToTerraform(_ context.Context, data *NetworkManagedTLSCertificateResourceModel, remote *managedTLSCertificateResponse) {
	if remote == nil {
		return
	}

	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.CommonName = convert.GraphToFrameworkString(remote.commonName)
	data.OrganizationName = convert.GraphToFrameworkString(remote.organizationName)
	data.ValidityMonths = convert.GraphToFrameworkInt32(remote.validityMonths)
	data.Status = convert.GraphToFrameworkString(remote.status)
	if remote.status != nil {
		switch strings.ToLower(*remote.status) {
		case "active", "enabled", "enrolling", "expiring":
			data.Enabled = types.BoolValue(true)
		case "unknownfuturevalue", "creating", "disabled", "expired", "revoked":
			data.Enabled = types.BoolValue(false)
		}
	}
	data.CreatedDateTime = convert.GraphToFrameworkString(remote.createdDateTime)
	data.Certificate = convert.GraphToFrameworkString(remote.certificate)

	if remote.validity != nil {
		data.ValidityStartDateTime = convert.GraphToFrameworkString(remote.validity.startDateTime)
		data.ValidityEndDateTime = convert.GraphToFrameworkString(remote.validity.endDateTime)
	}
}

package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(_ context.Context, data *NetworkManagedTLSCertificateResourceModel, remote *managedTLSCertificateResponse) {
	if remote == nil {
		return
	}

	data.ID = convert.GraphToFrameworkString(remote.ID)
	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.CommonName = convert.GraphToFrameworkString(remote.CommonName)
	data.OrganizationName = convert.GraphToFrameworkString(remote.OrganizationName)
	data.ValidityMonths = convert.GraphToFrameworkInt32(remote.ValidityMonths)
	data.Status = convert.GraphToFrameworkString(remote.Status)
	if remote.Status != nil {
		switch strings.ToLower(*remote.Status) {
		case "active", "enabled", "enrolling", "expiring":
			data.Enabled = types.BoolValue(true)
		case "unknownfuturevalue", "creating", "disabled", "expired", "revoked":
			data.Enabled = types.BoolValue(false)
		}
	}
	data.CreatedDateTime = convert.GraphToFrameworkString(remote.CreatedDateTime)
	data.Certificate = convert.GraphToFrameworkString(remote.Certificate)

	if remote.Validity != nil {
		data.ValidityStartDateTime = convert.GraphToFrameworkString(remote.Validity.StartDateTime)
		data.ValidityEndDateTime = convert.GraphToFrameworkString(remote.Validity.EndDateTime)
	}
}

// mapManagedTLSCertificateComputedToState copies the computed values from a response, leaving configured
// attributes as planned, and replaces any value Graph did not return with null.
func mapManagedTLSCertificateComputedToState(ctx context.Context, object *NetworkManagedTLSCertificateResourceModel, remote *managedTLSCertificateResponse) {
	if remote != nil {
		remoteState := *object
		MapRemoteStateToTerraform(ctx, &remoteState, remote)
		object.Status = remoteState.Status
		object.CreatedDateTime = remoteState.CreatedDateTime
		object.ValidityStartDateTime = remoteState.ValidityStartDateTime
		object.ValidityEndDateTime = remoteState.ValidityEndDateTime
		object.Certificate = remoteState.Certificate
	}
	for _, v := range []*types.String{
		&object.Status, &object.CreatedDateTime, &object.ValidityStartDateTime,
		&object.ValidityEndDateTime, &object.Certificate,
	} {
		if v.IsUnknown() {
			*v = types.StringNull()
		}
	}
}

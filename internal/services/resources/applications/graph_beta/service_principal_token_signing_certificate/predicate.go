package graphBetaServicePrincipalTokenSigningCertificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// tokenSigningCertificateConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the certificate's key credential has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate checks the full read state against the expected state captured at write time:
//
//   - id must be non-empty — the composite servicePrincipalId/keyId key; confirms Read
//     successfully located the key credential on the service principal.
//   - service_principal_id must match expected — confirms the read is for the correct service principal.
//   - key_id must match expected — confirms the credential propagated to the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func tokenSigningCertificateConsistencyPredicate(expected *ServicePrincipalTokenSigningCertificateResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalTokenSigningCertificateResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.Id.IsNull() || actual.Id.IsUnknown() || actual.Id.ValueString() == "" {
			return false
		}

		if actual.ServicePrincipalID.ValueString() != expected.ServicePrincipalID.ValueString() {
			return false
		}

		if actual.KeyId.ValueString() != expected.KeyId.ValueString() {
			return false
		}

		return true
	}
}

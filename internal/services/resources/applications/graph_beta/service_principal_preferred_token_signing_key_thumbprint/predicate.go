package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// preferredTokenSigningKeyThumbprintConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the thumbprint has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate checks the full read state against the expected state captured at write time:
//
//   - id must be non-empty — confirms Read successfully located the service principal.
//   - service_principal_id must match expected — confirms the read is for the correct service principal.
//   - thumbprint must match expected (case-insensitive, Graph normalizes to lowercase) —
//     confirms the value propagated to the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func preferredTokenSigningKeyThumbprintConsistencyPredicate(expected *ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.Id.IsNull() || actual.Id.IsUnknown() || actual.Id.ValueString() == "" {
			return false
		}

		if actual.ServicePrincipalID.ValueString() != expected.ServicePrincipalID.ValueString() {
			return false
		}

		if !strings.EqualFold(actual.Thumbprint.ValueString(), expected.Thumbprint.ValueString()) {
			return false
		}

		return true
	}
}

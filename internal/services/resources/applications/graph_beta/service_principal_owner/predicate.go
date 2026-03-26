package graphBetaServicePrincipalOwner

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// servicePrincipalOwnerConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the service principal owner resource state has propagated across Microsoft Entra
// replicas before accepting the read as authoritative.
//
// The predicate checks the full read state against the expected state captured at write time:
//
//   - id must be non-empty — the composite servicePrincipalId/ownerId key; confirms Read
//     located the owner in the service principal's owners list.
//   - service_principal_id must match expected — confirms the read is for the correct SP.
//   - owner_id must match expected — confirms the correct owner is in state.
//   - owner_type must not be empty or "Unknown" — a Computed field resolved by the API; when
//     set to a real type ("User" or "ServicePrincipal") the read has returned full owner data.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func servicePrincipalOwnerConsistencyPredicate(expected *ServicePrincipalOwnerResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalOwnerResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.ServicePrincipalID.ValueString() != expected.ServicePrincipalID.ValueString() {
			return false
		}

		if actual.OwnerID.ValueString() != expected.OwnerID.ValueString() {
			return false
		}

		// owner_type is Computed — resolved by Read from the API ("User" or "ServicePrincipal").
		ownerType := actual.OwnerType.ValueString()
		if ownerType == "" || ownerType == "Unknown" {
			return false
		}

		return true
	}
}

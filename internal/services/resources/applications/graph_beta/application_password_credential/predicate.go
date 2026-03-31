package graphBetaApplicationPasswordCredential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// applicationPasswordCredentialConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the password credential write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate checks that the KeyID is populated in state, confirming that the credential
// was found in the application's password credential list on the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func applicationPasswordCredentialConsistencyPredicate(expected *ApplicationPasswordCredentialResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ApplicationPasswordCredentialResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		// KeyID is set by the API after addPassword. If it is absent the responding
		// replica has not yet returned the new credential.
		if actual.KeyID.IsNull() || actual.KeyID.IsUnknown() || actual.KeyID.ValueString() == "" {
			return false
		}

		// Confirm the key_id matches the one from the create response.
		if !expected.KeyID.IsNull() && !expected.KeyID.IsUnknown() && expected.KeyID.ValueString() != "" {
			if actual.KeyID.ValueString() != expected.KeyID.ValueString() {
				return false
			}
		}

		return true
	}
}

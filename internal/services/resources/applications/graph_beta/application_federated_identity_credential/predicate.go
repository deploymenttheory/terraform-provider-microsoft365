package graphBetaApplicationFederatedIdentityCredential

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// applicationFederatedIdentityCredentialConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the federated identity credential write has propagated across
// Microsoft Entra replicas before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func applicationFederatedIdentityCredentialConsistencyPredicate(expected *ApplicationFederatedIdentityCredentialResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ApplicationFederatedIdentityCredentialResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.ApplicationID.ValueString() != expected.ApplicationID.ValueString() {
			return false
		}

		// Compare all user-specified mutable fields. If any field set in the plan does not
		// match the read-back state, the responding replica has not yet received the write.
		if !actual.Name.Equal(expected.Name) {
			return false
		}
		if !actual.Issuer.Equal(expected.Issuer) {
			return false
		}
		if !actual.Subject.Equal(expected.Subject) {
			return false
		}
		if !expected.Audiences.IsNull() && !expected.Audiences.IsUnknown() {
			if !actual.Audiences.Equal(expected.Audiences) {
				return false
			}
		}
		if !expected.Description.IsNull() && !expected.Description.IsUnknown() {
			if !actual.Description.Equal(expected.Description) {
				return false
			}
		}
		if !expected.ClaimsMatchingExpression.IsNull() && !expected.ClaimsMatchingExpression.IsUnknown() {
			if !actual.ClaimsMatchingExpression.Equal(expected.ClaimsMatchingExpression) {
				return false
			}
		}

		return true
	}
}

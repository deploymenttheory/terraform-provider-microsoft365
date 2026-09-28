package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// b2bManagementPolicyConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// only accepts a read once it reflects the values just written.
//
// Writes to /policies/b2bManagementPolicies are eventually consistent: for several seconds after a
// successful POST or PATCH, reads served by different Microsoft Entra replicas alternate between the
// previous and the new values. Accepting a stale read would store the old values and produce a
// spurious diff on the next plan.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func b2bManagementPolicyConsistencyPredicate(expected *B2bManagementPolicyResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		if state.Raw.IsNull() {
			return false
		}

		var actual B2bManagementPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		return actual.DisplayName.Equal(expected.DisplayName) &&
			actual.Definition.Equal(expected.Definition) &&
			actual.IsOrganizationDefault.Equal(expected.IsOrganizationDefault)
	}
}

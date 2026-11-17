package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var policyRetryInterval = 5 * time.Second

// stableConsistencyPredicate requires four matching reads, spanning at least 15 seconds at
// the normal retry interval. Curl probes showed two current reads followed by old values again.
func stableConsistencyPredicate(expected *B2bManagementPolicyResourceModel) func(context.Context, tfsdk.State) bool {
	matches := b2bManagementPolicyConsistencyPredicate(expected)
	consecutive := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		if !matches(ctx, state) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= 4
	}
}

// b2bManagementPolicyConsistencyPredicate compares a read with the written values.
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

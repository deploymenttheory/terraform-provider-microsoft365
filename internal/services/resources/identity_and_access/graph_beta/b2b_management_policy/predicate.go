package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// b2bManagementPolicyConsistencyPredicate accepts a read only once it reflects the written values;
// for several seconds after a write, Entra replicas alternate between the old and new values.
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

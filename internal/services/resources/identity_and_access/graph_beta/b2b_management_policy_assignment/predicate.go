package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// b2bManagementPolicyAssignmentConsistencyPredicate rejects reads from replicas that do not list the
// new assignment yet; accepting them would leave Create without state.
func b2bManagementPolicyAssignmentConsistencyPredicate(expected *B2bManagementPolicyAssignmentResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		if state.Raw.IsNull() {
			return false
		}

		var actual B2bManagementPolicyAssignmentResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		return actual.B2bManagementPolicyID.ValueString() == expected.B2bManagementPolicyID.ValueString() &&
			actual.DirectoryObjectID.ValueString() == expected.DirectoryObjectID.ValueString()
	}
}

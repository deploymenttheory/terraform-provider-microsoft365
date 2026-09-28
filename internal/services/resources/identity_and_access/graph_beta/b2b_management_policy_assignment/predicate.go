package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// b2bManagementPolicyAssignmentConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the assignment write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// POST /{applications|servicePrincipals}/{id}/policies/$ref is eventually consistent: immediately
// after a successful assignment, GET /policies/b2bManagementPolicies/{id}/appliesTo served by a
// stale replica may not include the directory object yet, in which case Read removes the resource
// from state. Without this predicate that null state would be accepted as a successful read, and
// Create would return no state ("Missing Resource State After Create") despite the assignment
// existing in Entra.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
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

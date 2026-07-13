package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// servicePrincipalTokenLifetimePolicyAssignmentConsistencyPredicate returns a consistency
// predicate for ReadWithRetry that verifies the policy assignment write has propagated across
// Microsoft Entra replicas before accepting the read as authoritative.
//
// POST /servicePrincipals/{id}/tokenLifetimePolicies/$ref is eventually consistent: immediately
// after a successful assignment, GET /servicePrincipals/{id}/tokenLifetimePolicies served by a
// stale replica may not include the policy yet, in which case Read removes the resource from
// state. Without this predicate that null state would be accepted as a successful read, and
// Create would return no state ("Missing Resource State After Create") despite the assignment
// existing in Entra.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func servicePrincipalTokenLifetimePolicyAssignmentConsistencyPredicate(expected *ServicePrincipalTokenLifetimePolicyAssignmentResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		if state.Raw.IsNull() {
			return false
		}

		var actual ServicePrincipalTokenLifetimePolicyAssignmentResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.ServicePrincipalID.ValueString() != expected.ServicePrincipalID.ValueString() {
			return false
		}

		if actual.TokenLifetimePolicyID.ValueString() != expected.TokenLifetimePolicyID.ValueString() {
			return false
		}

		return true
	}
}

package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// servicePrincipalAppRoleAssignedToConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the app role assignment write has propagated across Microsoft
// Entra replicas before accepting the read as authoritative.
//
// The predicate checks that the assignment ID and resource display name are populated in state,
// confirming the assignment is visible on the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func servicePrincipalAppRoleAssignedToConsistencyPredicate(expected *ServicePrincipalAppRoleAssignedToResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalAppRoleAssignedToResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// Confirm the assignment ID matches the one from the create response.
		if !expected.ID.IsNull() && !expected.ID.IsUnknown() && expected.ID.ValueString() != "" {
			if actual.ID.ValueString() != expected.ID.ValueString() {
				return false
			}
		}

		// ResourceDisplayName is populated by the API read; its presence confirms the
		// responding replica has the assignment record.
		if actual.ResourceDisplayName.IsNull() || actual.ResourceDisplayName.IsUnknown() || actual.ResourceDisplayName.ValueString() == "" {
			return false
		}

		return true
	}
}

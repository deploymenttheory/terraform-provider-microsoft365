package graphBetaCrossTenantAccessPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// crossTenantAccessPolicyConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the cross-tenant access policy write has propagated across Microsoft Entra
// replicas before accepting the read as authoritative.
//
// This is a singleton resource. The predicate compares all user-specified mutable fields from
// the expected plan against the state returned by the read. A stale replica returns the
// pre-write version, causing at least one field to differ and triggering a retry.
//
// restore_defaults_on_destroy is a Terraform-only field and is not compared.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func crossTenantAccessPolicyConsistencyPredicate(expected *CrossTenantAccessPolicyResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual CrossTenantAccessPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// display_name is always populated by the API and confirms the replica returned the
		// full policy object.
		if actual.DisplayName.IsNull() || actual.DisplayName.IsUnknown() || actual.DisplayName.ValueString() == "" {
			return false
		}

		// Compare all user-specified mutable fields. If any field set in the plan does not
		// match the read-back state, the responding replica has not yet received the write.
		if !expected.DisplayName.IsNull() && !expected.DisplayName.IsUnknown() {
			if !actual.DisplayName.Equal(expected.DisplayName) {
				return false
			}
		}
		if !expected.AllowedCloudEndpoints.IsNull() && !expected.AllowedCloudEndpoints.IsUnknown() {
			if !actual.AllowedCloudEndpoints.Equal(expected.AllowedCloudEndpoints) {
				return false
			}
		}

		return true
	}
}

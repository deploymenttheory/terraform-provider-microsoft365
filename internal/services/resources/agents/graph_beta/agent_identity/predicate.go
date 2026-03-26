package graphBetaAgentIdentity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// agentIdentityConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the agent identity resource write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// Computed-only fields (created_by_app_id, created_date_time, disabled_by_microsoft_status) are
// not compared — they are set by the API and are not meaningful signals for write propagation.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func agentIdentityConsistencyPredicate(expected *AgentIdentityResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AgentIdentityResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// service_principal_type is Computed — only set by the API, never from plan.
		// Its presence confirms the replica has the full object representation.
		if actual.ServicePrincipalType.IsNull() || actual.ServicePrincipalType.IsUnknown() || actual.ServicePrincipalType.ValueString() == "" {
			return false
		}

		// Compare all user-specified mutable fields. If any field set in the plan does not
		// match the read-back state, the responding replica has not yet received the write.
		if !actual.DisplayName.Equal(expected.DisplayName) {
			return false
		}
		if !expected.AccountEnabled.IsNull() && !expected.AccountEnabled.IsUnknown() {
			if !actual.AccountEnabled.Equal(expected.AccountEnabled) {
				return false
			}
		}
		if !expected.SponsorIds.IsNull() && !expected.SponsorIds.IsUnknown() {
			if !actual.SponsorIds.Equal(expected.SponsorIds) {
				return false
			}
		}
		if !expected.OwnerIds.IsNull() && !expected.OwnerIds.IsUnknown() {
			if !actual.OwnerIds.Equal(expected.OwnerIds) {
				return false
			}
		}

		return true
	}
}

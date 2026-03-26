package graphBetaAdministrativeUnitMembership

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// administrativeUnitMembershipConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the administrative unit membership write has propagated across
// Microsoft Entra replicas before accepting the read as authoritative.
//
// The predicate compares the full member set from the expected plan against the state returned
// by the read. The addMembersById / removeMembersById APIs are eventually consistent; a stale
// replica may return an outdated membership list that does not reflect the written state.
// Retrying until the member set matches ensures the accepted read is authoritative.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func administrativeUnitMembershipConsistencyPredicate(expected *AdministrativeUnitMembershipResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AdministrativeUnitMembershipResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.AdministrativeUnitID.ValueString() != expected.AdministrativeUnitID.ValueString() {
			return false
		}

		// members must match exactly — addMembersById and removeMembersById are eventually
		// consistent; a stale replica may return an incomplete or outdated membership set.
		if !expected.Members.IsNull() && !expected.Members.IsUnknown() {
			if !actual.Members.Equal(expected.Members) {
				return false
			}
		} else {
			// expected is null/unknown — accept only if actual is also null or empty
			if !actual.Members.IsNull() && !actual.Members.IsUnknown() && len(actual.Members.Elements()) > 0 {
				return false
			}
		}

		return true
	}
}

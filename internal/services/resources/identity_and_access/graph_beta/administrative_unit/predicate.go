package graphBetaAdministrativeUnit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// administrativeUnitConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the administrative unit resource write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func administrativeUnitConsistencyPredicate(expected *AdministrativeUnitResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AdministrativeUnitResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// Compare all user-specified mutable fields. If any field set in the plan does not
		// match the read-back state, the responding replica has not yet received the write.
		if !actual.DisplayName.Equal(expected.DisplayName) {
			return false
		}
		if !expected.Description.IsNull() && !expected.Description.IsUnknown() {
			if !actual.Description.Equal(expected.Description) {
				return false
			}
		}
		if !expected.IsMemberManagementRestricted.IsNull() && !expected.IsMemberManagementRestricted.IsUnknown() {
			if !actual.IsMemberManagementRestricted.Equal(expected.IsMemberManagementRestricted) {
				return false
			}
		}
		if !expected.MembershipRule.IsNull() && !expected.MembershipRule.IsUnknown() {
			if !actual.MembershipRule.Equal(expected.MembershipRule) {
				return false
			}
		}
		if !expected.MembershipRuleProcessingState.IsNull() && !expected.MembershipRuleProcessingState.IsUnknown() {
			if !actual.MembershipRuleProcessingState.Equal(expected.MembershipRuleProcessingState) {
				return false
			}
		}
		if !expected.MembershipType.IsNull() && !expected.MembershipType.IsUnknown() {
			if !actual.MembershipType.Equal(expected.MembershipType) {
				return false
			}
		}
		if !expected.Visibility.IsNull() && !expected.Visibility.IsUnknown() {
			if !actual.Visibility.Equal(expected.Visibility) {
				return false
			}
		}

		return true
	}
}

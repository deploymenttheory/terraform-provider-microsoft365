package graphBetaRoleManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// roleManagementPolicyConsistencyPredicate returns a consistency predicate for ReadWithRetry that accepts a read only
// once every configured rule shows the written values.
//
// Role management policy reads are eventually consistent: probes against Graph beta showed a GET issued within about a
// second of a successful PATCH returning the previous rule values. The policy-level lastModifiedDateTime cannot be used
// instead because it does not change when a PATCH writes values the policy already has.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func roleManagementPolicyConsistencyPredicate(expected *RoleManagementPolicyResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual RoleManagementPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		return rulesMatch(expected, &actual)
	}
}

// rulesMatch reports whether every rule configured in expected has the same values in actual.
func rulesMatch(expected, actual *RoleManagementPolicyResourceModel) bool {
	for ruleId, want := range expected.ExpirationRules {
		got, ok := actual.ExpirationRules[ruleId]
		if !ok || !got.IsExpirationRequired.Equal(want.IsExpirationRequired) || !durationsEquivalent(got.MaximumDuration.ValueString(), want.MaximumDuration.ValueString()) {
			return false
		}
	}

	for ruleId, want := range expected.EnablementRules {
		got, ok := actual.EnablementRules[ruleId]
		if !ok || !got.EnabledRules.Equal(want.EnabledRules) {
			return false
		}
	}

	for ruleId, want := range expected.NotificationRules {
		got, ok := actual.NotificationRules[ruleId]
		if !ok || !got.NotificationLevel.Equal(want.NotificationLevel) || !got.IsDefaultRecipientsEnabled.Equal(want.IsDefaultRecipientsEnabled) {
			return false
		}
		if !got.NotificationRecipients.Equal(want.NotificationRecipients) {
			return false
		}
	}

	for ruleId, want := range expected.AuthenticationContextRules {
		got, ok := actual.AuthenticationContextRules[ruleId]
		if !ok || !got.IsEnabled.Equal(want.IsEnabled) || !got.ClaimValue.Equal(want.ClaimValue) {
			return false
		}
	}

	return true
}

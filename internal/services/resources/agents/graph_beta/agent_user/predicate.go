package graphBetaAgentUser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// agentUserConsistencyPredicate returns a consistency predicate for ReadWithRetry that verifies
// the agent user resource write has propagated across Microsoft Entra replicas before accepting
// the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// Computed-only fields (mail, user_type, created_date_time, creation_type) are not compared —
// they are set by the API and are not meaningful signals for write propagation.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func agentUserConsistencyPredicate(expected *AgentUserResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AgentUserResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// user_principal_name is always returned by the API — its presence confirms the
		// replica has the full user object.
		if actual.UserPrincipalName.IsNull() || actual.UserPrincipalName.IsUnknown() || actual.UserPrincipalName.ValueString() == "" {
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
		if !expected.MailNickname.IsNull() && !expected.MailNickname.IsUnknown() {
			if !actual.MailNickname.Equal(expected.MailNickname) {
				return false
			}
		}
		if !expected.GivenName.IsNull() && !expected.GivenName.IsUnknown() {
			if !actual.GivenName.Equal(expected.GivenName) {
				return false
			}
		}
		if !expected.Surname.IsNull() && !expected.Surname.IsUnknown() {
			if !actual.Surname.Equal(expected.Surname) {
				return false
			}
		}
		if !expected.JobTitle.IsNull() && !expected.JobTitle.IsUnknown() {
			if !actual.JobTitle.Equal(expected.JobTitle) {
				return false
			}
		}
		if !expected.Department.IsNull() && !expected.Department.IsUnknown() {
			if !actual.Department.Equal(expected.Department) {
				return false
			}
		}
		if !expected.CompanyName.IsNull() && !expected.CompanyName.IsUnknown() {
			if !actual.CompanyName.Equal(expected.CompanyName) {
				return false
			}
		}
		if !expected.OfficeLocation.IsNull() && !expected.OfficeLocation.IsUnknown() {
			if !actual.OfficeLocation.Equal(expected.OfficeLocation) {
				return false
			}
		}
		if !expected.City.IsNull() && !expected.City.IsUnknown() {
			if !actual.City.Equal(expected.City) {
				return false
			}
		}
		if !expected.State.IsNull() && !expected.State.IsUnknown() {
			if !actual.State.Equal(expected.State) {
				return false
			}
		}
		if !expected.Country.IsNull() && !expected.Country.IsUnknown() {
			if !actual.Country.Equal(expected.Country) {
				return false
			}
		}
		if !expected.PostalCode.IsNull() && !expected.PostalCode.IsUnknown() {
			if !actual.PostalCode.Equal(expected.PostalCode) {
				return false
			}
		}
		if !expected.StreetAddress.IsNull() && !expected.StreetAddress.IsUnknown() {
			if !actual.StreetAddress.Equal(expected.StreetAddress) {
				return false
			}
		}
		if !expected.UsageLocation.IsNull() && !expected.UsageLocation.IsUnknown() {
			if !actual.UsageLocation.Equal(expected.UsageLocation) {
				return false
			}
		}
		if !expected.PreferredLanguage.IsNull() && !expected.PreferredLanguage.IsUnknown() {
			if !actual.PreferredLanguage.Equal(expected.PreferredLanguage) {
				return false
			}
		}
		if !expected.SponsorIds.IsNull() && !expected.SponsorIds.IsUnknown() {
			if !actual.SponsorIds.Equal(expected.SponsorIds) {
				return false
			}
		}

		return true
	}
}

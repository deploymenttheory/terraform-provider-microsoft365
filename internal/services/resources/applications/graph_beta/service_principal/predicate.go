package graphBetaServicePrincipal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// servicePrincipalConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the service principal resource write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// app_id identifies the application registration that backs this service principal and is an
// immutable key. service_principal_type is Computed (set by Entra) and is not compared.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func servicePrincipalConsistencyPredicate(expected *ServicePrincipalResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.AppID.ValueString() != expected.AppID.ValueString() {
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
		if !expected.AppRoleAssignmentRequired.IsNull() && !expected.AppRoleAssignmentRequired.IsUnknown() {
			if !actual.AppRoleAssignmentRequired.Equal(expected.AppRoleAssignmentRequired) {
				return false
			}
		}
		if !expected.Description.IsNull() && !expected.Description.IsUnknown() {
			if !actual.Description.Equal(expected.Description) {
				return false
			}
		}
		if !expected.Homepage.IsNull() && !expected.Homepage.IsUnknown() {
			if !actual.Homepage.Equal(expected.Homepage) {
				return false
			}
		}
		if !expected.LoginURL.IsNull() && !expected.LoginURL.IsUnknown() {
			if !actual.LoginURL.Equal(expected.LoginURL) {
				return false
			}
		}
		if !expected.LogoutURL.IsNull() && !expected.LogoutURL.IsUnknown() {
			if !actual.LogoutURL.Equal(expected.LogoutURL) {
				return false
			}
		}
		if !expected.Notes.IsNull() && !expected.Notes.IsUnknown() {
			if !actual.Notes.Equal(expected.Notes) {
				return false
			}
		}
		if !expected.NotificationEmailAddresses.IsNull() && !expected.NotificationEmailAddresses.IsUnknown() {
			if !actual.NotificationEmailAddresses.Equal(expected.NotificationEmailAddresses) {
				return false
			}
		}
		if !expected.PreferredSingleSignOnMode.IsNull() && !expected.PreferredSingleSignOnMode.IsUnknown() {
			if !actual.PreferredSingleSignOnMode.Equal(expected.PreferredSingleSignOnMode) {
				return false
			}
		}
		if !expected.Tags.IsNull() && !expected.Tags.IsUnknown() {
			if !actual.Tags.Equal(expected.Tags) {
				return false
			}
		}

		return true
	}
}

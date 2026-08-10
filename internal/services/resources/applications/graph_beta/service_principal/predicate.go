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
// immutable key. Computed properties are not compared: those set by Entra (such as
// service_principal_type) and those reflected from the associated application (such as
// display_name) are not written by this resource, and the reflected ones can change during
// the same apply when the backing application is updated.
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
		//
		// Properties reflected from the associated application (display_name, homepage,
		// logout_url, and the other Computed reflections) are deliberately not compared:
		// this resource never writes them, and they legitimately change when the backing
		// application is updated in the same apply. Asserting on them would retry until
		// the limit and fail an otherwise successful update.
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
		if !expected.LoginURL.IsNull() && !expected.LoginURL.IsUnknown() {
			if !actual.LoginURL.Equal(expected.LoginURL) {
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
		if !expected.AlternativeNames.IsNull() && !expected.AlternativeNames.IsUnknown() {
			if !actual.AlternativeNames.Equal(expected.AlternativeNames) {
				return false
			}
		}
		// Compared even when the expected value is null: updates clear the property with an
		// explicit null, so a stale replica still returning the old object must retry.
		if !expected.SamlSingleSignOnSettings.IsUnknown() {
			if !actual.SamlSingleSignOnSettings.Equal(expected.SamlSingleSignOnSettings) {
				return false
			}
		}
		if !expected.TokenEncryptionKeyID.IsUnknown() {
			if !actual.TokenEncryptionKeyID.Equal(expected.TokenEncryptionKeyID) {
				return false
			}
		}

		return true
	}
}

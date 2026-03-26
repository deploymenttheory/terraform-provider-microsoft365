package graphBetaUsersUser

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// userConsistencyPredicate returns a consistency predicate for ReadWithRetry that verifies the
// user resource write has propagated across Microsoft Entra replicas before accepting the read
// as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing at least one field to differ and triggering a retry.
//
// Computed-only fields (created_date_time, deleted_date_time, creation_type, external_user_state,
// external_user_state_change_date_time, mail, on_premises_*, proxy_addresses,
// security_identifier, sign_in_sessions_valid_from_date_time) are not compared.
// password_profile is write-only and never returned by the API — it is not compared.
// hard_delete is a Terraform-only field and is not compared.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func userConsistencyPredicate(expected *UserResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual UserResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// user_principal_name is a key identity field always returned by the API.
		// Its presence confirms the replica returned the full user object.
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
		if !expected.UserPrincipalName.IsNull() && !expected.UserPrincipalName.IsUnknown() {
			if !actual.UserPrincipalName.Equal(expected.UserPrincipalName) {
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
		if !expected.MobilePhone.IsNull() && !expected.MobilePhone.IsUnknown() {
			if !actual.MobilePhone.Equal(expected.MobilePhone) {
				return false
			}
		}
		if !expected.BusinessPhones.IsNull() && !expected.BusinessPhones.IsUnknown() {
			if !actual.BusinessPhones.Equal(expected.BusinessPhones) {
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
		if !expected.AgeGroup.IsNull() && !expected.AgeGroup.IsUnknown() {
			if !actual.AgeGroup.Equal(expected.AgeGroup) {
				return false
			}
		}
		if !expected.ManagerId.IsNull() && !expected.ManagerId.IsUnknown() {
			if !actual.ManagerId.Equal(expected.ManagerId) {
				return false
			}
		}
		if !expected.EmployeeHireDate.IsNull() && !expected.EmployeeHireDate.IsUnknown() {
			if !actual.EmployeeHireDate.Equal(expected.EmployeeHireDate) {
				return false
			}
		}
		if !expected.EmployeeId.IsNull() && !expected.EmployeeId.IsUnknown() {
			if !actual.EmployeeId.Equal(expected.EmployeeId) {
				return false
			}
		}
		if !expected.EmployeeType.IsNull() && !expected.EmployeeType.IsUnknown() {
			if !actual.EmployeeType.Equal(expected.EmployeeType) {
				return false
			}
		}
		if !expected.FaxNumber.IsNull() && !expected.FaxNumber.IsUnknown() {
			if !actual.FaxNumber.Equal(expected.FaxNumber) {
				return false
			}
		}
		if !expected.OtherMails.IsNull() && !expected.OtherMails.IsUnknown() {
			if !actual.OtherMails.Equal(expected.OtherMails) {
				return false
			}
		}
		if !expected.PasswordPolicies.IsNull() && !expected.PasswordPolicies.IsUnknown() {
			if !actual.PasswordPolicies.Equal(expected.PasswordPolicies) {
				return false
			}
		}
		if !expected.PreferredDataLocation.IsNull() && !expected.PreferredDataLocation.IsUnknown() {
			if !actual.PreferredDataLocation.Equal(expected.PreferredDataLocation) {
				return false
			}
		}
		if !expected.PreferredName.IsNull() && !expected.PreferredName.IsUnknown() {
			if !actual.PreferredName.Equal(expected.PreferredName) {
				return false
			}
		}
		if !expected.ShowInAddressList.IsNull() && !expected.ShowInAddressList.IsUnknown() {
			if !actual.ShowInAddressList.Equal(expected.ShowInAddressList) {
				return false
			}
		}
		if !expected.UserType.IsNull() && !expected.UserType.IsUnknown() {
			if !actual.UserType.Equal(expected.UserType) {
				return false
			}
		}
		if !expected.ConsentProvidedForMinor.IsNull() && !expected.ConsentProvidedForMinor.IsUnknown() {
			if !actual.ConsentProvidedForMinor.Equal(expected.ConsentProvidedForMinor) {
				return false
			}
		}
		if !expected.OnPremisesImmutableId.IsNull() && !expected.OnPremisesImmutableId.IsUnknown() {
			if !actual.OnPremisesImmutableId.Equal(expected.OnPremisesImmutableId) {
				return false
			}
		}

		return true
	}
}

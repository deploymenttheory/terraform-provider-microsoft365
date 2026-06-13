package graphBetaUser

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructUserItems constructs a list of UserItemModel from a list of Userable
func ConstructUserItems(users []graphmodels.Userable) []UserItemModel {
	if users == nil {
		return []UserItemModel{}
	}

	items := make([]UserItemModel, 0, len(users))
	for _, user := range users {
		if user != nil {
			items = append(items, ConstructUserItem(user))
		}
	}

	return items
}

// ConstructUserItem constructs a UserItemModel from a Userable
func ConstructUserItem(user graphmodels.Userable) UserItemModel {
	return UserItemModel{
		ID:                              convert.GraphToFrameworkString(user.GetId()),
		AboutMe:                         convert.GraphToFrameworkString(user.GetAboutMe()),
		AccountEnabled:                  convert.GraphToFrameworkBool(user.GetAccountEnabled()),
		AgeGroup:                        convert.GraphToFrameworkString(user.GetAgeGroup()),
		BusinessPhones:                  convert.GraphToFrameworkStringSlice(user.GetBusinessPhones()),
		City:                            convert.GraphToFrameworkString(user.GetCity()),
		CompanyName:                     convert.GraphToFrameworkString(user.GetCompanyName()),
		ConsentProvidedForMinor:         convert.GraphToFrameworkString(user.GetConsentProvidedForMinor()),
		Country:                         convert.GraphToFrameworkString(user.GetCountry()),
		CreatedDateTime:                 convert.GraphToFrameworkTime(user.GetCreatedDateTime()),
		CreationType:                    convert.GraphToFrameworkString(user.GetCreationType()),
		DeletedDateTime:                 convert.GraphToFrameworkTime(user.GetDeletedDateTime()),
		Department:                      convert.GraphToFrameworkString(user.GetDepartment()),
		DisplayName:                     convert.GraphToFrameworkString(user.GetDisplayName()),
		EmployeeHireDate:                convert.GraphToFrameworkTime(user.GetEmployeeHireDate()),
		EmployeeId:                      convert.GraphToFrameworkString(user.GetEmployeeId()),
		EmployeeType:                    convert.GraphToFrameworkString(user.GetEmployeeType()),
		ExternalUserState:               convert.GraphToFrameworkString(user.GetExternalUserState()),
		ExternalUserStateChangeDateTime: convert.GraphToFrameworkString(user.GetExternalUserStateChangeDateTime()),
		FaxNumber:                       convert.GraphToFrameworkString(user.GetFaxNumber()),
		GivenName:                       convert.GraphToFrameworkString(user.GetGivenName()),
		JobTitle:                        convert.GraphToFrameworkString(user.GetJobTitle()),
		Mail:                            convert.GraphToFrameworkString(user.GetMail()),
		MailNickname:                    convert.GraphToFrameworkString(user.GetMailNickname()),
		MobilePhone:                     convert.GraphToFrameworkString(user.GetMobilePhone()),
		OfficeLocation:                  convert.GraphToFrameworkString(user.GetOfficeLocation()),
		OnPremisesDistinguishedName:     convert.GraphToFrameworkString(user.GetOnPremisesDistinguishedName()),
		OnPremisesDomainName:            convert.GraphToFrameworkString(user.GetOnPremisesDomainName()),
		OnPremisesImmutableId:           convert.GraphToFrameworkString(user.GetOnPremisesImmutableId()),
		OnPremisesLastSyncDateTime:      convert.GraphToFrameworkTime(user.GetOnPremisesLastSyncDateTime()),
		OnPremisesSamAccountName:        convert.GraphToFrameworkString(user.GetOnPremisesSamAccountName()),
		OnPremisesSecurityIdentifier:    convert.GraphToFrameworkString(user.GetOnPremisesSecurityIdentifier()),
		OnPremisesSyncEnabled:           convert.GraphToFrameworkBool(user.GetOnPremisesSyncEnabled()),
		OnPremisesUserPrincipalName:     convert.GraphToFrameworkString(user.GetOnPremisesUserPrincipalName()),
		OtherMails:                      convert.GraphToFrameworkStringSlice(user.GetOtherMails()),
		PasswordPolicies:                convert.GraphToFrameworkString(user.GetPasswordPolicies()),
		PostalCode:                      convert.GraphToFrameworkString(user.GetPostalCode()),
		PreferredDataLocation:           convert.GraphToFrameworkString(user.GetPreferredDataLocation()),
		PreferredLanguage:               convert.GraphToFrameworkString(user.GetPreferredLanguage()),
		PreferredName:                   convert.GraphToFrameworkString(user.GetPreferredName()),
		ProxyAddresses:                  convert.GraphToFrameworkStringSlice(user.GetProxyAddresses()),
		SecurityIdentifier:              convert.GraphToFrameworkString(user.GetSecurityIdentifier()),
		ShowInAddressList:               convert.GraphToFrameworkBool(user.GetShowInAddressList()),
		SignInSessionsValidFromDateTime: convert.GraphToFrameworkTime(user.GetSignInSessionsValidFromDateTime()),
		State:                           convert.GraphToFrameworkString(user.GetState()),
		StreetAddress:                   convert.GraphToFrameworkString(user.GetStreetAddress()),
		Surname:                         convert.GraphToFrameworkString(user.GetSurname()),
		UsageLocation:                   convert.GraphToFrameworkString(user.GetUsageLocation()),
		UserPrincipalName:               convert.GraphToFrameworkString(user.GetUserPrincipalName()),
		UserType:                        convert.GraphToFrameworkString(user.GetUserType()),
	}
}

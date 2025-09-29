package graphBetaUsersUser

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *UserResourceModel) (graphmodels.Userable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	user := graphmodels.NewUser()

	// Set required properties
	convert.FrameworkToGraphString(data.DisplayName, user.SetDisplayName)
	convert.FrameworkToGraphBool(data.AccountEnabled, user.SetAccountEnabled)
	convert.FrameworkToGraphString(data.UserPrincipalName, user.SetUserPrincipalName)

	// Set optional properties
	convert.FrameworkToGraphString(data.AboutMe, user.SetAboutMe)
	convert.FrameworkToGraphString(data.AgeGroup, user.SetAgeGroup)
	convert.FrameworkToGraphString(data.City, user.SetCity)
	convert.FrameworkToGraphString(data.CompanyName, user.SetCompanyName)
	convert.FrameworkToGraphString(data.ConsentProvidedForMinor, user.SetConsentProvidedForMinor)
	convert.FrameworkToGraphString(data.Country, user.SetCountry)
	convert.FrameworkToGraphString(data.Department, user.SetDepartment)

	if err := convert.FrameworkToGraphTime(data.EmployeeHireDate, user.SetEmployeeHireDate); err != nil {
		return nil, fmt.Errorf("error converting employee hire date: %v", err)
	}

	convert.FrameworkToGraphString(data.EmployeeId, user.SetEmployeeId)
	convert.FrameworkToGraphString(data.EmployeeType, user.SetEmployeeType)
	convert.FrameworkToGraphString(data.ExternalUserState, user.SetExternalUserState)
	convert.FrameworkToGraphString(data.ExternalUserStateChangeDateTime, user.SetExternalUserStateChangeDateTime)
	convert.FrameworkToGraphString(data.FaxNumber, user.SetFaxNumber)
	convert.FrameworkToGraphString(data.GivenName, user.SetGivenName)
	convert.FrameworkToGraphString(data.JobTitle, user.SetJobTitle)
	convert.FrameworkToGraphString(data.Mail, user.SetMail)
	convert.FrameworkToGraphString(data.MailNickname, user.SetMailNickname)
	convert.FrameworkToGraphString(data.MobilePhone, user.SetMobilePhone)
	convert.FrameworkToGraphString(data.OfficeLocation, user.SetOfficeLocation)
	convert.FrameworkToGraphString(data.OnPremisesDistinguishedName, user.SetOnPremisesDistinguishedName)
	convert.FrameworkToGraphString(data.OnPremisesDomainName, user.SetOnPremisesDomainName)
	convert.FrameworkToGraphString(data.OnPremisesImmutableId, user.SetOnPremisesImmutableId)

	if err := convert.FrameworkToGraphTime(data.OnPremisesLastSyncDateTime, user.SetOnPremisesLastSyncDateTime); err != nil {
		return nil, fmt.Errorf("error converting on-premises last sync date time: %v", err)
	}

	convert.FrameworkToGraphString(data.OnPremisesSamAccountName, user.SetOnPremisesSamAccountName)
	convert.FrameworkToGraphString(data.OnPremisesSecurityIdentifier, user.SetOnPremisesSecurityIdentifier)
	convert.FrameworkToGraphBool(data.OnPremisesSyncEnabled, user.SetOnPremisesSyncEnabled)
	convert.FrameworkToGraphString(data.OnPremisesUserPrincipalName, user.SetOnPremisesUserPrincipalName)
	convert.FrameworkToGraphString(data.PasswordPolicies, user.SetPasswordPolicies)
	convert.FrameworkToGraphString(data.PostalCode, user.SetPostalCode)
	convert.FrameworkToGraphString(data.PreferredDataLocation, user.SetPreferredDataLocation)
	convert.FrameworkToGraphString(data.PreferredLanguage, user.SetPreferredLanguage)
	convert.FrameworkToGraphString(data.PreferredName, user.SetPreferredName)
	convert.FrameworkToGraphString(data.SecurityIdentifier, user.SetSecurityIdentifier)
	convert.FrameworkToGraphBool(data.ShowInAddressList, user.SetShowInAddressList)

	if err := convert.FrameworkToGraphTime(data.SignInSessionsValidFromDateTime, user.SetSignInSessionsValidFromDateTime); err != nil {
		return nil, fmt.Errorf("error converting sign-in sessions valid from date time: %v", err)
	}

	convert.FrameworkToGraphString(data.State, user.SetState)
	convert.FrameworkToGraphString(data.StreetAddress, user.SetStreetAddress)
	convert.FrameworkToGraphString(data.Surname, user.SetSurname)
	convert.FrameworkToGraphString(data.UsageLocation, user.SetUsageLocation)
	convert.FrameworkToGraphString(data.UserType, user.SetUserType)

	// Set collection properties
	if err := convert.FrameworkToGraphStringSet(ctx, data.BusinessPhones, user.SetBusinessPhones); err != nil {
		return nil, fmt.Errorf("error converting business phones: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ImAddresses, user.SetImAddresses); err != nil {
		return nil, fmt.Errorf("error converting IM addresses: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.OtherMails, user.SetOtherMails); err != nil {
		return nil, fmt.Errorf("error converting other mails: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ProxyAddresses, user.SetProxyAddresses); err != nil {
		return nil, fmt.Errorf("error converting proxy addresses: %v", err)
	}

	if data.PasswordProfile != nil {
		passwordProfile := graphmodels.NewPasswordProfile()

		convert.FrameworkToGraphString(data.PasswordProfile.Password, passwordProfile.SetPassword)
		convert.FrameworkToGraphBool(data.PasswordProfile.ForceChangePasswordNextSignIn, passwordProfile.SetForceChangePasswordNextSignIn)
		convert.FrameworkToGraphBool(data.PasswordProfile.ForceChangePasswordNextSignInWithMfa, passwordProfile.SetForceChangePasswordNextSignInWithMfa)

		user.SetPasswordProfile(passwordProfile)
	}

	if !data.Identities.IsNull() && !data.Identities.IsUnknown() {
		var identities []graphmodels.ObjectIdentityable

		identityElements := []ObjectIdentity{}
		if diags := data.Identities.ElementsAs(ctx, &identityElements, false); diags.HasError() {
			return nil, fmt.Errorf("error extracting identities: %v", diags)
		}

		for _, identity := range identityElements {
			objectIdentity := graphmodels.NewObjectIdentity()

			convert.FrameworkToGraphString(identity.SignInType, objectIdentity.SetSignInType)
			convert.FrameworkToGraphString(identity.Issuer, objectIdentity.SetIssuer)
			convert.FrameworkToGraphString(identity.IssuerAssignedId, objectIdentity.SetIssuerAssignedId)

			identities = append(identities, objectIdentity)
		}

		user.SetIdentities(identities)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), user); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return user, nil
}

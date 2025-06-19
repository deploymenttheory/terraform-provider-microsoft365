package user

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *UserResourceModel, remoteResource graphmodels.Userable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AboutMe = convert.GraphToFrameworkString(remoteResource.GetAboutMe())
	data.AccountEnabled = convert.GraphToFrameworkBool(remoteResource.GetAccountEnabled())
	data.AgeGroup = convert.GraphToFrameworkString(remoteResource.GetAgeGroup())
	data.City = convert.GraphToFrameworkString(remoteResource.GetCity())
	data.CompanyName = convert.GraphToFrameworkString(remoteResource.GetCompanyName())
	data.ConsentProvidedForMinor = convert.GraphToFrameworkString(remoteResource.GetConsentProvidedForMinor())
	data.Country = convert.GraphToFrameworkString(remoteResource.GetCountry())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.CreationType = convert.GraphToFrameworkString(remoteResource.GetCreationType())
	data.DeletedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeletedDateTime())
	data.Department = convert.GraphToFrameworkString(remoteResource.GetDepartment())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.EmployeeHireDate = convert.GraphToFrameworkTime(remoteResource.GetEmployeeHireDate())
	data.EmployeeId = convert.GraphToFrameworkString(remoteResource.GetEmployeeId())
	data.EmployeeType = convert.GraphToFrameworkString(remoteResource.GetEmployeeType())
	data.ExternalUserState = convert.GraphToFrameworkString(remoteResource.GetExternalUserState())
	data.ExternalUserStateChangeDateTime = convert.GraphToFrameworkString(remoteResource.GetExternalUserStateChangeDateTime())
	data.FaxNumber = convert.GraphToFrameworkString(remoteResource.GetFaxNumber())
	data.GivenName = convert.GraphToFrameworkString(remoteResource.GetGivenName())
	data.JobTitle = convert.GraphToFrameworkString(remoteResource.GetJobTitle())
	data.Mail = convert.GraphToFrameworkString(remoteResource.GetMail())
	data.MailNickname = convert.GraphToFrameworkString(remoteResource.GetMailNickname())
	data.MobilePhone = convert.GraphToFrameworkString(remoteResource.GetMobilePhone())
	data.OfficeLocation = convert.GraphToFrameworkString(remoteResource.GetOfficeLocation())
	data.OnPremisesDistinguishedName = convert.GraphToFrameworkString(remoteResource.GetOnPremisesDistinguishedName())
	data.OnPremisesDomainName = convert.GraphToFrameworkString(remoteResource.GetOnPremisesDomainName())
	data.OnPremisesImmutableId = convert.GraphToFrameworkString(remoteResource.GetOnPremisesImmutableId())
	data.OnPremisesLastSyncDateTime = convert.GraphToFrameworkTime(remoteResource.GetOnPremisesLastSyncDateTime())
	data.OnPremisesSamAccountName = convert.GraphToFrameworkString(remoteResource.GetOnPremisesSamAccountName())
	data.OnPremisesSecurityIdentifier = convert.GraphToFrameworkString(remoteResource.GetOnPremisesSecurityIdentifier())
	data.OnPremisesSyncEnabled = convert.GraphToFrameworkBool(remoteResource.GetOnPremisesSyncEnabled())
	data.OnPremisesUserPrincipalName = convert.GraphToFrameworkString(remoteResource.GetOnPremisesUserPrincipalName())
	data.PasswordPolicies = convert.GraphToFrameworkString(remoteResource.GetPasswordPolicies())
	data.PostalCode = convert.GraphToFrameworkString(remoteResource.GetPostalCode())
	data.PreferredDataLocation = convert.GraphToFrameworkString(remoteResource.GetPreferredDataLocation())
	data.PreferredLanguage = convert.GraphToFrameworkString(remoteResource.GetPreferredLanguage())
	data.PreferredName = convert.GraphToFrameworkString(remoteResource.GetPreferredName())
	data.SecurityIdentifier = convert.GraphToFrameworkString(remoteResource.GetSecurityIdentifier())
	data.ShowInAddressList = convert.GraphToFrameworkBool(remoteResource.GetShowInAddressList())
	data.SignInSessionsValidFromDateTime = convert.GraphToFrameworkTime(remoteResource.GetSignInSessionsValidFromDateTime())
	data.State = convert.GraphToFrameworkString(remoteResource.GetState())
	data.StreetAddress = convert.GraphToFrameworkString(remoteResource.GetStreetAddress())
	data.Surname = convert.GraphToFrameworkString(remoteResource.GetSurname())
	data.UsageLocation = convert.GraphToFrameworkString(remoteResource.GetUsageLocation())
	data.UserPrincipalName = convert.GraphToFrameworkString(remoteResource.GetUserPrincipalName())
	data.UserType = convert.GraphToFrameworkString(remoteResource.GetUserType())

	businessPhones := remoteResource.GetBusinessPhones()
	if businessPhones != nil {
		data.BusinessPhones = convert.GraphToFrameworkStringSet(ctx, businessPhones)
	}

	imAddresses := remoteResource.GetImAddresses()
	if imAddresses != nil {
		data.ImAddresses = convert.GraphToFrameworkStringSet(ctx, imAddresses)
	}

	otherMails := remoteResource.GetOtherMails()
	if otherMails != nil {
		data.OtherMails = convert.GraphToFrameworkStringSet(ctx, otherMails)
	}

	proxyAddresses := remoteResource.GetProxyAddresses()
	if proxyAddresses != nil {
		data.ProxyAddresses = convert.GraphToFrameworkStringSet(ctx, proxyAddresses)
	}

	identities := remoteResource.GetIdentities()
	if identities != nil {
		identityElements := []ObjectIdentity{}
		for _, identity := range identities {
			identityElement := ObjectIdentity{
				SignInType:       convert.GraphToFrameworkString(identity.GetSignInType()),
				Issuer:           convert.GraphToFrameworkString(identity.GetIssuer()),
				IssuerAssignedId: convert.GraphToFrameworkString(identity.GetIssuerAssignedId()),
			}
			identityElements = append(identityElements, identityElement)
		}

		identitySet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"sign_in_type":       types.StringType,
			"issuer":             types.StringType,
			"issuer_assigned_id": types.StringType,
		}}, identityElements)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to convert identities to set", map[string]interface{}{
				"errors": diags.Errors(),
			})
		} else {
			data.Identities = identitySet
		}
	}

	passwordProfile := remoteResource.GetPasswordProfile()
	if passwordProfile != nil {
		if data.PasswordProfile == nil {
			data.PasswordProfile = &PasswordProfile{}
		}

		// We can't read the password back from the API
		// Only set the ForceChangePasswordNextSignIn and ForceChangePasswordNextSignInWithMfa values
		data.PasswordProfile.ForceChangePasswordNextSignIn = convert.GraphToFrameworkBool(passwordProfile.GetForceChangePasswordNextSignIn())
		data.PasswordProfile.ForceChangePasswordNextSignInWithMfa = convert.GraphToFrameworkBool(passwordProfile.GetForceChangePasswordNextSignInWithMfa())
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

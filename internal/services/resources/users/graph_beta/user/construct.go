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

	convert.FrameworkToGraphString(data.DisplayName, user.SetDisplayName)
	convert.FrameworkToGraphBool(data.AccountEnabled, user.SetAccountEnabled)
	convert.FrameworkToGraphString(data.UserPrincipalName, user.SetUserPrincipalName)
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

	if err := convert.FrameworkToGraphStringSet(ctx, data.BusinessPhones, user.SetBusinessPhones); err != nil {
		return nil, fmt.Errorf("error converting business phones: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.OtherMails, user.SetOtherMails); err != nil {
		return nil, fmt.Errorf("error converting other mails: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ProxyAddresses, user.SetProxyAddresses); err != nil {
		return nil, fmt.Errorf("error converting proxy addresses: %v", err)
	}

	// password_profile is write-only - only included on Create, never on Update
	if data.PasswordProfile != nil {
		passwordProfile := graphmodels.NewPasswordProfile()
		convert.FrameworkToGraphString(data.PasswordProfile.Password, passwordProfile.SetPassword)
		convert.FrameworkToGraphBool(data.PasswordProfile.ForceChangePasswordNextSignIn, passwordProfile.SetForceChangePasswordNextSignIn)

		// Handle force_change_password_next_sign_in_with_mfa with default of false
		if !data.PasswordProfile.ForceChangePasswordNextSignInWithMfa.IsNull() && !data.PasswordProfile.ForceChangePasswordNextSignInWithMfa.IsUnknown() {
			convert.FrameworkToGraphBool(data.PasswordProfile.ForceChangePasswordNextSignInWithMfa, passwordProfile.SetForceChangePasswordNextSignInWithMfa)
		} else {
			// Default to false if not specified
			falseValue := false
			passwordProfile.SetForceChangePasswordNextSignInWithMfa(&falseValue)
		}

		user.SetPasswordProfile(passwordProfile)
	}

	if len(data.CustomSecurityAttributes) > 0 {
		customSecurityAttributes := constructCustomSecurityAttributes(ctx, data.CustomSecurityAttributes)
		if customSecurityAttributes != nil {
			user.SetCustomSecurityAttributes(customSecurityAttributes)
		}
	}

	// Set manager using @odata.bind pattern in additionalData
	if !data.ManagerId.IsNull() && !data.ManagerId.IsUnknown() {
		managerId := data.ManagerId.ValueString()
		if managerId != "" {
			additionalData := make(map[string]any)
			additionalData["manager@odata.bind"] = fmt.Sprintf("https://graph.microsoft.com/beta/users/%s", managerId)
			user.SetAdditionalData(additionalData)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), user); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return user, nil
}

// constructCustomSecurityAttributes converts the Terraform custom security attributes to the Graph API format
// The Graph API expects a structure like:
//
//	{
//	  "Engineering": {
//	    "@odata.type": "#microsoft.graph.customSecurityAttributeValue",
//	    "Project@odata.type": "#Collection(String)",
//	    "Project": ["Baker", "Cascade"],
//	    "Certification": true
//	  }
//	}
func constructCustomSecurityAttributes(ctx context.Context, attributeSets []CustomSecurityAttributeSet) graphmodels.CustomSecurityAttributeValueable {
	if len(attributeSets) == 0 {
		return nil
	}

	tflog.Debug(ctx, "Constructing custom security attributes", map[string]any{
		"attributeSetCount": len(attributeSets),
	})

	customSecurityAttributeValue := graphmodels.NewCustomSecurityAttributeValue()
	additionalData := make(map[string]any)

	for _, attrSet := range attributeSets {
		if attrSet.AttributeSet.IsNull() || attrSet.AttributeSet.IsUnknown() {
			continue
		}

		attributeSetName := attrSet.AttributeSet.ValueString()
		attributeSetData := make(map[string]any)
		attributeSetData["@odata.type"] = "#microsoft.graph.customSecurityAttributeValue"

		for _, attr := range attrSet.Attributes {
			if attr.Name.IsNull() || attr.Name.IsUnknown() {
				continue
			}

			attrName := attr.Name.ValueString()

			// Handle single-valued string
			if !attr.StringValue.IsNull() && !attr.StringValue.IsUnknown() {
				attributeSetData[attrName] = attr.StringValue.ValueString()
				continue
			}

			// Handle single-valued integer
			if !attr.IntValue.IsNull() && !attr.IntValue.IsUnknown() {
				attributeSetData[attrName+"@odata.type"] = "#Int32"
				attributeSetData[attrName] = attr.IntValue.ValueInt32()
				continue
			}

			// Handle boolean
			if !attr.BoolValue.IsNull() && !attr.BoolValue.IsUnknown() {
				attributeSetData[attrName] = attr.BoolValue.ValueBool()
				continue
			}

			// Handle multi-valued strings
			if !attr.StringValues.IsNull() && !attr.StringValues.IsUnknown() {
				var stringValues []string
				diags := attr.StringValues.ElementsAs(ctx, &stringValues, false)
				if !diags.HasError() {
					attributeSetData[attrName+"@odata.type"] = "#Collection(String)"
					attributeSetData[attrName] = stringValues
				}
				continue
			}

			// Handle multi-valued integers
			if !attr.IntValues.IsNull() && !attr.IntValues.IsUnknown() {
				var int32Values []int32
				diags := attr.IntValues.ElementsAs(ctx, &int32Values, false)
				if !diags.HasError() {
					attributeSetData[attrName+"@odata.type"] = "#Collection(Int32)"
					attributeSetData[attrName] = int32Values
				}
				continue
			}
		}

		additionalData[attributeSetName] = attributeSetData
	}

	customSecurityAttributeValue.SetAdditionalData(additionalData)

	tflog.Debug(ctx, "Finished constructing custom security attributes", map[string]any{
		"additionalDataKeys": len(additionalData),
	})

	return customSecurityAttributeValue
}

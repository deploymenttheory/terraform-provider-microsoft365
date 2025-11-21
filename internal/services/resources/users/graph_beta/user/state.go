package graphBetaUsersUser

import (
	"context"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the Terraform state
// userWithCSA is an optional second user object containing custom security attributes from a separate API call
// manager is an optional user object representing the user's manager from a separate API call
func MapRemoteStateToTerraform(ctx context.Context, data *UserResourceModel, remoteResource graphmodels.Userable, userWithCSA graphmodels.Userable, manager graphmodels.DirectoryObjectable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
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

	// Map manager if available
	// If manager is nil, preserve the existing state value (manager might exist but we couldn't read it)
	// Only set to null if we had no manager_id in state originally
	if manager != nil {
		data.ManagerId = convert.GraphToFrameworkString(manager.GetId())
	} else if data.ManagerId.IsNull() || data.ManagerId.IsUnknown() {
		// If state already has no manager, keep it as null
		data.ManagerId = types.StringNull()
	}
	// else: preserve existing data.ManagerId value (don't overwrite with null)

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

	// Handle collection fields - initialize with empty sets if null
	data.BusinessPhones = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetBusinessPhones())
	if data.BusinessPhones.IsNull() {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.BusinessPhones = emptySet
	}

	data.OtherMails = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetOtherMails())
	if data.OtherMails.IsNull() {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.OtherMails = emptySet
	}

	data.ProxyAddresses = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetProxyAddresses())
	if data.ProxyAddresses.IsNull() {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		data.ProxyAddresses = emptySet
	}

	// password_profile is entirely write-only - used only for initial user provisioning
	// Do not map any password_profile fields from the API response
	// The framework automatically handles write-only fields by keeping them null in state

	// Handle custom security attributes from the separate API call
	data.CustomSecurityAttributes = mapCustomSecurityAttributes(ctx, userWithCSA)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})
}

// mapCustomSecurityAttributes extracts custom security attributes from a user object
// The Graph API returns a structure like:
//
//	{
//	  "Engineering": {
//	    "@odata.type": "#microsoft.graph.customSecurityAttributeValue",
//	    "Project@odata.type": "#Collection(String)",
//	    "Project": ["Baker", "Cascade"],
//	    "Certification": true
//	  }
//	}
func mapCustomSecurityAttributes(ctx context.Context, userWithCSA graphmodels.Userable) []CustomSecurityAttributeSet {
	if userWithCSA == nil {
		return nil
	}

	// Try GetCustomSecurityAttributes() first
	if csa := userWithCSA.GetCustomSecurityAttributes(); csa != nil {
		if data := csa.GetAdditionalData(); len(data) > 0 {
			return parseCustomSecurityAttributesMap(ctx, data)
		}
	}

	// Fall back to checking AdditionalData directly (SDK may not deserialize automatically)
	if additionalData := userWithCSA.GetAdditionalData(); additionalData != nil {
		if csaData, ok := additionalData["customSecurityAttributes"].(map[string]any); ok {
			return parseCustomSecurityAttributesMap(ctx, csaData)
		}
	}

	return nil
}

// parseCustomSecurityAttributesMap parses the raw map structure into Terraform state
func parseCustomSecurityAttributesMap(ctx context.Context, data map[string]any) []CustomSecurityAttributeSet {
	var attributeSets []CustomSecurityAttributeSet

	for setName, setValue := range data {
		// Skip @odata.type at root level
		if setName == "@odata.type" {
			continue
		}

		setMap, ok := setValue.(map[string]any)
		if !ok {
			continue
		}

		attributeSet := CustomSecurityAttributeSet{
			AttributeSet: types.StringValue(setName),
			Attributes:   parseAttributeItems(ctx, setMap),
		}

		if len(attributeSet.Attributes) > 0 {
			attributeSets = append(attributeSets, attributeSet)
		}
	}

	return attributeSets
}

// parseAttributeItems parses individual attributes within an attribute set
func parseAttributeItems(ctx context.Context, setMap map[string]any) []CustomSecurityAttributeItem {
	var items []CustomSecurityAttributeItem
	processed := make(map[string]bool)

	for attrName, attrValue := range setMap {
		// Skip @odata.type annotations
		if attrName == "@odata.type" || strings.HasSuffix(attrName, "@odata.type") {
			continue
		}

		if processed[attrName] {
			continue
		}
		processed[attrName] = true

		item := CustomSecurityAttributeItem{
			Name:         types.StringValue(attrName),
			StringValue:  types.StringNull(),
			IntValue:     types.Int32Null(),
			BoolValue:    types.BoolNull(),
			StringValues: types.SetNull(types.StringType),
			IntValues:    types.SetNull(types.Int32Type),
		}

		// Get @odata.type annotation to determine collection types
		odataType := getOdataType(setMap, attrName)

		switch odataType {
		case "#Collection(String)":
			if values := parseStringSlice(attrValue); len(values) > 0 {
				if set, diags := types.SetValueFrom(ctx, types.StringType, values); !diags.HasError() {
					item.StringValues = set
				}
			}

		case "#Collection(Int32)":
			if values := parseInt32Slice(attrValue); len(values) > 0 {
				if set, diags := types.SetValueFrom(ctx, types.Int32Type, values); !diags.HasError() {
					item.IntValues = set
				}
			}

		case "#Int32":
			if v := parseInt32Value(attrValue); v != nil {
				item.IntValue = types.Int32Value(*v)
			}

		default:
			// Infer type from the value itself
			switch v := attrValue.(type) {
			case string:
				item.StringValue = types.StringValue(v)
			case *string:
				if v != nil {
					item.StringValue = types.StringValue(*v)
				}
			case bool:
				item.BoolValue = types.BoolValue(v)
			case *bool:
				if v != nil {
					item.BoolValue = types.BoolValue(*v)
				}
			case float64:
				item.IntValue = types.Int32Value(int32(v))
			case *float64:
				if v != nil {
					item.IntValue = types.Int32Value(int32(*v))
				}
			case int32:
				item.IntValue = types.Int32Value(v)
			case *int32:
				if v != nil {
					item.IntValue = types.Int32Value(*v)
				}
			}
		}

		// Only add if a value was set
		if hasValue(item) {
			items = append(items, item)
		}
	}

	return items
}

// getOdataType extracts the @odata.type annotation for an attribute
func getOdataType(setMap map[string]any, attrName string) string {
	raw, ok := setMap[attrName+"@odata.type"]
	if !ok {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return v
	case *string:
		if v != nil {
			return *v
		}
	}
	return ""
}

// parseStringSlice extracts strings from a slice (handles both string and *string)
func parseStringSlice(value any) []string {
	slice, ok := value.([]any)
	if !ok {
		return nil
	}
	var result []string
	for _, v := range slice {
		switch s := v.(type) {
		case string:
			result = append(result, s)
		case *string:
			if s != nil {
				result = append(result, *s)
			}
		}
	}
	return result
}

// parseInt32Slice extracts int32 values from a slice (handles various numeric types)
func parseInt32Slice(value any) []int32 {
	slice, ok := value.([]any)
	if !ok {
		return nil
	}
	var result []int32
	for _, v := range slice {
		if i := parseInt32Value(v); i != nil {
			result = append(result, *i)
		}
	}
	return result
}

// parseInt32Value converts various numeric types to int32
func parseInt32Value(value any) *int32 {
	var result int32
	switch v := value.(type) {
	case float64:
		result = int32(v)
	case *float64:
		if v == nil {
			return nil
		}
		result = int32(*v)
	case int:
		result = int32(v)
	case *int:
		if v == nil {
			return nil
		}
		result = int32(*v)
	case int32:
		result = v
	case *int32:
		if v == nil {
			return nil
		}
		result = *v
	case int64:
		result = int32(v)
	case *int64:
		if v == nil {
			return nil
		}
		result = int32(*v)
	default:
		return nil
	}
	return &result
}

// hasValue checks if any value field in the item is set (not null)
func hasValue(item CustomSecurityAttributeItem) bool {
	return !item.StringValue.IsNull() ||
		!item.IntValue.IsNull() ||
		!item.BoolValue.IsNull() ||
		!item.StringValues.IsNull() ||
		!item.IntValues.IsNull()
}

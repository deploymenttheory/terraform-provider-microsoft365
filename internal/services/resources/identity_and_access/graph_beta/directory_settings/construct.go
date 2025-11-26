package graphBetaDirectorySettings

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Directory setting template types and their corresponding IDs
const (
	TemplateTypeGroupUnifiedGuest         = "Group.Unified.Guest"
	TemplateTypeApplication               = "Application"
	TemplateTypePasswordRuleSettings      = "Password Rule Settings"
	TemplateTypeGroupUnified              = "Group.Unified"
	TemplateTypeProhibitedNamesSettings   = "Prohibited Names Settings"
	TemplateTypeCustomPolicySettings      = "Custom Policy Settings"
	TemplateTypeProhibitedNamesRestricted = "Prohibited Names Restricted Settings"
	TemplateTypeConsentPolicySettings     = "Consent Policy Settings"
)

// templateTypeToID maps template type names to their UUID
var templateTypeToID = map[string]string{
	TemplateTypeGroupUnifiedGuest:         "08d542b9-071f-4e16-94b0-74abb372e3d9",
	TemplateTypeApplication:               "4bc7f740-180e-4586-adb6-38b2e9024e6b",
	TemplateTypePasswordRuleSettings:      "5cf42378-d67d-4f36-ba46-e8b86229381d",
	TemplateTypeGroupUnified:              "62375ab9-6b52-47ed-826b-58e47e0e304b",
	TemplateTypeProhibitedNamesSettings:   "80661d51-be2f-4d46-9713-98a2fcaec5bc",
	TemplateTypeCustomPolicySettings:      "898f1161-d651-43d1-805c-3b0b388a9fc2",
	TemplateTypeProhibitedNamesRestricted: "aad3907d-1d1a-448b-b3ef-7bf7f63db63b",
	TemplateTypeConsentPolicySettings:     "dffd5d46-495d-40a9-8e21-954ff55e198a",
}

// getTemplateID returns the template UUID for a given template type
func getTemplateID(templateType string) string {
	return templateTypeToID[templateType]
}

// constructResource maps the Terraform schema to the SDK model for DirectorySetting.
// This function creates a DirectorySetting object based on the template type specified.
func constructResource(ctx context.Context, data *DirectorySettingsResourceModel) (graphmodels.DirectorySettingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDirectorySetting()

	var values []graphmodels.SettingValueable

	switch data.TemplateType.ValueString() {
	case TemplateTypeGroupUnifiedGuest:
		values = constructGroupUnifiedGuestSettings(data)
	case TemplateTypeApplication:
		values = constructApplicationSettings(data)
	case TemplateTypePasswordRuleSettings:
		values = constructPasswordRuleSettings(data)
	case TemplateTypeGroupUnified:
		values = constructGroupUnifiedSettings(data)
	case TemplateTypeProhibitedNamesSettings:
		values = constructProhibitedNamesSettings(data)
	case TemplateTypeCustomPolicySettings:
		values = constructCustomPolicySettings(data)
	case TemplateTypeProhibitedNamesRestricted:
		values = constructProhibitedNamesRestrictedSettings(data)
	case TemplateTypeConsentPolicySettings:
		values = constructConsentPolicySettings(data)
	default:
		return nil, fmt.Errorf("unsupported template type: %s", data.TemplateType.ValueString())
	}

	requestBody.SetValues(values)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructGroupUnifiedGuestSettings builds settings for Group.Unified.Guest template
func constructGroupUnifiedGuestSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.GroupUnifiedGuestSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("AllowToAddGuests", frameworkBoolToGraphString(data.GroupUnifiedGuestSettings.AllowToAddGuests)),
	}
}

// constructApplicationSettings builds settings for Application template
func constructApplicationSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.ApplicationSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("EnableAccessCheckForPrivilegedApplicationUpdates", frameworkBoolToGraphString(data.ApplicationSettings.EnableAccessCheckForPrivilegedApplicationUpdates)),
	}
}

// constructPasswordRuleSettings builds settings for Password Rule Settings template
func constructPasswordRuleSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.PasswordRuleSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("BannedPasswordCheckOnPremisesMode", frameworkStringToGraphString(data.PasswordRuleSettings.BannedPasswordCheckOnPremisesMode)),
		createSettingValue("EnableBannedPasswordCheckOnPremises", frameworkBoolToGraphString(data.PasswordRuleSettings.EnableBannedPasswordCheckOnPremises)),
		createSettingValue("EnableBannedPasswordCheck", frameworkBoolToGraphString(data.PasswordRuleSettings.EnableBannedPasswordCheck)),
		createSettingValue("LockoutDurationInSeconds", frameworkInt32ToGraphString(data.PasswordRuleSettings.LockoutDurationInSeconds)),
		createSettingValue("LockoutThreshold", frameworkInt32ToGraphString(data.PasswordRuleSettings.LockoutThreshold)),
		createSettingValue("BannedPasswordList", frameworkStringToGraphString(data.PasswordRuleSettings.BannedPasswordList)),
	}
}

// constructGroupUnifiedSettings builds settings for Group.Unified template
func constructGroupUnifiedSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.GroupUnifiedSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("NewUnifiedGroupWritebackDefault", frameworkBoolToGraphString(data.GroupUnifiedSettings.NewUnifiedGroupWritebackDefault)),
		createSettingValue("EnableMIPLabels", frameworkBoolToGraphString(data.GroupUnifiedSettings.EnableMIPLabels)),
		createSettingValue("CustomBlockedWordsList", frameworkStringToGraphString(data.GroupUnifiedSettings.CustomBlockedWordsList)),
		createSettingValue("EnableMSStandardBlockedWords", frameworkBoolToGraphString(data.GroupUnifiedSettings.EnableMSStandardBlockedWords)),
		createSettingValue("ClassificationDescriptions", frameworkStringToGraphString(data.GroupUnifiedSettings.ClassificationDescriptions)),
		createSettingValue("DefaultClassification", frameworkStringToGraphString(data.GroupUnifiedSettings.DefaultClassification)),
		createSettingValue("PrefixSuffixNamingRequirement", frameworkStringToGraphString(data.GroupUnifiedSettings.PrefixSuffixNamingRequirement)),
		createSettingValue("AllowGuestsToBeGroupOwner", frameworkBoolToGraphString(data.GroupUnifiedSettings.AllowGuestsToBeGroupOwner)),
		createSettingValue("AllowGuestsToAccessGroups", frameworkBoolToGraphString(data.GroupUnifiedSettings.AllowGuestsToAccessGroups)),
		createSettingValue("GuestUsageGuidelinesUrl", frameworkStringToGraphString(data.GroupUnifiedSettings.GuestUsageGuidelinesUrl)),
		createSettingValue("GroupCreationAllowedGroupId", frameworkStringToGraphString(data.GroupUnifiedSettings.GroupCreationAllowedGroupId)),
		createSettingValue("AllowToAddGuests", frameworkBoolToGraphString(data.GroupUnifiedSettings.AllowToAddGuests)),
		createSettingValue("UsageGuidelinesUrl", frameworkStringToGraphString(data.GroupUnifiedSettings.UsageGuidelinesUrl)),
		createSettingValue("ClassificationList", frameworkStringToGraphString(data.GroupUnifiedSettings.ClassificationList)),
		createSettingValue("EnableGroupCreation", frameworkBoolToGraphString(data.GroupUnifiedSettings.EnableGroupCreation)),
	}
}

// constructProhibitedNamesSettings builds settings for Prohibited Names Settings template
func constructProhibitedNamesSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.ProhibitedNamesSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("CustomBlockedSubStringsList", frameworkStringToGraphString(data.ProhibitedNamesSettings.CustomBlockedSubStringsList)),
		createSettingValue("CustomBlockedWholeWordsList", frameworkStringToGraphString(data.ProhibitedNamesSettings.CustomBlockedWholeWordsList)),
	}
}

// constructCustomPolicySettings builds settings for Custom Policy Settings template
func constructCustomPolicySettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.CustomPolicySettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("CustomConditionalAccessPolicyUrl", frameworkStringToGraphString(data.CustomPolicySettings.CustomConditionalAccessPolicyUrl)),
	}
}

// constructProhibitedNamesRestrictedSettings builds settings for Prohibited Names Restricted Settings template
func constructProhibitedNamesRestrictedSettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.ProhibitedNamesRestrictedSettings == nil {
		return []graphmodels.SettingValueable{}
	}

	return []graphmodels.SettingValueable{
		createSettingValue("CustomAllowedSubStringsList", frameworkStringToGraphString(data.ProhibitedNamesRestrictedSettings.CustomAllowedSubStringsList)),
		createSettingValue("CustomAllowedWholeWordsList", frameworkStringToGraphString(data.ProhibitedNamesRestrictedSettings.CustomAllowedWholeWordsList)),
		createSettingValue("DoNotValidateAgainstTrademark", frameworkBoolToGraphString(data.ProhibitedNamesRestrictedSettings.DoNotValidateAgainstTrademark)),
	}
}

// constructConsentPolicySettings builds settings for Consent Policy Settings template
func constructConsentPolicySettings(data *DirectorySettingsResourceModel) []graphmodels.SettingValueable {
	if data.ConsentPolicySettings == nil {
		return []graphmodels.SettingValueable{}
	}

	// Note: EnableGroupSpecificConsent and ConstrainGroupSpecificConsentToMembersOfGroupId are
	// intentionally excluded from the request body. The Microsoft Graph API rejects these fields
	// when included in POST/PATCH requests, indicating they are read-only or require special
	// licensing/tenant configuration that is not universally available.
	return []graphmodels.SettingValueable{
		createSettingValue("BlockUserConsentForRiskyApps", frameworkBoolToGraphString(data.ConsentPolicySettings.BlockUserConsentForRiskyApps)),
		createSettingValue("EnableAdminConsentRequests", frameworkBoolToGraphString(data.ConsentPolicySettings.EnableAdminConsentRequests)),
	}
}

// createSettingValue creates a SettingValue object with the given name and value
func createSettingValue(name, value string) graphmodels.SettingValueable {
	settingValue := graphmodels.NewSettingValue()
	settingValue.SetName(&name)
	settingValue.SetValue(&value)
	return settingValue
}

// frameworkStringToGraphString converts a Terraform Framework string to a Graph API string value.
// Returns the string value, or empty string if null/unknown.
func frameworkStringToGraphString(value basetypes.StringValue) string {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueString()
	}
	return ""
}

// frameworkBoolToGraphString converts a Terraform Framework bool to a Graph API string value ("true" or "false").
// Returns "false" if null/unknown.
func frameworkBoolToGraphString(value basetypes.BoolValue) string {
	if !value.IsNull() && !value.IsUnknown() {
		if value.ValueBool() {
			return "true"
		}
		return "false"
	}
	return "false"
}

// frameworkNullableBoolToGraphString converts a Terraform Framework bool to a Graph API string value for nullable boolean fields.
// Returns empty string "" if null/unknown (representing null in Graph API), or "true"/"false" if set.
// This is used for System.Nullable`1[System.Boolean] type fields where the API expects "" instead of "false" for null values.
func frameworkNullableBoolToGraphString(value basetypes.BoolValue) string {
	if !value.IsNull() && !value.IsUnknown() {
		if value.ValueBool() {
			return "true"
		}
		return "false"
	}
	return ""
}

// frameworkInt32ToGraphString converts a Terraform Framework int32 to a Graph API string value.
// Returns "0" if null/unknown.
func frameworkInt32ToGraphString(value basetypes.Int32Value) string {
	if !value.IsNull() && !value.IsUnknown() {
		return strconv.FormatInt(int64(value.ValueInt32()), 10)
	}
	return "0"
}

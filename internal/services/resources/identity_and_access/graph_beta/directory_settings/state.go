package graphBetaDirectorySettings

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the API response to Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	// Set the ID
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	// Map values based on template type
	switch data.TemplateType.ValueString() {
	case TemplateTypeGroupUnifiedGuest:
		mapGroupUnifiedGuestSettings(ctx, data, remoteResource)
	case TemplateTypeApplication:
		mapApplicationSettings(ctx, data, remoteResource)
	case TemplateTypePasswordRuleSettings:
		mapPasswordRuleSettings(ctx, data, remoteResource)
	case TemplateTypeGroupUnified:
		mapGroupUnifiedSettings(ctx, data, remoteResource)
	case TemplateTypeProhibitedNamesSettings:
		mapProhibitedNamesSettings(ctx, data, remoteResource)
	case TemplateTypeCustomPolicySettings:
		mapCustomPolicySettings(ctx, data, remoteResource)
	case TemplateTypeProhibitedNamesRestricted:
		mapProhibitedNamesRestrictedSettings(ctx, data, remoteResource)
	case TemplateTypeConsentPolicySettings:
		mapConsentPolicySettings(ctx, data, remoteResource)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapGroupUnifiedGuestSettings maps Group.Unified.Guest settings
func mapGroupUnifiedGuestSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.GroupUnifiedGuestSettings == nil {
		data.GroupUnifiedGuestSettings = &GroupUnifiedGuestSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "AllowToAddGuests":
			data.GroupUnifiedGuestSettings.AllowToAddGuests = graphStringToBool(value.GetValue())
		}
	}
}

// mapApplicationSettings maps Application settings
func mapApplicationSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.ApplicationSettings == nil {
		data.ApplicationSettings = &ApplicationSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "EnableAccessCheckForPrivilegedApplicationUpdates":
			data.ApplicationSettings.EnableAccessCheckForPrivilegedApplicationUpdates = graphStringToBool(value.GetValue())
		}
	}
}

// mapPasswordRuleSettings maps Password Rule Settings
func mapPasswordRuleSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.PasswordRuleSettings == nil {
		data.PasswordRuleSettings = &PasswordRuleSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "BannedPasswordCheckOnPremisesMode":
			data.PasswordRuleSettings.BannedPasswordCheckOnPremisesMode = convert.GraphToFrameworkString(value.GetValue())
		case "EnableBannedPasswordCheckOnPremises":
			data.PasswordRuleSettings.EnableBannedPasswordCheckOnPremises = graphStringToBool(value.GetValue())
		case "EnableBannedPasswordCheck":
			data.PasswordRuleSettings.EnableBannedPasswordCheck = graphStringToBool(value.GetValue())
		case "LockoutDurationInSeconds":
			data.PasswordRuleSettings.LockoutDurationInSeconds = graphStringToInt32(value.GetValue())
		case "LockoutThreshold":
			data.PasswordRuleSettings.LockoutThreshold = graphStringToInt32(value.GetValue())
		case "BannedPasswordList":
			data.PasswordRuleSettings.BannedPasswordList = convert.GraphToFrameworkString(value.GetValue())
		}
	}
}

// mapGroupUnifiedSettings maps Group.Unified settings
func mapGroupUnifiedSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.GroupUnifiedSettings == nil {
		data.GroupUnifiedSettings = &GroupUnifiedSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "NewUnifiedGroupWritebackDefault":
			data.GroupUnifiedSettings.NewUnifiedGroupWritebackDefault = graphStringToBool(value.GetValue())
		case "EnableMIPLabels":
			data.GroupUnifiedSettings.EnableMIPLabels = graphStringToBool(value.GetValue())
		case "CustomBlockedWordsList":
			data.GroupUnifiedSettings.CustomBlockedWordsList = convert.GraphToFrameworkString(value.GetValue())
		case "EnableMSStandardBlockedWords":
			data.GroupUnifiedSettings.EnableMSStandardBlockedWords = graphStringToBool(value.GetValue())
		case "ClassificationDescriptions":
			data.GroupUnifiedSettings.ClassificationDescriptions = convert.GraphToFrameworkString(value.GetValue())
		case "DefaultClassification":
			data.GroupUnifiedSettings.DefaultClassification = convert.GraphToFrameworkString(value.GetValue())
		case "PrefixSuffixNamingRequirement":
			data.GroupUnifiedSettings.PrefixSuffixNamingRequirement = convert.GraphToFrameworkString(value.GetValue())
		case "AllowGuestsToBeGroupOwner":
			data.GroupUnifiedSettings.AllowGuestsToBeGroupOwner = graphStringToBool(value.GetValue())
		case "AllowGuestsToAccessGroups":
			data.GroupUnifiedSettings.AllowGuestsToAccessGroups = graphStringToBool(value.GetValue())
		case "GuestUsageGuidelinesUrl":
			data.GroupUnifiedSettings.GuestUsageGuidelinesUrl = convert.GraphToFrameworkString(value.GetValue())
		case "GroupCreationAllowedGroupId":
			data.GroupUnifiedSettings.GroupCreationAllowedGroupId = convert.GraphToFrameworkString(value.GetValue())
		case "AllowToAddGuests":
			data.GroupUnifiedSettings.AllowToAddGuests = graphStringToBool(value.GetValue())
		case "UsageGuidelinesUrl":
			data.GroupUnifiedSettings.UsageGuidelinesUrl = convert.GraphToFrameworkString(value.GetValue())
		case "ClassificationList":
			data.GroupUnifiedSettings.ClassificationList = convert.GraphToFrameworkString(value.GetValue())
		case "EnableGroupCreation":
			data.GroupUnifiedSettings.EnableGroupCreation = graphStringToBool(value.GetValue())
		}
	}
}

// mapProhibitedNamesSettings maps Prohibited Names Settings
func mapProhibitedNamesSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.ProhibitedNamesSettings == nil {
		data.ProhibitedNamesSettings = &ProhibitedNamesSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "CustomBlockedSubStringsList":
			data.ProhibitedNamesSettings.CustomBlockedSubStringsList = convert.GraphToFrameworkString(value.GetValue())
		case "CustomBlockedWholeWordsList":
			data.ProhibitedNamesSettings.CustomBlockedWholeWordsList = convert.GraphToFrameworkString(value.GetValue())
		}
	}
}

// mapCustomPolicySettings maps Custom Policy Settings
func mapCustomPolicySettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.CustomPolicySettings == nil {
		data.CustomPolicySettings = &CustomPolicySettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "CustomConditionalAccessPolicyUrl":
			data.CustomPolicySettings.CustomConditionalAccessPolicyUrl = convert.GraphToFrameworkString(value.GetValue())
		}
	}
}

// mapProhibitedNamesRestrictedSettings maps Prohibited Names Restricted Settings
func mapProhibitedNamesRestrictedSettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.ProhibitedNamesRestrictedSettings == nil {
		data.ProhibitedNamesRestrictedSettings = &ProhibitedNamesRestrictedSettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "CustomAllowedSubStringsList":
			data.ProhibitedNamesRestrictedSettings.CustomAllowedSubStringsList = convert.GraphToFrameworkString(value.GetValue())
		case "CustomAllowedWholeWordsList":
			data.ProhibitedNamesRestrictedSettings.CustomAllowedWholeWordsList = convert.GraphToFrameworkString(value.GetValue())
		case "DoNotValidateAgainstTrademark":
			data.ProhibitedNamesRestrictedSettings.DoNotValidateAgainstTrademark = graphStringToBool(value.GetValue())
		}
	}
}

// mapConsentPolicySettings maps Consent Policy Settings
func mapConsentPolicySettings(ctx context.Context, data *DirectorySettingsResourceModel, remoteResource graphmodels.DirectorySettingable) {
	if data.ConsentPolicySettings == nil {
		data.ConsentPolicySettings = &ConsentPolicySettingsModel{}
	}

	for _, value := range remoteResource.GetValues() {
		name := value.GetName()
		if name == nil {
			continue
		}

		switch *name {
		case "EnableGroupSpecificConsent":
			data.ConsentPolicySettings.EnableGroupSpecificConsent = graphStringToNullableBool(value.GetValue())
		case "BlockUserConsentForRiskyApps":
			data.ConsentPolicySettings.BlockUserConsentForRiskyApps = graphStringToBool(value.GetValue())
		case "EnableAdminConsentRequests":
			data.ConsentPolicySettings.EnableAdminConsentRequests = graphStringToBool(value.GetValue())
		case "ConstrainGroupSpecificConsentToMembersOfGroupId":
			data.ConsentPolicySettings.ConstrainGroupSpecificConsentToMembersOfGroupId = convert.GraphToFrameworkString(value.GetValue())
		}
	}
}

// graphStringToBool converts a Graph API string value ("true"/"false") to a Terraform Framework bool.
// Returns BoolValue(false) if the value is nil or not "true".
func graphStringToBool(value *string) types.Bool {
	if value != nil && *value == "true" {
		return types.BoolValue(true)
	}
	return types.BoolValue(false)
}

// graphStringToNullableBool converts a Graph API string value to a Terraform Framework bool for nullable boolean fields.
// Returns BoolNull() if the value is nil or empty string, BoolValue(true) if "true", BoolValue(false) if "false".
// This is used for System.Nullable`1[System.Boolean] type fields where the API returns "" for null values.
func graphStringToNullableBool(value *string) types.Bool {
	if value == nil || *value == "" {
		return types.BoolNull()
	}
	if *value == "true" {
		return types.BoolValue(true)
	}
	return types.BoolValue(false)
}

// graphStringToInt32 converts a Graph API string value to a Terraform Framework int32.
// Returns Int32Value(0) if the value is nil or cannot be parsed.
func graphStringToInt32(value *string) types.Int32 {
	if value != nil {
		if i, err := strconv.ParseInt(*value, 10, 32); err == nil {
			return types.Int32Value(int32(i))
		}
	}
	return types.Int32Value(0)
}

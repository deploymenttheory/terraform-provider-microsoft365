// REF: https://learn.microsoft.com/en-us/entra/identity/users/groups-naming-policy
// REF: https://learn.microsoft.com/en-us/graph/api/resources/directorysetting?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directorysettingtemplate-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directorysetting-get?view=graph-rest-beta&tabs=http
package graphBetaDirectorySettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DirectorySettingsResourceModel is the root model containing all possible template settings
type DirectorySettingsResourceModel struct {
	ID                        types.String   `tfsdk:"id"`
	TemplateType              types.String   `tfsdk:"template_type"`
	OverwriteExistingSettings types.Bool     `tfsdk:"overwrite_existing_settings"`
	Timeouts                  timeouts.Value `tfsdk:"timeouts"`

	// Template-specific settings (only one will be populated based on template_type)
	GroupUnifiedGuestSettings         *GroupUnifiedGuestSettingsModel         `tfsdk:"group_unified_guest"`
	ApplicationSettings               *ApplicationSettingsModel               `tfsdk:"application"`
	PasswordRuleSettings              *PasswordRuleSettingsModel              `tfsdk:"password_rule_settings"`
	GroupUnifiedSettings              *GroupUnifiedSettingsModel              `tfsdk:"group_unified"`
	ProhibitedNamesSettings           *ProhibitedNamesSettingsModel           `tfsdk:"prohibited_names_settings"`
	CustomPolicySettings              *CustomPolicySettingsModel              `tfsdk:"custom_policy_settings"`
	ProhibitedNamesRestrictedSettings *ProhibitedNamesRestrictedSettingsModel `tfsdk:"prohibited_names_restricted_settings"`
	ConsentPolicySettings             *ConsentPolicySettingsModel             `tfsdk:"consent_policy_settings"`
}

// GroupUnifiedGuestSettingsModel - Settings for a specific Unified Group
// Template ID: 08d542b9-071f-4e16-94b0-74abb372e3d9
type GroupUnifiedGuestSettingsModel struct {
	AllowToAddGuests types.Bool `tfsdk:"allow_to_add_guests"`
}

// ApplicationSettingsModel - Settings for managing tenant-wide application behavior
// Template ID: 4bc7f740-180e-4586-adb6-38b2e9024e6b
type ApplicationSettingsModel struct {
	EnableAccessCheckForPrivilegedApplicationUpdates types.Bool `tfsdk:"enable_access_check_for_privileged_application_updates"`
}

// PasswordRuleSettingsModel - Settings for managing tenant-wide password rule settings
// Template ID: 5cf42378-d67d-4f36-ba46-e8b86229381d
type PasswordRuleSettingsModel struct {
	BannedPasswordCheckOnPremisesMode   types.String `tfsdk:"banned_password_check_on_premises_mode"`
	EnableBannedPasswordCheckOnPremises types.Bool   `tfsdk:"enable_banned_password_check_on_premises"`
	EnableBannedPasswordCheck           types.Bool   `tfsdk:"enable_banned_password_check"`
	LockoutDurationInSeconds            types.Int32  `tfsdk:"lockout_duration_in_seconds"`
	LockoutThreshold                    types.Int32  `tfsdk:"lockout_threshold"`
	BannedPasswordList                  types.String `tfsdk:"banned_password_list"`
}

// GroupUnifiedSettingsModel - Settings for Unified Groups
// Template ID: 62375ab9-6b52-47ed-826b-58e47e0e304b
type GroupUnifiedSettingsModel struct {
	NewUnifiedGroupWritebackDefault types.Bool   `tfsdk:"new_unified_group_writeback_default"`
	EnableMIPLabels                 types.Bool   `tfsdk:"enable_mip_labels"`
	CustomBlockedWordsList          types.String `tfsdk:"custom_blocked_words_list"`
	EnableMSStandardBlockedWords    types.Bool   `tfsdk:"enable_ms_standard_blocked_words"`
	ClassificationDescriptions      types.String `tfsdk:"classification_descriptions"`
	DefaultClassification           types.String `tfsdk:"default_classification"`
	PrefixSuffixNamingRequirement   types.String `tfsdk:"prefix_suffix_naming_requirement"`
	AllowGuestsToBeGroupOwner       types.Bool   `tfsdk:"allow_guests_to_be_group_owner"`
	AllowGuestsToAccessGroups       types.Bool   `tfsdk:"allow_guests_to_access_groups"`
	GuestUsageGuidelinesUrl         types.String `tfsdk:"guest_usage_guidelines_url"`
	GroupCreationAllowedGroupId     types.String `tfsdk:"group_creation_allowed_group_id"`
	AllowToAddGuests                types.Bool   `tfsdk:"allow_to_add_guests"`
	UsageGuidelinesUrl              types.String `tfsdk:"usage_guidelines_url"`
	ClassificationList              types.String `tfsdk:"classification_list"`
	EnableGroupCreation             types.Bool   `tfsdk:"enable_group_creation"`
}

// ProhibitedNamesSettingsModel - Settings for managing tenant-wide prohibited names settings
// Template ID: 80661d51-be2f-4d46-9713-98a2fcaec5bc
type ProhibitedNamesSettingsModel struct {
	CustomBlockedSubStringsList types.String `tfsdk:"custom_blocked_sub_strings_list"`
	CustomBlockedWholeWordsList types.String `tfsdk:"custom_blocked_whole_words_list"`
}

// CustomPolicySettingsModel - Settings for managing tenant-wide custom policy settings
// Template ID: 898f1161-d651-43d1-805c-3b0b388a9fc2
type CustomPolicySettingsModel struct {
	CustomConditionalAccessPolicyUrl types.String `tfsdk:"custom_conditional_access_policy_url"`
}

// ProhibitedNamesRestrictedSettingsModel - Settings for managing tenant-wide prohibited names restricted settings
// Template ID: aad3907d-1d1a-448b-b3ef-7bf7f63db63b
type ProhibitedNamesRestrictedSettingsModel struct {
	CustomAllowedSubStringsList   types.String `tfsdk:"custom_allowed_sub_strings_list"`
	CustomAllowedWholeWordsList   types.String `tfsdk:"custom_allowed_whole_words_list"`
	DoNotValidateAgainstTrademark types.Bool   `tfsdk:"do_not_validate_against_trademark"`
}

// ConsentPolicySettingsModel - Settings for managing tenant-wide consent policy
// Template ID: dffd5d46-495d-40a9-8e21-954ff55e198a
type ConsentPolicySettingsModel struct {
	EnableGroupSpecificConsent                      types.Bool   `tfsdk:"enable_group_specific_consent"`
	BlockUserConsentForRiskyApps                    types.Bool   `tfsdk:"block_user_consent_for_risky_apps"`
	EnableAdminConsentRequests                      types.Bool   `tfsdk:"enable_admin_consent_requests"`
	ConstrainGroupSpecificConsentToMembersOfGroupId types.String `tfsdk:"constrain_group_specific_consent_to_members_of_group_id"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depenrollmentbaseprofile?view=graph-rest-beta
package graphBetaMacOSDepEnrollmentProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSDepEnrollmentProfileResourceModel models a macOS Automated Device Enrollment (DEP/ADE)
// enrollment profile (#microsoft.graph.depMacOSEnrollmentProfile) under a DEP onboarding setting.
//
// Endpoint shape: POST /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles
// Fields verified against the Microsoft Graph beta metadata for depMacOSEnrollmentProfile
// (which inherits from depEnrollmentBaseProfile / enrollmentProfile).
type MacOSDepEnrollmentProfileResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	DepOnboardingSettingsID types.String `tfsdk:"dep_onboarding_settings_id"`

	// enrollmentProfile (base)
	DisplayName                                         types.String `tfsdk:"display_name"`
	Description                                         types.String `tfsdk:"description"`
	RequiresUserAuthentication                          types.Bool   `tfsdk:"requires_user_authentication"`
	ConfigurationEndpointURL                            types.String `tfsdk:"configuration_endpoint_url"`
	EnableAuthenticationViaCompanyPortal                types.Bool   `tfsdk:"enable_authentication_via_company_portal"`
	RequireCompanyPortalOnSetupAssistantEnrolledDevices types.Bool   `tfsdk:"require_company_portal_on_setup_assistant_enrolled_devices"`

	// depEnrollmentBaseProfile (inherited)
	IsDefault                     types.Bool   `tfsdk:"is_default"`
	IsMandatory                   types.Bool   `tfsdk:"is_mandatory"`
	SupervisedModeEnabled         types.Bool   `tfsdk:"supervised_mode_enabled"`
	SupportDepartment             types.String `tfsdk:"support_department"`
	SupportPhoneNumber            types.String `tfsdk:"support_phone_number"`
	DeviceNameTemplate            types.String `tfsdk:"device_name_template"`
	ProfileRemovalDisabled        types.Bool   `tfsdk:"profile_removal_disabled"`
	ConfigurationWebURL           types.Bool   `tfsdk:"configuration_web_url"`
	AwaitDeviceConfigured         types.Bool   `tfsdk:"await_device_configured"`
	EnabledSkipKeys               types.Set    `tfsdk:"enabled_skip_keys"`
	EnrollmentTimeAzureAdGroupIds types.Set    `tfsdk:"enrollment_time_azure_ad_group_ids"`

	// Setup Assistant pane skip booleans (inherited from depEnrollmentBaseProfile).
	// These drive the enabledSkipKeys array sent to Graph (single source of truth);
	// enabled_skip_keys is a computed reflection of the keys the provider generates.
	LocationDisabled           types.Bool `tfsdk:"location_disabled"`
	RestoreBlocked             types.Bool `tfsdk:"restore_blocked"`
	AppleIdDisabled            types.Bool `tfsdk:"apple_id_disabled"`
	TermsAndConditionsDisabled types.Bool `tfsdk:"terms_and_conditions_disabled"`
	TouchIdDisabled            types.Bool `tfsdk:"touch_id_disabled"`
	ApplePayDisabled           types.Bool `tfsdk:"apple_pay_disabled"`
	SiriDisabled               types.Bool `tfsdk:"siri_disabled"`
	DiagnosticsDisabled        types.Bool `tfsdk:"diagnostics_disabled"`
	DisplayToneSetupDisabled   types.Bool `tfsdk:"display_tone_setup_disabled"`
	PrivacyPaneDisabled        types.Bool `tfsdk:"privacy_pane_disabled"`
	ScreenTimeScreenDisabled   types.Bool `tfsdk:"screen_time_screen_disabled"`

	// macOS-specific (depMacOSEnrollmentProfile)
	WelcomeScreenDisabled              types.Bool `tfsdk:"welcome_screen_disabled"`
	RegistrationDisabled               types.Bool `tfsdk:"registration_disabled"`
	FileVaultDisabled                  types.Bool `tfsdk:"file_vault_disabled"`
	ICloudDiagnosticsDisabled          types.Bool `tfsdk:"icloud_diagnostics_disabled"`
	PassCodeDisabled                   types.Bool `tfsdk:"pass_code_disabled"`
	ZoomDisabled                       types.Bool `tfsdk:"zoom_disabled"`
	ICloudStorageDisabled              types.Bool `tfsdk:"icloud_storage_disabled"`
	ChooseYourLockScreenDisabled       types.Bool `tfsdk:"choose_your_lock_screen_disabled"`
	AccessibilityScreenDisabled        types.Bool `tfsdk:"accessibility_screen_disabled"`
	AutoUnlockWithWatchDisabled        types.Bool `tfsdk:"auto_unlock_with_watch_disabled"`
	AutoAdvanceSetupEnabled            types.Bool `tfsdk:"auto_advance_setup_enabled"`
	RequestRequiresNetworkTether       types.Bool `tfsdk:"request_requires_network_tether"`
	UsePlatformSSODuringSetupAssistant types.Bool `tfsdk:"use_platform_sso_during_setup_assistant"`

	// Primary (managed local) account auto-creation
	SkipPrimarySetupAccountCreation     types.Bool   `tfsdk:"skip_primary_setup_account_creation"`
	SetPrimarySetupAccountAsRegularUser types.Bool   `tfsdk:"set_primary_setup_account_as_regular_user"`
	DontAutoPopulatePrimaryAccountInfo  types.Bool   `tfsdk:"dont_auto_populate_primary_account_info"`
	PrimaryAccountFullName              types.String `tfsdk:"primary_account_full_name"`
	PrimaryAccountUserName              types.String `tfsdk:"primary_account_user_name"`
	EnableRestrictEditing               types.Bool   `tfsdk:"enable_restrict_editing"`

	// Admin (local) account auto-creation
	AdminAccountUserName types.String `tfsdk:"admin_account_user_name"`
	AdminAccountFullName types.String `tfsdk:"admin_account_full_name"`
	AdminAccountPassword types.String `tfsdk:"admin_account_password"`
	HideAdminAccount     types.Bool   `tfsdk:"hide_admin_account"`

	// Nested admin password auto-rotation settings
	AdminAccountPasswordRotation *AdminAccountPasswordRotationModel `tfsdk:"admin_account_password_rotation"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// AdminAccountPasswordRotationModel models depProfileAdminAccountPasswordRotationSetting.
type AdminAccountPasswordRotationModel struct {
	AutoRotationPeriodInDays                  types.Int32 `tfsdk:"auto_rotation_period_in_days"`
	OnRetrievalAutoRotatePasswordEnabled      types.Bool  `tfsdk:"on_retrieval_auto_rotate_password_enabled"`
	OnRetrievalDelayAutoRotatePasswordInHours types.Int32 `tfsdk:"on_retrieval_delay_auto_rotate_password_in_hours"`
}

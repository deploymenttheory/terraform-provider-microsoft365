// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentconfiguration?view=graph-rest-beta
package graphBetaDeviceEnrollmentConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceEnrollmentConfigurationResourceModel struct {
	ID                                types.String                            `tfsdk:"id"`
	DisplayName                       types.String                            `tfsdk:"display_name"`
	Description                       types.String                            `tfsdk:"description"`
	Priority                          types.Int32                             `tfsdk:"priority"`
	CreatedDateTime                   types.String                            `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String                            `tfsdk:"last_modified_date_time"`
	Version                           types.Int32                             `tfsdk:"version"`
	DeviceEnrollmentConfigurationType types.String                            `tfsdk:"device_enrollment_configuration_type"`
	RoleScopeTagIds                   types.Set                               `tfsdk:"role_scope_tag_ids"`
	NewPlatformRestriction            *NewPlatformRestrictionModel            `tfsdk:"new_platform_restriction"`
	PlatformRestriction               *PlatformRestrictionModel               `tfsdk:"platform_restriction"`
	Windows10EnrollmentCompletionPage *Windows10EnrollmentCompletionPageModel `tfsdk:"windows10_enrollment_completion_page"`
	WindowsHelloForBusiness           *WindowsHelloForBusinessModel           `tfsdk:"windows_hello_for_business"`
	EnrollmentNotifications           *EnrollmentNotificationsModel           `tfsdk:"enrollment_notifications"`
	DeviceComanagementAuthority       *DeviceComanagementAuthorityModel       `tfsdk:"device_comanagement_authority"`

	Assignments []AssignmentResourceModel `tfsdk:"assignment"`
	Timeouts    timeouts.Value            `tfsdk:"timeouts"`
}

type NewPlatformRestrictionModel struct {
	PlatformType types.String                         `tfsdk:"platform_type"`
	Restriction  *DeviceEnrollmentPlatformRestriction `tfsdk:"restriction"`
}

type PlatformRestrictionModel struct {
	AndroidRestriction        *DeviceEnrollmentPlatformRestriction `tfsdk:"android_restriction"`
	AndroidForWorkRestriction *DeviceEnrollmentPlatformRestriction `tfsdk:"android_for_work_restriction"`
	IOSRestriction            *DeviceEnrollmentPlatformRestriction `tfsdk:"ios_restriction"`
	MacRestriction            *DeviceEnrollmentPlatformRestriction `tfsdk:"mac_restriction"`
	MacOSRestriction          *DeviceEnrollmentPlatformRestriction `tfsdk:"macos_restriction"`
	WindowsRestriction        *DeviceEnrollmentPlatformRestriction `tfsdk:"windows_restriction"`
	WindowsMobileRestriction  *DeviceEnrollmentPlatformRestriction `tfsdk:"windows_mobile_restriction"`
	WindowsHomeSkuRestriction *DeviceEnrollmentPlatformRestriction `tfsdk:"windows_home_sku_restriction"`
	TVOSRestriction           *DeviceEnrollmentPlatformRestriction `tfsdk:"tvos_restriction"`
	VisionOSRestriction       *DeviceEnrollmentPlatformRestriction `tfsdk:"vision_os_restriction"`
}

type DeviceEnrollmentPlatformRestriction struct {
	PlatformBlocked                 types.Bool   `tfsdk:"platform_blocked"`
	PersonalDeviceEnrollmentBlocked types.Bool   `tfsdk:"personal_device_enrollment_blocked"`
	OSMinimumVersion                types.String `tfsdk:"os_minimum_version"`
	OSMaximumVersion                types.String `tfsdk:"os_maximum_version"`
	BlockedManufacturers            types.Set    `tfsdk:"blocked_manufacturers"`
	BlockedSkus                     types.Set    `tfsdk:"blocked_skus"`
}

// Windows10EnrollmentCompletionPageModel represents Windows 10 enrollment completion page settings
type Windows10EnrollmentCompletionPageModel struct {
	AllowDeviceResetOnInstallFailure        types.Bool   `tfsdk:"allow_device_reset_on_install_failure"`
	AllowDeviceUseOnInstallFailure          types.Bool   `tfsdk:"allow_device_use_on_install_failure"`
	AllowLogCollectionOnInstallFailure      types.Bool   `tfsdk:"allow_log_collection_on_install_failure"`
	AllowNonBlockingAppInstallation         types.Bool   `tfsdk:"allow_non_blocking_app_installation"`
	BlockDeviceSetupRetryByUser             types.Bool   `tfsdk:"block_device_setup_retry_by_user"`
	CustomErrorMessage                      types.String `tfsdk:"custom_error_message"`
	DisableUserStatusTrackingAfterFirstUser types.Bool   `tfsdk:"disable_user_status_tracking_after_first_user"`
	InstallProgressTimeoutInMinutes         types.Int32  `tfsdk:"install_progress_timeout_in_minutes"`
	InstallQualityUpdates                   types.Bool   `tfsdk:"install_quality_updates"`
	SelectedMobileAppIds                    types.Set    `tfsdk:"selected_mobile_app_ids"`
	ShowInstallationProgress                types.Bool   `tfsdk:"show_installation_progress"`
	TrackInstallProgressForAutopilotOnly    types.Bool   `tfsdk:"track_install_progress_for_autopilot_only"`
}

// WindowsHelloForBusinessModel represents Windows Hello for Business settings
type WindowsHelloForBusinessModel struct {
	State                       types.String `tfsdk:"state"`                          // Possible values: "notConfigured", "enabled", "disabled"
	EnhancedBiometricsState     types.String `tfsdk:"enhanced_biometrics_state"`      // Possible values: "notConfigured", "enabled", "disabled"
	SecurityKeyForSignIn        types.String `tfsdk:"security_key_for_sign_in"`       // Possible values: "notConfigured", "enabled", "disabled"
	PinLowercaseCharactersUsage types.String `tfsdk:"pin_lowercase_characters_usage"` // Possible values: "allowed", "required", "disallowed"
	PinUppercaseCharactersUsage types.String `tfsdk:"pin_uppercase_characters_usage"` // Possible values: "allowed", "required", "disallowed"
	PinSpecialCharactersUsage   types.String `tfsdk:"pin_special_characters_usage"`   // Possible values: "allowed", "required", "disallowed"
	EnhancedSignInSecurity      types.Int32  `tfsdk:"enhanced_sign_in_security"`      // Default is Not Configured
	PinMinimumLength            types.Int32  `tfsdk:"pin_minimum_length"`             // Between 4 and 127, inclusive
	PinMaximumLength            types.Int32  `tfsdk:"pin_maximum_length"`             // Between 4 and 127, inclusive and >= pinMinimumLength
	PinExpirationInDays         types.Int32  `tfsdk:"pin_expiration_in_days"`         // Between 0 and 730, inclusive. 0 = never expire
	PinPreviousBlockCount       types.Int32  `tfsdk:"pin_previous_block_count"`       // Between 0 and 50, inclusive
	RemotePassportEnabled       types.Bool   `tfsdk:"remote_passport_enabled"`        // Controls Remote Windows Hello for Business
	SecurityDeviceRequired      types.Bool   `tfsdk:"security_device_required"`       // Require TPM for Windows Hello for Business
	UnlockWithBiometricsEnabled types.Bool   `tfsdk:"unlock_with_biometrics_enabled"` // Allow biometric gestures (face, fingerprint)
}

type EnrollmentNotificationsModel struct {
	IncludeCompanyPortalLink      types.Bool   `tfsdk:"include_company_portal_link"`
	SendPushNotification          types.Bool   `tfsdk:"send_push_notification"`
	NotificationTitle             types.String `tfsdk:"notification_title"`
	NotificationBody              types.String `tfsdk:"notification_body"`
	NotificationSender            types.String `tfsdk:"notification_sender"`
	DefaultLocale                 types.String `tfsdk:"default_locale"`
	BrandingOptions               types.String `tfsdk:"branding_options"`                 // Enum
	PlatformType                  types.String `tfsdk:"platform_type"`                    // Enum
	TemplateType                  types.String `tfsdk:"template_type"`                    // Enum
	NotificationMessageTemplateId types.String `tfsdk:"notification_message_template_id"` // UUID
	NotificationTemplates         types.Set    `tfsdk:"notification_templates"`           // Set of strings
}

type DeviceComanagementAuthorityModel struct {
	ConfigurationManagerAgentCommandLineArgument types.String `tfsdk:"configuration_manager_agent_command_line_argument"`
	InstallConfigurationManagerAgent             types.Bool   `tfsdk:"install_configuration_manager_agent"`
	ManagedDeviceAuthority                       types.Int32  `tfsdk:"managed_device_authority"`
}

// AssignmentResourceModel defines a single assignment block within the primary resource
type AssignmentResourceModel struct {
	Target   types.String `tfsdk:"target"`    // "include" or "exclude"
	GroupIds types.Set    `tfsdk:"group_ids"` // Set of Microsoft Entra ID group IDs
}

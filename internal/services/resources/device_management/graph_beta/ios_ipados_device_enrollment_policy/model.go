// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios

package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IOSiPadOSDeviceEnrollmentPolicyResourceModel holds the configuration for an iOS/iPadOS Automated
// Device Enrollment (ADE) profile implemented via the settings catalog backed
// `/deviceManagement/configurationPolicies` endpoint. This is the modern replacement for the
// legacy `depIOSEnrollmentProfile` resource.
type IOSiPadOSDeviceEnrollmentPolicyResourceModel struct {
	// Base policy fields from DeviceManagementConfigurationPolicy
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	RoleScopeTagIds      types.Set    `tfsdk:"role_scope_tag_ids"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	SettingsCount        types.Int32  `tfsdk:"settings_count"`
	IsAssigned           types.Bool   `tfsdk:"is_assigned"`
	Platforms            types.String `tfsdk:"platforms"`
	Technologies         types.String `tfsdk:"technologies"`
	TemplateId           types.String `tfsdk:"template_id"`
	TemplateFamily       types.String `tfsdk:"template_family"`

	// DepOnboardingSettingsId is the Apple ABM/ASM DEP token that owns this profile. Auto-resolved
	// on Create if omitted; used to build the `creationSource` field sent to Graph.
	DepOnboardingSettingsId types.String `tfsdk:"dep_onboarding_settings_id"`

	// IsDefaultPolicyAssignment controls whether this policy is the default iOS/iPadOS enrollment
	// profile for its dep_onboarding_settings_id, via the dedicated setDefaultProfile action on
	// /deviceManagement/depOnboardingSettings/{id}/enrollmentProfiles/{enrollmentProfileId}.
	// Reflects the DEP token's actual current default on every Read, regardless of configuration.
	IsDefaultPolicyAssignment types.Bool `tfsdk:"is_default_policy_assignment"`

	// DeviceSecurityGroup is the enrollment time device membership target (the "Device group" tab
	// in the Intune admin center). It is set/cleared via the dedicated
	// setEnrollmentTimeDeviceMembershipTarget/clearEnrollmentTimeDeviceMembershipTarget actions on
	// /deviceManagement/configurationPolicies/{id}, not via the settings catalog.
	DeviceSecurityGroup types.String `tfsdk:"device_security_group"`

	// User affinity / authentication
	RequiresUserAuthentication                    types.Bool `tfsdk:"requires_user_authentication"`
	EnableAuthenticationViaCompanyPortal          types.Bool `tfsdk:"enable_authentication_via_company_portal"`
	RequireSetupAssistantWithModernAuthentication types.Bool `tfsdk:"require_setup_assistant_with_modern_authentication"`
	AwaitFinalConfiguration                       types.Bool `tfsdk:"await_final_configuration"`

	LockedEnrollmentEnabled types.Bool `tfsdk:"locked_enrollment_enabled"`

	// DeviceNameTemplate maps to the ade_devicenametemplatechoices/ade_appledevicenametemplate
	// subtree. When null, device naming is left to the device (choice sent as "_0").
	DeviceNameTemplate types.String `tfsdk:"device_name_template"`

	// CellularDataActivationUrl maps to the ade_activatecellulardatachoices/ade_activatecellulardata
	// subtree. When null, cellular data plan activation is not configured (choice sent as "_0").
	CellularDataActivationUrl types.String `tfsdk:"cellular_data_activation_url"`

	SupportDepartment  types.String `tfsdk:"support_department"`
	SupportPhoneNumber types.String `tfsdk:"support_phone_number"`

	// Setup Assistant screen toggles. Naming mirrors the macos_device_enrollment_policy resource
	// where the same pane exists there, and the legacy depIOSEnrollmentProfile fields elsewhere.
	PasscodeDisabled                      types.Bool `tfsdk:"passcode_disabled"`
	LocationServicesDisabled              types.Bool `tfsdk:"location_services_disabled"`
	RestoreDisabled                       types.Bool `tfsdk:"restore_disabled"`
	AppleIdDisabled                       types.Bool `tfsdk:"apple_id_disabled"`
	TermsAndConditionsDisabled            types.Bool `tfsdk:"terms_and_conditions_disabled"`
	TouchIdDisabled                       types.Bool `tfsdk:"touch_id_disabled"`
	ApplePayDisabled                      types.Bool `tfsdk:"apple_pay_disabled"`
	SiriDisabled                          types.Bool `tfsdk:"siri_disabled"`
	DiagnosticsDisabled                   types.Bool `tfsdk:"diagnostics_disabled"`
	PrivacyPaneDisabled                   types.Bool `tfsdk:"privacy_pane_disabled"`
	RestoreFromAndroidDisabled            types.Bool `tfsdk:"restore_from_android_disabled"`
	IMessageAndFaceTimeDisabled           types.Bool `tfsdk:"imessage_and_facetime_disabled"`
	ScreenTimeScreenDisabled              types.Bool `tfsdk:"screen_time_screen_disabled"`
	SimSetupScreenDisabled                types.Bool `tfsdk:"sim_setup_screen_disabled"`
	SoftwareUpdateScreenDisabled          types.Bool `tfsdk:"software_update_screen_disabled"`
	WatchMigrationScreenDisabled          types.Bool `tfsdk:"watch_migration_screen_disabled"`
	AppearanceScreenDisabled              types.Bool `tfsdk:"appearance_screen_disabled"`
	DeviceToDeviceMigrationDisabled       types.Bool `tfsdk:"device_to_device_migration_disabled"`
	RestoreCompletedScreenDisabled        types.Bool `tfsdk:"restore_completed_screen_disabled"`
	SoftwareUpdateCompletedScreenDisabled types.Bool `tfsdk:"software_update_completed_screen_disabled"`
	GetStartedScreenDisabled              types.Bool `tfsdk:"get_started_screen_disabled"`
	ActionButtonScreenDisabled            types.Bool `tfsdk:"action_button_screen_disabled"`
	SafetyScreenDisabled                  types.Bool `tfsdk:"safety_screen_disabled"`
	TermsOfAddressScreenDisabled          types.Bool `tfsdk:"terms_of_address_screen_disabled"`
	AppleIntelligenceDisabled             types.Bool `tfsdk:"apple_intelligence_disabled"`
	LockdownModeDisabled                  types.Bool `tfsdk:"lockdown_mode_disabled"`
	AppStoreDisabled                      types.Bool `tfsdk:"app_store_disabled"`
	CameraButtonScreenDisabled            types.Bool `tfsdk:"camera_button_screen_disabled"`
	MultitaskingScreenDisabled            types.Bool `tfsdk:"multitasking_screen_disabled"`
	OsShowcaseScreenDisabled              types.Bool `tfsdk:"os_showcase_screen_disabled"`
	SafetyAndHandlingScreenDisabled       types.Bool `tfsdk:"safety_and_handling_screen_disabled"`
	WebContentFilteringDisabled           types.Bool `tfsdk:"web_content_filtering_disabled"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

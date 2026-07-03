// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-macos

package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSDeviceEnrollmentPolicyResourceModel holds the configuration for a macOS Automated Device
// Enrollment (ADE) profile implemented via the settings catalog backed
// `/deviceManagement/configurationPolicies` endpoint. This is the modern replacement for the
// legacy `depMacOSEnrollmentProfile` resource (see macos_dep_enrollment_profile).
type MacOSDeviceEnrollmentPolicyResourceModel struct {
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

	// User affinity / authentication
	RequiresUserAuthentication                          types.Bool `tfsdk:"requires_user_authentication"`
	EnableAuthenticationViaCompanyPortal                types.Bool `tfsdk:"enable_authentication_via_company_portal"`
	RequireCompanyPortalOnSetupAssistantEnrolledDevices types.Bool `tfsdk:"require_company_portal_on_setup_assistant_enrolled_devices"`

	// Await final configuration / local account creation
	AwaitDeviceConfigured types.Bool         `tfsdk:"await_device_configured"`
	AdminAccount          *AdminAccountModel `tfsdk:"admin_account"`

	LockedEnrollmentEnabled types.Bool `tfsdk:"locked_enrollment_enabled"`

	SupportDepartment  types.String `tfsdk:"support_department"`
	SupportPhoneNumber types.String `tfsdk:"support_phone_number"`

	// Setup Assistant screen toggles. Naming mirrors the legacy macos_dep_enrollment_profile
	// resource where the same concept exists there.
	LocationServicesDisabled              types.Bool `tfsdk:"location_services_disabled"`
	RestoreDisabled                       types.Bool `tfsdk:"restore_disabled"`
	AppleIdDisabled                       types.Bool `tfsdk:"apple_id_disabled"`
	TermsAndConditionsDisabled            types.Bool `tfsdk:"terms_and_conditions_disabled"`
	TouchIdDisabled                       types.Bool `tfsdk:"touch_id_disabled"`
	ApplePayDisabled                      types.Bool `tfsdk:"apple_pay_disabled"`
	SiriDisabled                          types.Bool `tfsdk:"siri_disabled"`
	DiagnosticsDisabled                   types.Bool `tfsdk:"diagnostics_disabled"`
	FileVaultDisabled                     types.Bool `tfsdk:"file_vault_disabled"`
	ICloudDiagnosticsDisabled             types.Bool `tfsdk:"icloud_diagnostics_disabled"`
	ICloudStorageDisabled                 types.Bool `tfsdk:"icloud_storage_disabled"`
	DisplayToneSetupDisabled              types.Bool `tfsdk:"display_tone_setup_disabled"`
	ScreenTimeScreenDisabled              types.Bool `tfsdk:"screen_time_screen_disabled"`
	PrivacyPaneDisabled                   types.Bool `tfsdk:"privacy_pane_disabled"`
	AccessibilityScreenDisabled           types.Bool `tfsdk:"accessibility_screen_disabled"`
	AutoUnlockWithWatchDisabled           types.Bool `tfsdk:"auto_unlock_with_watch_disabled"`
	LockdownModeDisabled                  types.Bool `tfsdk:"lockdown_mode_disabled"`
	SoftwareUpdateScreenDisabled          types.Bool `tfsdk:"software_update_screen_disabled"`
	SoftwareUpdateCompletedScreenDisabled types.Bool `tfsdk:"software_update_completed_screen_disabled"`
	TermsOfAddressScreenDisabled          types.Bool `tfsdk:"terms_of_address_screen_disabled"`
	AppleIntelligenceDisabled             types.Bool `tfsdk:"apple_intelligence_disabled"`
	OsShowcaseScreenDisabled              types.Bool `tfsdk:"os_showcase_screen_disabled"`
	AppStoreDisabled                      types.Bool `tfsdk:"app_store_disabled"`

	Assignments types.Set      `tfsdk:"assignments"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

// AdminAccountModel represents the local admin account that ADE creates on the device when
// await_device_configured is true. Maps to the ade_accountsettings_createlocaladmin subtree.
type AdminAccountModel struct {
	CreateLocalAdminAccount   types.Bool           `tfsdk:"create_local_admin_account"`
	UserName                  types.String         `tfsdk:"user_name"`
	FullName                  types.String         `tfsdk:"full_name"`
	HideAccount               types.Bool           `tfsdk:"hide_account"`
	PasswordRotationInDays    types.Int64          `tfsdk:"password_rotation_in_days"`
	CreateLocalPrimaryAccount types.Bool           `tfsdk:"create_local_primary_account"`
	PrimaryAccount            *PrimaryAccountModel `tfsdk:"primary_account"`
}

// PrimaryAccountModel represents the secondary/standard local account that ADE optionally creates
// alongside the local admin account. Maps to the ade_accountsettings_createlocalprimary subtree.
type PrimaryAccountModel struct {
	PrefillAccountInfo types.Bool   `tfsdk:"prefill_account_info"`
	RestrictEditing    types.Bool   `tfsdk:"restrict_editing"`
	FullName           types.String `tfsdk:"full_name"`
	UserName           types.String `tfsdk:"user_name"`
}

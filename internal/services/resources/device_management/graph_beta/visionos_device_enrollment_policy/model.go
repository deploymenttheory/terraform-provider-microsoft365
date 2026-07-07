// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios

package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VisionOSDeviceEnrollmentPolicyResourceModel holds the configuration for a visionOS Automated
// Device Enrollment (ADE) profile implemented via the settings catalog backed
// `/deviceManagement/configurationPolicies` endpoint.
type VisionOSDeviceEnrollmentPolicyResourceModel struct {
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

	// IsDefaultPolicyAssignment controls whether this policy is the default visionOS enrollment
	// profile for its dep_onboarding_settings_id, via the dedicated setDefaultProfile action on
	// /deviceManagement/depOnboardingSettings/{id}/enrollmentProfiles/{enrollmentProfileId}.
	// Reflects the DEP token's actual current default on every Read, regardless of configuration.
	IsDefaultPolicyAssignment types.Bool `tfsdk:"is_default_policy_assignment"`

	// DeviceSecurityGroup is the enrollment time device membership target (the "Device group" tab
	// in the Intune admin center). It is set/cleared via the dedicated
	// setEnrollmentTimeDeviceMembershipTarget/clearEnrollmentTimeDeviceMembershipTarget actions on
	// /deviceManagement/configurationPolicies/{id}, not via the settings catalog.
	DeviceSecurityGroup types.String `tfsdk:"device_security_group"`

	// User affinity / await configuration. visionOS ADE only supports enrollment without user
	// affinity - Graph rejects ade_useraffinitybasic_1 - so RequiresUserAuthentication is
	// Optional/Computed and defaults to false; it is not expected to ever be true.
	UserAffinity          types.Bool `tfsdk:"user_affinity"`
	AwaitDeviceConfigured types.Bool `tfsdk:"await_device_configured"`

	LockedEnrollmentEnabled types.Bool `tfsdk:"locked_enrollment_enabled"`

	SupportDepartment  types.String `tfsdk:"support_department"`
	SupportPhoneNumber types.String `tfsdk:"support_phone_number"`

	// Setup Assistant screen toggles. Naming mirrors the ios_ipados_device_enrollment_policy
	// resource where the same pane exists there.
	AppleIdDisabled              types.Bool `tfsdk:"apple_id_disabled"`
	ApplePayDisabled             types.Bool `tfsdk:"apple_pay_disabled"`
	DiagnosticsDisabled          types.Bool `tfsdk:"diagnostics_disabled"`
	GetStartedScreenDisabled     types.Bool `tfsdk:"get_started_screen_disabled"`
	AppleIntelligenceDisabled    types.Bool `tfsdk:"apple_intelligence_disabled"`
	LocationServicesDisabled     types.Bool `tfsdk:"location_services_disabled"`
	PasscodeDisabled             types.Bool `tfsdk:"passcode_disabled"`
	PrivacyPaneDisabled          types.Bool `tfsdk:"privacy_pane_disabled"`
	ScreenTimeScreenDisabled     types.Bool `tfsdk:"screen_time_screen_disabled"`
	SiriDisabled                 types.Bool `tfsdk:"siri_disabled"`
	SoftwareUpdateScreenDisabled types.Bool `tfsdk:"software_update_screen_disabled"`
	TermsAndConditionsDisabled   types.Bool `tfsdk:"terms_and_conditions_disabled"`
	TipsScreenDisabled           types.Bool `tfsdk:"tips_screen_disabled"`
	TouchIdDisabled              types.Bool `tfsdk:"touch_id_disabled"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

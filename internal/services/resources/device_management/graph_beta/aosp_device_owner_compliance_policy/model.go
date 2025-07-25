// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecompliancepolicy?view=graph-rest-beta
package graphBetaDeviceCompliancePolicies

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	Type                       types.String   `tfsdk:"type"`
	DisplayName                types.String   `tfsdk:"display_name"`
	Description                types.String   `tfsdk:"description"`
	RoleScopeTagIds            types.Set      `tfsdk:"role_scope_tag_ids"`
	OsMinimumVersion           types.String   `tfsdk:"os_minimum_version"`
	OsMaximumVersion           types.String   `tfsdk:"os_maximum_version"`
	PasswordRequired           types.Bool     `tfsdk:"password_required"`
	PasswordRequiredType       types.String   `tfsdk:"password_required_type"`
	ScheduledActionsForRule    types.Set      `tfsdk:"scheduled_actions_for_rule"`
	LocalActions               types.List     `tfsdk:"local_actions"`
	AospDeviceOwnerSettings    types.Object   `tfsdk:"aosp_device_owner_settings"`
	AndroidDeviceOwnerSettings types.Object   `tfsdk:"android_device_owner_settings"`
	IosSettings                types.Object   `tfsdk:"ios_settings"`
	Windows10Settings          types.Object   `tfsdk:"windows10_settings"`
	MacOsSettings              types.Object   `tfsdk:"macos_settings"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}

type ScheduledActionForRuleModel struct {
	RuleName                      types.String `tfsdk:"rule_name"`
	ScheduledActionConfigurations types.List   `tfsdk:"scheduled_action_configurations"`
}

type ScheduledActionConfigurationModel struct {
	ActionType                types.String `tfsdk:"action_type"`
	GracePeriodHours          types.Int32  `tfsdk:"grace_period_hours"`
	NotificationTemplateId    types.String `tfsdk:"notification_template_id"`
	NotificationMessageCcList types.Set    `tfsdk:"notification_message_cc_list"`
}

type AospDeviceOwnerSettingsModel struct {
	MinAndroidSecurityPatchLevel          types.String `tfsdk:"min_android_security_patch_level"`
	SecurityBlockJailbrokenDevices        types.Bool   `tfsdk:"security_block_jailbroken_devices"`
	StorageRequireEncryption              types.Bool   `tfsdk:"storage_require_encryption"`
	PasswordMinimumLength                 types.Int32  `tfsdk:"password_minimum_length"`
	PasswordMinutesOfInactivityBeforeLock types.Int32  `tfsdk:"password_minutes_of_inactivity_before_lock"`
}

type AndroidDeviceOwnerSettingsModel struct {
	MinAndroidSecurityPatchLevel                       types.String `tfsdk:"min_android_security_patch_level"`
	SecurityBlockJailbrokenDevices                     types.Bool   `tfsdk:"security_block_jailbroken_devices"`
	StorageRequireEncryption                           types.Bool   `tfsdk:"storage_require_encryption"`
	PasswordMinimumLength                              types.Int32  `tfsdk:"password_minimum_length"`
	PasswordMinutesOfInactivityBeforeLock              types.Int32  `tfsdk:"password_minutes_of_inactivity_before_lock"`
	DeviceThreatProtectionRequiredSecurityLevel        types.String `tfsdk:"device_threat_protection_required_security_level"`
	AdvancedThreatProtectionRequiredSecurityLevel      types.String `tfsdk:"advanced_threat_protection_required_security_level"`
	PasswordExpirationDays                             types.Int32  `tfsdk:"password_expiration_days"`
	PasswordPreviousPasswordCountToBlock               types.Int32  `tfsdk:"password_previous_password_count_to_block"`
	SecurityRequiredAndroidSafetyNetEvaluationType     types.String `tfsdk:"security_required_android_safety_net_evaluation_type"`
	SecurityRequireIntuneAppIntegrity                  types.Bool   `tfsdk:"security_require_intune_app_integrity"`
	DeviceThreatProtectionEnabled                      types.Bool   `tfsdk:"device_threat_protection_enabled"`
	SecurityRequireSafetyNetAttestationBasicIntegrity  types.Bool   `tfsdk:"security_require_safety_net_attestation_basic_integrity"`
	SecurityRequireSafetyNetAttestationCertifiedDevice types.Bool   `tfsdk:"security_require_safety_net_attestation_certified_device"`
}

type IosSettingsModel struct {
	DeviceThreatProtectionRequiredSecurityLevel    types.String `tfsdk:"device_threat_protection_required_security_level"`
	AdvancedThreatProtectionRequiredSecurityLevel  types.String `tfsdk:"advanced_threat_protection_required_security_level"`
	DeviceThreatProtectionEnabled                  types.Bool   `tfsdk:"device_threat_protection_enabled"`
	PasscodeRequiredType                           types.String `tfsdk:"passcode_required_type"`
	ManagedEmailProfileRequired                    types.Bool   `tfsdk:"managed_email_profile_required"`
	SecurityBlockJailbrokenDevices                 types.Bool   `tfsdk:"security_block_jailbroken_devices"`
	OsMinimumBuildVersion                          types.String `tfsdk:"os_minimum_build_version"`
	OsMaximumBuildVersion                          types.String `tfsdk:"os_maximum_build_version"`
	PasscodeMinimumCharacterSetCount               types.Int32  `tfsdk:"passcode_minimum_character_set_count"`
	PasscodeMinutesOfInactivityBeforeLock          types.Int32  `tfsdk:"passcode_minutes_of_inactivity_before_lock"`
	PasscodeMinutesOfInactivityBeforeScreenTimeout types.Int32  `tfsdk:"passcode_minutes_of_inactivity_before_screen_timeout"`
	PasscodeExpirationDays                         types.Int32  `tfsdk:"passcode_expiration_days"`
	PasscodePreviousPasscodeBlockCount             types.Int32  `tfsdk:"passcode_previous_passcode_block_count"`
	RestrictedApps                                 types.List   `tfsdk:"restricted_apps"`
}

type RestrictedAppModel struct {
	Name  types.String `tfsdk:"name"`
	AppId types.String `tfsdk:"app_id"`
}

type MacOsSettingsModel struct {
	GatekeeperAllowedAppSource            types.String `tfsdk:"gatekeeper_allowed_app_source"`
	SystemIntegrityProtectionEnabled      types.Bool   `tfsdk:"system_integrity_protection_enabled"`
	OsMinimumBuildVersion                 types.String `tfsdk:"os_minimum_build_version"`
	OsMaximumBuildVersion                 types.String `tfsdk:"os_maximum_build_version"`
	PasswordBlockSimple                   types.Bool   `tfsdk:"password_block_simple"`
	PasswordMinimumCharacterSetCount      types.Int32  `tfsdk:"password_minimum_character_set_count"`
	PasswordMinutesOfInactivityBeforeLock types.Int32  `tfsdk:"password_minutes_of_inactivity_before_lock"`
	StorageRequireEncryption              types.Bool   `tfsdk:"storage_require_encryption"`
	FirewallEnabled                       types.Bool   `tfsdk:"firewall_enabled"`
	FirewallBlockAllIncoming              types.Bool   `tfsdk:"firewall_block_all_incoming"`
	FirewallEnableStealthMode             types.Bool   `tfsdk:"firewall_enable_stealth_mode"`
}

type Windows10SettingsModel struct {
	DeviceThreatProtectionRequiredSecurityLevel types.String `tfsdk:"device_threat_protection_required_security_level"`
	DeviceCompliancePolicyScript                types.String `tfsdk:"device_compliance_policy_script"`
	PasswordRequiredType                        types.String `tfsdk:"password_required_type"`
	WslDistributions                            types.List   `tfsdk:"wsl_distributions"`
	PasswordRequired                            types.Bool   `tfsdk:"password_required"`
	PasswordBlockSimple                         types.Bool   `tfsdk:"password_block_simple"`
	PasswordRequiredToUnlockFromIdle            types.Bool   `tfsdk:"password_required_to_unlock_from_idle"`
	StorageRequireEncryption                    types.Bool   `tfsdk:"storage_require_encryption"`
	PasswordMinutesOfInactivityBeforeLock       types.Int32  `tfsdk:"password_minutes_of_inactivity_before_lock"`
	PasswordMinimumCharacterSetCount            types.Int32  `tfsdk:"password_minimum_character_set_count"`
	ActiveFirewallRequired                      types.Bool   `tfsdk:"active_firewall_required"`
	TpmRequired                                 types.Bool   `tfsdk:"tpm_required"`
	AntivirusRequired                           types.Bool   `tfsdk:"antivirus_required"`
	AntiSpywareRequired                         types.Bool   `tfsdk:"anti_spyware_required"`
	DefenderEnabled                             types.Bool   `tfsdk:"defender_enabled"`
	SignatureOutOfDate                          types.Bool   `tfsdk:"signature_out_of_date"`
	RtpEnabled                                  types.Bool   `tfsdk:"rtp_enabled"`
	DefenderVersion                             types.String `tfsdk:"defender_version"`
	ConfigurationManagerComplianceRequired      types.Bool   `tfsdk:"configuration_manager_compliance_required"`
	OsMinimumVersion                            types.String `tfsdk:"os_minimum_version"`
	OsMaximumVersion                            types.String `tfsdk:"os_maximum_version"`
	MobileOsMinimumVersion                      types.String `tfsdk:"mobile_os_minimum_version"`
	MobileOsMaximumVersion                      types.String `tfsdk:"mobile_os_maximum_version"`
	SecureBootEnabled                           types.Bool   `tfsdk:"secure_boot_enabled"`
	BitLockerEnabled                            types.Bool   `tfsdk:"bit_locker_enabled"`
	CodeIntegrityEnabled                        types.Bool   `tfsdk:"code_integrity_enabled"`
	DeviceThreatProtectionEnabled               types.Bool   `tfsdk:"device_threat_protection_enabled"`
}

type WslDistributionModel struct {
	Distribution     types.String `tfsdk:"distribution"`
	MinimumOSVersion types.String `tfsdk:"minimum_os_version"`
	MaximumOSVersion types.String `tfsdk:"maximum_os_version"`
}

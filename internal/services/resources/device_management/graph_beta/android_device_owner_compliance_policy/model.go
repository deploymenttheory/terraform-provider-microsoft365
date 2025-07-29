// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownercompliancepolicy?view=graph-rest-beta
package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                                                 types.String   `tfsdk:"id"`
	DisplayName                                        types.String   `tfsdk:"display_name"`
	Description                                        types.String   `tfsdk:"description"`
	RoleScopeTagIds                                    types.Set      `tfsdk:"role_scope_tag_ids"`
	DeviceThreatProtectionEnabled                      types.Bool     `tfsdk:"device_threat_protection_enabled"`
	DeviceThreatProtectionRequiredSecurityLevel        types.String   `tfsdk:"device_threat_protection_required_security_level"`
	AdvancedThreatProtectionRequiredSecurityLevel      types.String   `tfsdk:"advanced_threat_protection_required_security_level"`
	SecurityBlockJailbrokenDevices                     types.Bool     `tfsdk:"security_block_jailbroken_devices"`
	SecurityRequireSafetyNetAttestationBasicIntegrity  types.Bool     `tfsdk:"security_require_safety_net_attestation_basic_integrity"`
	SecurityRequireSafetyNetAttestationCertifiedDevice types.Bool     `tfsdk:"security_require_safety_net_attestation_certified_device"`
	OsMinimumVersion                                   types.String   `tfsdk:"os_minimum_version"`
	OsMaximumVersion                                   types.String   `tfsdk:"os_maximum_version"`
	MinAndroidSecurityPatchLevel                       types.String   `tfsdk:"min_android_security_patch_level"`
	PasswordRequired                                   types.Bool     `tfsdk:"password_required"`
	PasswordMinimumLength                              types.Int32    `tfsdk:"password_minimum_length"`
	PasswordMinimumLetterCharacters                    types.Int32    `tfsdk:"password_minimum_letter_characters"`
	PasswordMinimumLowerCaseCharacters                 types.Int32    `tfsdk:"password_minimum_lower_case_characters"`
	PasswordMinimumNonLetterCharacters                 types.Int32    `tfsdk:"password_minimum_non_letter_characters"`
	PasswordMinimumNumericCharacters                   types.Int32    `tfsdk:"password_minimum_numeric_characters"`
	PasswordMinimumSymbolCharacters                    types.Int32    `tfsdk:"password_minimum_symbol_characters"`
	PasswordMinimumUpperCaseCharacters                 types.Int32    `tfsdk:"password_minimum_upper_case_characters"`
	PasswordRequiredType                               types.String   `tfsdk:"password_required_type"`
	PasswordMinutesOfInactivityBeforeLock              types.Int32    `tfsdk:"password_minutes_of_inactivity_before_lock"`
	PasswordExpirationDays                             types.Int32    `tfsdk:"password_expiration_days"`
	PasswordPreviousPasswordCountToBlock               types.Int32    `tfsdk:"password_previous_password_count_to_block"`
	StorageRequireEncryption                           types.Bool     `tfsdk:"storage_require_encryption"`
	SecurityRequireIntuneAppIntegrity                  types.Bool     `tfsdk:"security_require_intune_app_integrity"`
	RequireNoPendingSystemUpdates                      types.Bool     `tfsdk:"require_no_pending_system_updates"`
	SecurityRequiredAndroidSafetyNetEvaluationType     types.String   `tfsdk:"security_required_android_safety_net_evaluation_type"`
	ScheduledActionsForRule                            types.List     `tfsdk:"scheduled_actions_for_rule"`
	Assignments                                        types.Set      `tfsdk:"assignments"`
	Timeouts                                           timeouts.Value `tfsdk:"timeouts"`
}

type ScheduledActionForRuleModel struct {
	RuleName                      types.String `tfsdk:"rule_name"`
	ScheduledActionConfigurations types.Set    `tfsdk:"scheduled_action_configurations"`
}

type ScheduledActionConfigurationModel struct {
	ActionType                types.String `tfsdk:"action_type"`
	GracePeriodHours          types.Int32  `tfsdk:"grace_period_hours"`
	NotificationTemplateId    types.String `tfsdk:"notification_template_id"`
	NotificationMessageCcList types.List   `tfsdk:"notification_message_cc_list"`
}

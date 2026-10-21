// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10compliancepolicy?view=graph-rest-beta
//
// CHANGED FROM ORIGINAL: added DeviceCompliancePolicyScript / CustomComplianceRequired
// fields and the DeviceCompliancePolicyScriptModel struct, copied from
// windows_device_compliance_policy/model.go. Confirmed directly against
// msgraph-beta-sdk-go v0.166.0 that graphmodels.MacOSCompliancePolicy already
// has GetDeviceCompliancePolicyScript/SetDeviceCompliancePolicyScript methods
// (identical to Windows10CompliancePolicy) -- this was just never wired up on
// the Terraform side.
package graphBetaMacosDeviceCompliancePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                                             types.String   `tfsdk:"id"`
	DisplayName                                    types.String   `tfsdk:"display_name"`
	Description                                    types.String   `tfsdk:"description"`
	RoleScopeTagIds                                types.Set      `tfsdk:"role_scope_tag_ids"`
	PasswordRequired                               types.Bool     `tfsdk:"password_required"`
	PasswordBlockSimple                            types.Bool     `tfsdk:"password_block_simple"`
	PasswordMinutesOfInactivityBeforeLock          types.Int32    `tfsdk:"password_minutes_of_inactivity_before_lock"`
	PasswordExpirationDays                         types.Int32    `tfsdk:"password_expiration_days"`
	PasswordMinimumLength                          types.Int32    `tfsdk:"password_minimum_length"`
	PasswordMinimumCharacterSetCount               types.Int32    `tfsdk:"password_minimum_character_set_count"`
	PasswordRequiredType                           types.String   `tfsdk:"password_required_type"`
	PasswordPreviousPasswordBlockCount             types.Int32    `tfsdk:"password_previous_password_block_count"`
	OsMinimumVersion                                types.String   `tfsdk:"os_minimum_version"`
	OsMaximumVersion                                types.String   `tfsdk:"os_maximum_version"`
	OsMinimumBuildVersion                           types.String   `tfsdk:"os_minimum_build_version"`
	OsMaximumBuildVersion                           types.String   `tfsdk:"os_maximum_build_version"`
	SystemIntegrityProtectionEnabled                types.Bool     `tfsdk:"system_integrity_protection_enabled"`
	DeviceThreatProtectionEnabled                   types.Bool     `tfsdk:"device_threat_protection_enabled"`
	DeviceThreatProtectionRequiredSecurityLevel     types.String   `tfsdk:"device_threat_protection_required_security_level"`
	AdvancedThreatProtectionRequiredSecurityLevel   types.String   `tfsdk:"advanced_threat_protection_required_security_level"`
	StorageRequireEncryption                        types.Bool     `tfsdk:"storage_require_encryption"`
	GatekeeperAllowedAppSource                      types.String   `tfsdk:"gatekeeper_allowed_app_source"`
	FirewallEnabled                                 types.Bool     `tfsdk:"firewall_enabled"`
	FirewallBlockAllIncoming                        types.Bool     `tfsdk:"firewall_block_all_incoming"`
	FirewallEnableStealthMode                       types.Bool     `tfsdk:"firewall_enable_stealth_mode"`
	// NEW: custom compliance (script + rules), mirrors windows_device_compliance_policy.
	DeviceCompliancePolicyScript types.Object   `tfsdk:"device_compliance_policy_script"`
	CustomComplianceRequired     types.Bool     `tfsdk:"custom_compliance_required"`
	ScheduledActionsForRule      types.List     `tfsdk:"scheduled_actions_for_rule"`
	Assignments                  types.Set      `tfsdk:"assignments"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
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

// NEW: mirrors windows_device_compliance_policy.DeviceCompliancePolicyScriptModel
type DeviceCompliancePolicyScriptModel struct {
	DeviceComplianceScriptId types.String `tfsdk:"device_compliance_script_id"`
	RulesContent             types.String `tfsdk:"rules_content"`
}
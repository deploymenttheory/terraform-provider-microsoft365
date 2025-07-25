// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10compliancepolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-custom-json
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecompliancepolicyscript?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-wsldistributionconfiguration?view=graph-rest-beta
package graphBetaWindowsDeviceCompliancePolicies

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                                          types.String                                                `tfsdk:"id"`
	DisplayName                                 types.String                                                `tfsdk:"display_name"`
	Description                                 types.String                                                `tfsdk:"description"`
	RoleScopeTagIds                             types.Set                                                   `tfsdk:"role_scope_tag_ids"`
	PasswordRequired                            types.Bool                                                  `tfsdk:"password_required"`
	PasswordBlockSimple                         types.Bool                                                  `tfsdk:"password_block_simple"`
	PasswordRequiredToUnlockFromIdle            types.Bool                                                  `tfsdk:"password_required_to_unlock_from_idle"`
	PasswordMinutesOfInactivityBeforeLock       types.Int32                                                 `tfsdk:"password_minutes_of_inactivity_before_lock"`
	PasswordExpirationDays                      types.Int32                                                 `tfsdk:"password_expiration_days"`
	PasswordMinimumLength                       types.Int32                                                 `tfsdk:"password_minimum_length"`
	PasswordMinimumCharacterSetCount            types.Int32                                                 `tfsdk:"password_minimum_character_set_count"`
	PasswordRequiredType                        types.String                                                `tfsdk:"password_required_type"`
	PasswordPreviousPasswordBlockCount          types.Int32                                                 `tfsdk:"password_previous_password_block_count"`
	RequireHealthyDeviceReport                  types.Bool                                                  `tfsdk:"require_healthy_device_report"`
	OsMinimumVersion                            types.String                                                `tfsdk:"os_minimum_version"`
	OsMaximumVersion                            types.String                                                `tfsdk:"os_maximum_version"`
	MobileOsMinimumVersion                      types.String                                                `tfsdk:"mobile_os_minimum_version"`
	MobileOsMaximumVersion                      types.String                                                `tfsdk:"mobile_os_maximum_version"`
	EarlyLaunchAntiMalwareDriverEnabled         types.Bool                                                  `tfsdk:"early_launch_anti_malware_driver_enabled"`
	BitLockerEnabled                            types.Bool                                                  `tfsdk:"bit_locker_enabled"`
	SecureBootEnabled                           types.Bool                                                  `tfsdk:"secure_boot_enabled"`
	CodeIntegrityEnabled                        types.Bool                                                  `tfsdk:"code_integrity_enabled"`
	MemoryIntegrityEnabled                      types.Bool                                                  `tfsdk:"memory_integrity_enabled"`
	KernelDmaProtectionEnabled                  types.Bool                                                  `tfsdk:"kernel_dma_protection_enabled"`
	VirtualizationBasedSecurityEnabled          types.Bool                                                  `tfsdk:"virtualization_based_security_enabled"`
	FirmwareProtectionEnabled                   types.Bool                                                  `tfsdk:"firmware_protection_enabled"`
	StorageRequireEncryption                    types.Bool                                                  `tfsdk:"storage_require_encryption"`
	ActiveFirewallRequired                      types.Bool                                                  `tfsdk:"active_firewall_required"`
	DefenderEnabled                             types.Bool                                                  `tfsdk:"defender_enabled"`
	DefenderVersion                             types.String                                                `tfsdk:"defender_version"`
	SignatureOutOfDate                          types.Bool                                                  `tfsdk:"signature_out_of_date"`
	RtpEnabled                                  types.Bool                                                  `tfsdk:"rtp_enabled"`
	AntivirusRequired                           types.Bool                                                  `tfsdk:"antivirus_required"`
	AntiSpywareRequired                         types.Bool                                                  `tfsdk:"anti_spyware_required"`
	DeviceThreatProtectionEnabled               types.Bool                                                  `tfsdk:"device_threat_protection_enabled"`
	DeviceThreatProtectionRequiredSecurityLevel types.String                                                `tfsdk:"device_threat_protection_required_security_level"`
	ConfigurationManagerComplianceRequired      types.Bool                                                  `tfsdk:"configuration_manager_compliance_required"`
	TpmRequired                                 types.Bool                                                  `tfsdk:"tpm_required"`
	ValidOperatingSystemBuildRanges             types.List                                                  `tfsdk:"valid_operating_system_build_ranges"`
	DeviceCompliancePolicyScript                types.Object                                                `tfsdk:"device_compliance_policy_script"`
	CustomComplianceRequired                    types.Bool                                                  `tfsdk:"custom_compliance_required"`
	WslDistributions                            types.Set                                                   `tfsdk:"wsl_distributions"`
	ScheduledActionsForRule                     types.Set                                                   `tfsdk:"scheduled_actions_for_rule"`
	Assignments                                 *sharedmodels.DeviceManagementScriptAssignmentResourceModel `tfsdk:"assignments"`
	Timeouts                                    timeouts.Value                                              `tfsdk:"timeouts"`
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

type WslDistributionModel struct {
	Distribution     types.String `tfsdk:"distribution"`
	MinimumOSVersion types.String `tfsdk:"minimum_os_version"`
	MaximumOSVersion types.String `tfsdk:"maximum_os_version"`
}

type ValidOperatingSystemBuildRangeModel struct {
	LowOSVersion  types.String `tfsdk:"low_os_version"`
	HighOSVersion types.String `tfsdk:"high_os_version"`
}

type DeviceCompliancePolicyScriptModel struct {
	DeviceComplianceScriptId types.String `tfsdk:"device_compliance_script_id"`
	RulesContent             types.String `tfsdk:"rules_content"`
}

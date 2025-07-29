// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
package graphBetaLinuxDeviceCompliancePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinuxDeviceCompliancePolicyResourceModel defines the schema for a Linux Device Compliance Policy.
type LinuxDeviceCompliancePolicyResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Platforms            types.String `tfsdk:"platforms"`
	Technologies         types.String `tfsdk:"technologies"`
	RoleScopeTagIds      types.Set    `tfsdk:"role_scope_tag_ids"`
	SettingsCount        types.Int32  `tfsdk:"settings_count"`
	IsAssigned           types.Bool   `tfsdk:"is_assigned"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	// Individual Linux Compliance Settings
	DistributionAllowedDistros      types.List     `tfsdk:"distribution_allowed_distros"`
	CustomComplianceRequired        types.Bool     `tfsdk:"custom_compliance_required"`
	CustomComplianceDiscoveryScript types.String   `tfsdk:"custom_compliance_discovery_script"`
	CustomComplianceRules           types.String   `tfsdk:"custom_compliance_rules"`
	DeviceEncryptionRequired        types.Bool     `tfsdk:"device_encryption_required"`
	PasswordPolicyMinimumDigits     types.Int32    `tfsdk:"password_policy_minimum_digits"`
	PasswordPolicyMinimumLength     types.Int32    `tfsdk:"password_policy_minimum_length"`
	PasswordPolicyMinimumLowercase  types.Int32    `tfsdk:"password_policy_minimum_lowercase"`
	PasswordPolicyMinimumSymbols    types.Int32    `tfsdk:"password_policy_minimum_symbols"`
	PasswordPolicyMinimumUppercase  types.Int32    `tfsdk:"password_policy_minimum_uppercase"`
	Assignments                     types.Set      `tfsdk:"assignments"`
	ScheduledActions                types.List     `tfsdk:"scheduled_actions"`
	Timeouts                        timeouts.Value `tfsdk:"timeouts"`
}

// AllowedDistributionModel represents an allowed Linux distribution configuration
type AllowedDistributionModel struct {
	Type           types.String `tfsdk:"type"`
	MinimumVersion types.String `tfsdk:"minimum_version"`
	MaximumVersion types.String `tfsdk:"maximum_version"`
}

// ScheduledActionForRuleModel represents a scheduled action for compliance rules
type ScheduledActionForRuleModel struct {
	RuleName                      types.String `tfsdk:"rule_name"`
	ScheduledActionConfigurations types.Set    `tfsdk:"scheduled_action_configurations"`
}

// ScheduledActionConfigurationModel represents a scheduled action configuration
type ScheduledActionConfigurationModel struct {
	ActionType                types.String `tfsdk:"action_type"`
	GracePeriodHours          types.Int32  `tfsdk:"grace_period_hours"`
	NotificationTemplateId    types.String `tfsdk:"notification_template_id"`
	NotificationMessageCcList types.List   `tfsdk:"notification_message_cc_list"`
}

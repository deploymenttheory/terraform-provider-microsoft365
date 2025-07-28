// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscompliancepolicy?view=graph-rest-beta
package graphBetaIosDeviceCompliancePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                                             types.String   `tfsdk:"id"`
	DisplayName                                    types.String   `tfsdk:"display_name"`
	Description                                    types.String   `tfsdk:"description"`
	RoleScopeTagIds                                types.Set      `tfsdk:"role_scope_tag_ids"`
	PasscodeRequired                               types.Bool     `tfsdk:"passcode_required"`
	PasscodeBlockSimple                            types.Bool     `tfsdk:"passcode_block_simple"`
	PasscodeMinutesOfInactivityBeforeLock          types.Int32    `tfsdk:"passcode_minutes_of_inactivity_before_lock"`
	PasscodeMinutesOfInactivityBeforeScreenTimeout types.Int32    `tfsdk:"passcode_minutes_of_inactivity_before_screen_timeout"`
	PasscodeExpirationDays                         types.Int32    `tfsdk:"passcode_expiration_days"`
	PasscodeMinimumLength                          types.Int32    `tfsdk:"passcode_minimum_length"`
	PasscodeMinimumCharacterSetCount               types.Int32    `tfsdk:"passcode_minimum_character_set_count"`
	PasscodeRequiredType                           types.String   `tfsdk:"passcode_required_type"`
	PasscodePreviousPasscodeBlockCount             types.Int32    `tfsdk:"passcode_previous_passcode_block_count"`
	OsMinimumVersion                               types.String   `tfsdk:"os_minimum_version"`
	OsMaximumVersion                               types.String   `tfsdk:"os_maximum_version"`
	OsMinimumBuildVersion                          types.String   `tfsdk:"os_minimum_build_version"`
	OsMaximumBuildVersion                          types.String   `tfsdk:"os_maximum_build_version"`
	DeviceThreatProtectionEnabled                  types.Bool     `tfsdk:"device_threat_protection_enabled"`
	DeviceThreatProtectionRequiredSecurityLevel    types.String   `tfsdk:"device_threat_protection_required_security_level"`
	AdvancedThreatProtectionRequiredSecurityLevel  types.String   `tfsdk:"advanced_threat_protection_required_security_level"`
	ManagedEmailProfileRequired                    types.Bool     `tfsdk:"managed_email_profile_required"`
	SecurityBlockJailbrokenDevices                 types.Bool     `tfsdk:"security_block_jailbroken_devices"`
	RestrictedApps                                 types.Set      `tfsdk:"restricted_apps"`
	ScheduledActionsForRule                        types.List     `tfsdk:"scheduled_actions_for_rule"`
	Assignments                                    types.Set      `tfsdk:"assignments"`
	Timeouts                                       timeouts.Value `tfsdk:"timeouts"`
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

type RestrictedAppModel struct {
	Name        types.String `tfsdk:"name"`
	Publisher   types.String `tfsdk:"publisher"`
	AppId       types.String `tfsdk:"app_id"`
	AppStoreUrl types.String `tfsdk:"app_store_url"`
}

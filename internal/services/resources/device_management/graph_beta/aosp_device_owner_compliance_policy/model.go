// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-aospdeviceownercompliancepolicy?view=graph-rest-beta
package graphBetaAospDeviceOwnerCompliancePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCompliancePolicyResourceModel struct {
	ID                                    types.String   `tfsdk:"id"`
	DisplayName                           types.String   `tfsdk:"display_name"`
	Description                           types.String   `tfsdk:"description"`
	RoleScopeTagIds                       types.Set      `tfsdk:"role_scope_tag_ids"`
	PasscodeRequired                      types.Bool     `tfsdk:"passcode_required"`
	PasscodeMinimumLength                 types.Int32    `tfsdk:"passcode_minimum_length"`
	PasscodeMinutesOfInactivityBeforeLock types.Int32    `tfsdk:"passcode_minutes_of_inactivity_before_lock"`
	PasscodeRequiredType                  types.String   `tfsdk:"passcode_required_type"`
	OsMinimumVersion                      types.String   `tfsdk:"os_minimum_version"`
	OsMaximumVersion                      types.String   `tfsdk:"os_maximum_version"`
	SecurityBlockJailbrokenDevices        types.Bool     `tfsdk:"security_block_jailbroken_devices"`
	StorageRequireEncryption              types.Bool     `tfsdk:"storage_require_encryption"`
	MinAndroidSecurityPatchLevel          types.String   `tfsdk:"min_android_security_patch_level"`
	ScheduledActionsForRule               types.List     `tfsdk:"scheduled_actions_for_rule"`
	Assignments                           types.Set      `tfsdk:"assignments"`
	Timeouts                              timeouts.Value `tfsdk:"timeouts"`
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

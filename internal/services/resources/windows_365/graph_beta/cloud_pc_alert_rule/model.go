// REF: https://learn.microsoft.com/en-us/graph/api/resources/devicemanagement-alertrule?view=graph-rest-beta
package graphBetaCloudPcAlertRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcAlertRuleResourceModel struct {
	ID                   types.String               `tfsdk:"id"`
	AlertRuleTemplate    types.String               `tfsdk:"alert_rule_template"`
	Description          types.String               `tfsdk:"description"`
	DisplayName          types.String               `tfsdk:"display_name"`
	Enabled              types.Bool                 `tfsdk:"enabled"`
	IsSystemRule         types.Bool                 `tfsdk:"is_system_rule"`
	NotificationChannels []NotificationChannelModel `tfsdk:"notification_channels"`
	Severity             types.String               `tfsdk:"severity"`
	Threshold            *RuleThresholdModel        `tfsdk:"threshold"`
	Conditions           []RuleConditionModel       `tfsdk:"conditions"`
	Timeouts             timeouts.Value             `tfsdk:"timeouts"`
}

type NotificationChannelModel struct {
	NotificationChannelType types.String                `tfsdk:"notification_channel_type"`
	NotificationReceivers   []NotificationReceiverModel `tfsdk:"notification_receivers"`
}

type NotificationReceiverModel struct {
	ContactInformation types.String `tfsdk:"contact_information"`
	Locale             types.String `tfsdk:"locale"`
}

type RuleThresholdModel struct {
	Aggregation types.String `tfsdk:"aggregation"`
	Operator    types.String `tfsdk:"operator"`
	Target      types.Int32  `tfsdk:"target"`
}

type RuleConditionModel struct {
	RelationshipType  types.String `tfsdk:"relationship_type"`
	ConditionCategory types.String `tfsdk:"condition_category"`
	Aggregation       types.String `tfsdk:"aggregation"`
	Operator          types.String `tfsdk:"operator"`
	ThresholdValue    types.String `tfsdk:"threshold_value"`
}

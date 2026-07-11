package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkInternetAccessForwardingPolicyRuleResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	ForwardingPolicyID   types.String   `tfsdk:"forwarding_policy_id"`
	Name                 types.String   `tfsdk:"name"`
	Action               types.String   `tfsdk:"action"`
	RuleType             types.String   `tfsdk:"rule_type"`
	ClientFallbackAction types.String   `tfsdk:"client_fallback_action"`
	Ports                types.Set      `tfsdk:"ports"`
	Protocol             types.String   `tfsdk:"protocol"`
	Destinations         types.List     `tfsdk:"destinations"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type RuleDestinationModel struct {
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
	BeginAddress types.String `tfsdk:"begin_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

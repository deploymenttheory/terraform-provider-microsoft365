package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkThreatIntelligencePolicyRuleResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	ThreatIntelligencePolicyID types.String   `tfsdk:"threat_intelligence_policy_id"`
	Name                       types.String   `tfsdk:"name"`
	Description                types.String   `tfsdk:"description"`
	Action                     types.String   `tfsdk:"action"`
	Priority                   types.Int32    `tfsdk:"priority"`
	Enabled                    types.Bool     `tfsdk:"enabled"`
	Severity                   types.String   `tfsdk:"severity"`
	Status                     types.String   `tfsdk:"status"`
	Destinations               types.List     `tfsdk:"destinations"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}

type ThreatIntelligencePolicyRuleDestinationModel struct {
	Type   types.String `tfsdk:"type"`
	Values types.List   `tfsdk:"values"`
}

type ThreatIntelligencePolicyRuleIdentity struct {
	ID                         string `tfsdk:"id"`
	ThreatIntelligencePolicyID string `tfsdk:"threat_intelligence_policy_id"`
}

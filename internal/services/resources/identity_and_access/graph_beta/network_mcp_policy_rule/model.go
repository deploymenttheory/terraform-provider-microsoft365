package graphBetaNetworkMCPPolicyRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkMCPPolicyRuleResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	MCPPolicyID        types.String   `tfsdk:"mcp_policy_id"`
	Name               types.String   `tfsdk:"name"`
	Description        types.String   `tfsdk:"description"`
	Action             types.String   `tfsdk:"action"`
	Priority           types.Int32    `tfsdk:"priority"`
	Enabled            types.Bool     `tfsdk:"enabled"`
	Status             types.String   `tfsdk:"status"`
	MatchingConditions types.Object   `tfsdk:"matching_conditions"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

type MCPPolicyRuleIdentity struct {
	ID          string `tfsdk:"id"`
	MCPPolicyID string `tfsdk:"mcp_policy_id"`
}

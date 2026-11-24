package graphBetaNetworkMCPPolicyRule

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mcpPolicyRulesEndpoint is the portal-backed Graph collection of an MCP policy's rules; %s is the policy ID.
const mcpPolicyRulesEndpoint = "networkaccess/mcpPolicies/%s/policyRules"

// mcpPolicyRuleEndpoint addresses one MCP policy rule; the first %s is the policy ID and the second the rule ID.
const mcpPolicyRuleEndpoint = "networkaccess/mcpPolicies/%s/policyRules/%s"

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

// mcpPolicyRuleResponse is the JSON shape of an MCP policy rule returned by Graph. The SDK has no model for it (see the
// SDK deviation in crud.go). matchingConditions stays raw so unsupported fields are reported by MapRemoteStateToTerraform.
type mcpPolicyRuleResponse struct {
	ID          *string `json:"id"`
	ODataType   *string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Action      *string `json:"action"`
	Priority    *int64  `json:"priority"`
	Settings    *struct {
		Status *string `json:"status"`
	} `json:"settings"`
	Conditions json.RawMessage `json:"matchingConditions"`
}

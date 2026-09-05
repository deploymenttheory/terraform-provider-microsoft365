package graphBetaNetworkMCPPolicyRule

import "encoding/json"

// MCP is absent from the generated SDK; raw conditions retain unsupported fields for diagnostics.
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

package graphBetaNetworkMCPPolicyRule

import (
	"encoding/json"
	"fmt"
)

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

// decodeSendMCPPolicyRuleResponse validates and decodes a Microsoft Graph response.
func decodeSendMCPPolicyRuleResponse(result any, err error) (*mcpPolicyRuleResponse, error) {
	if err != nil {
		return nil, err
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 {
		return nil, errEmptyResponse
	}
	var rule *mcpPolicyRuleResponse
	if err := json.Unmarshal(raw, &rule); err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidResponse, err)
	}
	if rule == nil {
		return nil, errEmptyResponse
	}
	return rule, nil
}

// decodeUpdateMCPPolicyRuleResponse validates and decodes the Microsoft Graph response.
func decodeUpdateMCPPolicyRuleResponse(transportResult any, err error) error {
	_, err = decodeSendMCPPolicyRuleResponse(transportResult, err)
	return err
}

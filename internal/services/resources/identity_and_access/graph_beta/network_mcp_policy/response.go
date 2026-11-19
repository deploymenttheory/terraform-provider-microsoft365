package graphBetaNetworkMCPPolicy

import (
	"encoding/json"
	"fmt"
)

// Local response model for the portal-backed endpoint, absent from the generated SDK.
type mcpPolicyResponse struct {
	ID                   *string `json:"id"`
	Name                 *string `json:"name"`
	Description          *string `json:"description"`
	Version              *string `json:"version"`
	LastModifiedDateTime *string `json:"lastModifiedDateTime"`
	Settings             *struct {
		DefaultAction *string `json:"defaultAction"`
	} `json:"settings"`
}

// decodeSendMCPPolicyResponse validates and decodes a Microsoft Graph response.
func decodeSendMCPPolicyResponse(result any, err error) (*mcpPolicyResponse, error) {
	if err != nil {
		return nil, err
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 {
		return nil, errEmptyResponse
	}
	var policy *mcpPolicyResponse
	if err := json.Unmarshal(raw, &policy); err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidResponse, err)
	}
	if policy == nil {
		return nil, errEmptyResponse
	}
	return policy, nil
}

// decodeCreateMCPPolicyResponse validates and decodes the Microsoft Graph response.
func decodeCreateMCPPolicyResponse(transportResult any, err error) (*mcpPolicyResponse, error) {
	result, err := decodeSendMCPPolicyResponse(transportResult, err)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

// decodeUpdateMCPPolicyResponse validates and decodes the Microsoft Graph response.
func decodeUpdateMCPPolicyResponse(transportResult any, err error) error {
	_, err = decodeSendMCPPolicyResponse(transportResult, err)
	return err
}

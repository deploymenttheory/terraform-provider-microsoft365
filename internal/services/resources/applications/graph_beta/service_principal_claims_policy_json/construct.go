package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"encoding/json"
	"fmt"
)

func constructResource(value PolicyJSONValue) ([]byte, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, errPolicyKnown
	}
	object, err := decodePolicyJSON([]byte(value.ValueString()))
	if err != nil {
		return nil, err
	}
	stripServerPolicyFields(object)
	content, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("encode claims policy request: %w", err)
	}
	return content, nil
}

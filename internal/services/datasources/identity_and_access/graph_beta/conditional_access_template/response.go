package graphBetaConditionalAccessTemplate

import (
	"encoding/json"
	"errors"
	"fmt"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var errUnexpectedTemplateResponse = errors.New("unexpected template response type")

// templateResponse preserves evolvable values that the generated SDK cannot
// represent. In particular, its grant-control parser turns verifiedID into block,
// and its scenario parser discards aiAgents. Other fields retain SDK parsing.
type templateResponse struct {
	graphmodels.ConditionalAccessTemplateable
	scenarios       *string
	builtInControls []string
}

func (r *templateResponse) UnmarshalJSON(data []byte) error {
	var values struct {
		Scenarios *string `json:"scenarios"`
		Details   *struct {
			GrantControls *struct {
				BuiltInControls []string `json:"builtInControls"`
			} `json:"grantControls"`
		} `json:"details"`
	}
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("decode template enum values: %w", err)
	}
	node, err := jsonserialization.NewJsonParseNodeFactory().GetRootParseNode("application/json", data)
	if err != nil {
		return fmt.Errorf("parse template response: %w", err)
	}
	parsed, err := node.GetObjectValue(graphmodels.CreateConditionalAccessTemplateFromDiscriminatorValue)
	if err != nil {
		return fmt.Errorf("decode template response: %w", err)
	}
	template, ok := parsed.(graphmodels.ConditionalAccessTemplateable)
	if !ok || template == nil {
		return fmt.Errorf("%w: %T", errUnexpectedTemplateResponse, parsed)
	}
	r.ConditionalAccessTemplateable = template
	r.scenarios = values.Scenarios
	r.builtInControls = nil
	if values.Details != nil && values.Details.GrantControls != nil {
		r.builtInControls = values.Details.GrantControls.BuiltInControls
	}
	return nil
}

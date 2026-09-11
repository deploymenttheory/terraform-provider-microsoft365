package graphBetaNetworkThreatIntelligencePolicy

import (
	"context"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// The SDK omits nil strings. Explicit null is required to clear a description.
type threatIntelligencePolicyRequestBody struct {
	*models.ThreatIntelligencePolicy
	clearDescription bool
	changed          bool
}

func (b *threatIntelligencePolicyRequestBody) Serialize(w s.SerializationWriter) error {
	if err := b.ThreatIntelligencePolicy.Serialize(w); err != nil {
		return fmt.Errorf("serialize threat intelligence payload: %w", err)
	}
	if b.clearDescription {
		if err := w.WriteNullValue("description"); err != nil {
			return fmt.Errorf("serialize null description: %w", err)
		}
	}
	return nil
}

func (b *threatIntelligencePolicyRequestBody) hasChanges() bool { return b.changed }

func constructResource(
	ctx context.Context,
	data *NetworkThreatIntelligencePolicyResourceModel,
) (*threatIntelligencePolicyRequestBody, error) {
	return constructUpdateResource(ctx, data, nil)
}

func constructUpdateResource(
	_ context.Context,
	plan, state *NetworkThreatIntelligencePolicyResourceModel,
) (*threatIntelligencePolicyRequestBody, error) {
	b := &threatIntelligencePolicyRequestBody{
		ThreatIntelligencePolicy: models.NewThreatIntelligencePolicy(),
	}
	if state == nil || !plan.Name.Equal(state.Name) {
		convert.FrameworkToGraphString(plan.Name, b.SetName)
		b.changed = true
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		convert.FrameworkToGraphString(plan.Description, b.SetDescription)
		b.clearDescription = state != nil && plan.Description.IsNull()
		b.changed = true
	}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := models.NewThreatIntelligencePolicySettings()
		if err := convert.FrameworkToGraphEnum(
			plan.DefaultAction,
			models.ParseThreatIntelligenceAction,
			settings.SetDefaultAction,
		); err != nil {
			return nil, fmt.Errorf("construct default action: %w", err)
		}
		b.SetSettings(settings)
		b.changed = true
	}
	return b, nil
}

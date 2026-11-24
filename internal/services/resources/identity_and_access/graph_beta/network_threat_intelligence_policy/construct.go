package graphBetaNetworkThreatIntelligencePolicy

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource builds the threat intelligence policy body. Create passes a nil state and sends every field; Update
// passes the prior state and sends only the fields that changed.
func constructResource(_ context.Context, plan, state *NetworkThreatIntelligencePolicyResourceModel) (models.ThreatIntelligencePolicyable, error) {
	body := models.NewThreatIntelligencePolicy()
	if state == nil || !plan.Name.Equal(state.Name) {
		convert.FrameworkToGraphString(plan.Name, body.SetName)
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		convert.FrameworkToGraphString(plan.Description, body.SetDescription)
		if state != nil && plan.Description.IsNull() {
			// The generated setter omits nil strings, so an explicit JSON null clears the description through additional data.
			body.SetAdditionalData(map[string]any{"description": nil})
		}
	}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := models.NewThreatIntelligencePolicySettings()
		if err := convert.FrameworkToGraphEnum(plan.DefaultAction, models.ParseThreatIntelligenceAction, settings.SetDefaultAction); err != nil {
			return nil, fmt.Errorf("construct default action: %w", err)
		}
		body.SetSettings(settings)
	}
	return body, nil
}

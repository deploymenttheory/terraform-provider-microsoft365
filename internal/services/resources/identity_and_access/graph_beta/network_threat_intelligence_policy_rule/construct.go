package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource builds the threat intelligence rule body. Create passes a nil state and sends every field; Update passes
// the prior state and sends only the fields that changed.
func constructResource(ctx context.Context, plan, state *NetworkThreatIntelligencePolicyRuleResourceModel) (models.ThreatIntelligenceRuleable, error) {
	body := models.NewThreatIntelligenceRule()
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
	if state == nil || !plan.Priority.Equal(state.Priority) {
		priority := int64(plan.Priority.ValueInt32())
		body.SetPriority(&priority)
	}
	if state == nil || !plan.Action.Equal(state.Action) {
		if err := convert.FrameworkToGraphEnum(plan.Action, models.ParseThreatIntelligenceAction, body.SetAction); err != nil {
			return nil, fmt.Errorf("construct action: %w", err)
		}
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		status := models.DISABLED_SECURITYRULESTATUS
		if plan.Enabled.ValueBool() {
			status = models.ENABLED_SECURITYRULESTATUS
		}
		settings := models.NewThreatIntelligenceRuleSettings()
		settings.SetStatus(&status)
		body.SetSettings(settings)
	}
	if state == nil || !plan.Severity.Equal(state.Severity) || !plan.Destinations.Equal(state.Destinations) {
		conditions := models.NewThreatIntelligenceMatchingConditions()
		if err := convert.FrameworkToGraphEnum(plan.Severity, models.ParseThreatIntelligenceSeverity, conditions.SetSeverity); err != nil {
			return nil, fmt.Errorf("construct severity: %w", err)
		}
		var groups []ThreatIntelligencePolicyRuleDestinationModel
		if diags := plan.Destinations.ElementsAs(ctx, &groups, false); diags.HasError() {
			return nil, fmt.Errorf("%w: decode destinations: %v", errInvalidRequest, diags)
		}
		destinations := make([]models.ThreatIntelligenceDestinationable, 0, len(groups))
		for _, group := range groups {
			if group.Type.ValueString() != destinationTypeFQDN {
				return nil, fmt.Errorf("%w: unsupported destination type %q", errInvalidRequest, group.Type.ValueString())
			}
			var values []string
			if diags := group.Values.ElementsAs(ctx, &values, false); diags.HasError() {
				return nil, fmt.Errorf("%w: decode destination values: %v", errInvalidRequest, diags)
			}
			destination := models.NewThreatIntelligenceFqdnDestination()
			destination.SetValues(values)
			destinations = append(destinations, destination)
		}
		conditions.SetDestinations(destinations)
		body.SetMatchingConditions(conditions)
	}
	return body, nil
}

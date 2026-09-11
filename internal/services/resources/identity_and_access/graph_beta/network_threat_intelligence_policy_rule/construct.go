package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

type threatIntelligenceRuleRequestBody struct {
	*models.ThreatIntelligenceRule
	clearDescription bool
	changed          bool
}

func (b *threatIntelligenceRuleRequestBody) Serialize(w s.SerializationWriter) error {
	if err := b.ThreatIntelligenceRule.Serialize(w); err != nil {
		return fmt.Errorf("serialize threat intelligence payload: %w", err)
	}
	if b.clearDescription {
		if err := w.WriteNullValue("description"); err != nil {
			return fmt.Errorf("serialize null description: %w", err)
		}
	}
	return nil
}
func (b *threatIntelligenceRuleRequestBody) hasChanges() bool { return b.changed }

func constructResource(
	ctx context.Context,
	data *NetworkThreatIntelligencePolicyRuleResourceModel,
) (*threatIntelligenceRuleRequestBody, error) {
	return constructUpdateResource(ctx, data, nil)
}

func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkThreatIntelligencePolicyRuleResourceModel,
) (*threatIntelligenceRuleRequestBody, error) {
	b := &threatIntelligenceRuleRequestBody{
		ThreatIntelligenceRule: models.NewThreatIntelligenceRule(),
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
	if state == nil || !plan.Priority.Equal(state.Priority) {
		v := int64(plan.Priority.ValueInt32())
		b.SetPriority(&v)
		b.changed = true
	}
	if state == nil || !plan.Action.Equal(state.Action) {
		if err := convert.FrameworkToGraphEnum(
			plan.Action,
			models.ParseThreatIntelligenceAction,
			b.SetAction,
		); err != nil {
			return nil, fmt.Errorf("construct action: %w", err)
		}
		b.changed = true
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		v := models.DISABLED_SECURITYRULESTATUS
		if plan.Enabled.ValueBool() {
			v = models.ENABLED_SECURITYRULESTATUS
		}
		settings := models.NewThreatIntelligenceRuleSettings()
		settings.SetStatus(&v)
		b.SetSettings(settings)
		b.changed = true
	}
	if state == nil || !plan.Severity.Equal(state.Severity) ||
		!plan.Destinations.Equal(state.Destinations) {
		conditions := models.NewThreatIntelligenceMatchingConditions()
		if err := convert.FrameworkToGraphEnum(
			plan.Severity,
			models.ParseThreatIntelligenceSeverity,
			conditions.SetSeverity,
		); err != nil {
			return nil, fmt.Errorf("construct severity: %w", err)
		}
		var groups []ThreatIntelligencePolicyRuleDestinationModel
		if diags := plan.Destinations.ElementsAs(ctx, &groups, false); diags.HasError() {
			return nil, fmt.Errorf("%w: decode destinations: %v", errInvalidRequest, diags)
		}
		destinations := make([]models.ThreatIntelligenceDestinationable, 0, len(groups))
		for _, group := range groups {
			if group.Type.ValueString() != destinationTypeFQDN {
				return nil, fmt.Errorf(
					"%w: unsupported destination type %q",
					errInvalidRequest,
					group.Type.ValueString(),
				)
			}
			var values []string
			if diags := group.Values.ElementsAs(ctx, &values, false); diags.HasError() {
				return nil, fmt.Errorf(
					"%w: decode destination values: %v",
					errInvalidRequest,
					diags,
				)
			}
			destination := models.NewThreatIntelligenceFqdnDestination()
			destination.SetValues(values)
			destinations = append(destinations, destination)
		}
		conditions.SetDestinations(destinations)
		b.SetMatchingConditions(conditions)
		b.changed = true
	}
	return b, nil
}

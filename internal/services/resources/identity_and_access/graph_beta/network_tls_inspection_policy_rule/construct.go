package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource uses the generated TLS inspection rule model. action remains
// in additional data because it is currently missing from Graph metadata.
func constructResource(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
) (models.TlsInspectionRuleable, error) {
	body := models.NewTlsInspectionRule()
	body.SetName(data.Name.ValueStringPointer())
	body.SetDescription(data.Description.ValueStringPointer())
	body.GetAdditionalData()["action"] = data.Action.ValueString()
	priority := int64(data.Priority.ValueInt32())
	body.SetPriority(&priority)
	body.SetSettings(ruleSettings(data.Enabled.ValueBool()))
	conditions, err := ruleConditions(ctx, data)
	if err != nil {
		return nil, err
	}
	body.SetMatchingConditions(conditions)
	return body, nil
}

func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkTLSInspectionPolicyRuleResourceModel,
) (models.TlsInspectionRuleable, error) {
	body := models.NewTlsInspectionRule()
	if !plan.Name.Equal(state.Name) {
		body.SetName(plan.Name.ValueStringPointer())
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			body.GetAdditionalData()["description"] = nil
		} else {
			body.SetDescription(plan.Description.ValueStringPointer())
		}
	}
	if !plan.Action.Equal(state.Action) {
		body.GetAdditionalData()["action"] = plan.Action.ValueString()
	}
	if !plan.Priority.Equal(state.Priority) {
		priority := int64(plan.Priority.ValueInt32())
		body.SetPriority(&priority)
	}
	if !plan.Enabled.Equal(state.Enabled) {
		body.SetSettings(ruleSettings(plan.Enabled.ValueBool()))
	}
	if !plan.Destinations.Equal(state.Destinations) {
		conditions, err := ruleConditions(ctx, plan)
		if err != nil {
			return nil, err
		}
		body.SetMatchingConditions(conditions)
	}
	return body, nil
}

func ruleSettings(enabled bool) models.TlsInspectionRuleSettingsable {
	status := models.DISABLED_SECURITYRULESTATUS
	if enabled {
		status = models.ENABLED_SECURITYRULESTATUS
	}
	settings := models.NewTlsInspectionRuleSettings()
	settings.SetStatus(&status)
	return settings
}

func ruleConditions(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
) (models.TlsInspectionMatchingConditionsable, error) {
	var groups []TLSInspectionPolicyRuleDestinationModel
	if diags := data.Destinations.ElementsAs(ctx, &groups, false); diags.HasError() {
		return nil, fmt.Errorf("%w: decode groups: %v", errInvalidDestinations, diags)
	}
	destinations := make([]models.TlsInspectionDestinationable, 0, len(groups))
	for _, group := range groups {
		var values []string
		if diags := group.Values.ElementsAs(ctx, &values, false); diags.HasError() {
			return nil, fmt.Errorf("%w: decode values: %v", errInvalidDestinations, diags)
		}
		switch group.Type.ValueString() {
		case destinationTypeFQDN:
			destination := models.NewTlsInspectionFqdnDestination()
			destination.SetValues(values)
			destinations = append(destinations, destination)
		case destinationTypeWebCategory:
			destination := models.NewTlsInspectionWebCategoryDestination()
			destination.SetValues(values)
			destinations = append(destinations, destination)
		default:
			return nil, fmt.Errorf(
				"%w: unsupported type %q",
				errInvalidDestinations,
				group.Type.ValueString(),
			)
		}
	}
	conditions := models.NewTlsInspectionMatchingConditions()
	conditions.SetDestinations(destinations)
	return conditions, nil
}

func hasUpdateChanges(plan, state *NetworkTLSInspectionPolicyRuleResourceModel) bool {
	return !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) ||
		!plan.Action.Equal(state.Action) || !plan.Priority.Equal(state.Priority) ||
		!plan.Enabled.Equal(state.Enabled) || !plan.Destinations.Equal(state.Destinations)
}

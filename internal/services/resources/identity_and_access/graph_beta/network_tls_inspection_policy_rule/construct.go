package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource builds the TLS inspection rule body. Create passes a nil state and sends every field; Update passes the
// prior state and sends only the fields that changed. action goes through additional data because the SDK model lacks it
// (see the SDK deviations in crud.go).
func constructResource(ctx context.Context, plan, state *NetworkTLSInspectionPolicyRuleResourceModel) (models.TlsInspectionRuleable, error) {
	body := models.NewTlsInspectionRule()
	if state == nil || !plan.Name.Equal(state.Name) {
		body.SetName(plan.Name.ValueStringPointer())
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		if state != nil && plan.Description.IsNull() {
			// The generated setter omits nil strings, so an explicit JSON null clears the description through additional data.
			body.GetAdditionalData()["description"] = nil
		} else {
			body.SetDescription(plan.Description.ValueStringPointer())
		}
	}
	if state == nil || !plan.Action.Equal(state.Action) {
		body.GetAdditionalData()["action"] = plan.Action.ValueString()
	}
	if state == nil || !plan.Priority.Equal(state.Priority) {
		priority := int64(plan.Priority.ValueInt32())
		body.SetPriority(&priority)
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		status := models.DISABLED_SECURITYRULESTATUS
		if plan.Enabled.ValueBool() {
			status = models.ENABLED_SECURITYRULESTATUS
		}
		settings := models.NewTlsInspectionRuleSettings()
		settings.SetStatus(&status)
		body.SetSettings(settings)
	}
	if state == nil || !plan.Destinations.Equal(state.Destinations) {
		var groups []TLSInspectionPolicyRuleDestinationModel
		if diags := plan.Destinations.ElementsAs(ctx, &groups, false); diags.HasError() {
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
				return nil, fmt.Errorf("%w: unsupported type %q", errInvalidDestinations, group.Type.ValueString())
			}
		}
		conditions := models.NewTlsInspectionMatchingConditions()
		conditions.SetDestinations(destinations)
		body.SetMatchingConditions(conditions)
	}
	return body, nil
}

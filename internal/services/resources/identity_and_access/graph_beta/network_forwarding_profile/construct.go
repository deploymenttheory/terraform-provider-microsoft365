package graphBetaNetworkForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// constructResource builds the raw JSON body for the forwarding profile v2 properties that the SDK does not expose.
// With a nil state it builds the POST body; otherwise it builds a PATCH body that carries only the changed values.
func constructResource(ctx context.Context, plan NetworkForwardingProfileResourceModel, state *NetworkForwardingProfileResourceModel) (map[string]any, diag.Diagnostics) {
	body := map[string]any{}
	if state == nil {
		body["isCustomProfile"] = true
		body["trafficForwardingType"] = plan.TrafficForwardingType.ValueString()
	}
	if state == nil || !plan.Name.Equal(state.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsUnknown() && (state == nil || !plan.Description.Equal(state.Description)) {
		if plan.Description.IsNull() {
			if state != nil {
				body["description"] = nil
			}
		} else {
			body["description"] = plan.Description.ValueString()
		}
	}
	if state == nil || !plan.Priority.Equal(state.Priority) {
		body["priority"] = plan.Priority.ValueInt32()
	}
	planned := sharedmodels.ForwardingProfileSettings{
		State:                    plan.State,
		ClientFallbackAction:     plan.ClientFallbackAction,
		ForwardingSourceMatching: plan.ForwardingSourceMatching,
		LocalNetworkSettings:     plan.LocalNetworkSettings,
		AssignmentRules:          plan.AssignmentRules,
	}
	var prior *sharedmodels.ForwardingProfileSettings
	if state != nil {
		prior = &sharedmodels.ForwardingProfileSettings{
			State:                    state.State,
			ClientFallbackAction:     state.ClientFallbackAction,
			ForwardingSourceMatching: state.ForwardingSourceMatching,
			LocalNetworkSettings:     state.LocalNetworkSettings,
			AssignmentRules:          state.AssignmentRules,
		}
	}
	diags := sharedmodels.SettingsPatch(ctx, body, planned, prior)
	return body, diags
}

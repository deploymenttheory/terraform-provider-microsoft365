package graphBetaNetworkDeviceForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// constructResource builds the raw JSON body for the forwarding profile v2 properties that the SDK does not expose.
// With a nil prior state it builds the POST body; otherwise it builds a PATCH body that carries only the changed values.
func constructResource(ctx context.Context, plan, prior *NetworkDeviceForwardingProfileResourceModel) (map[string]any, diag.Diagnostics) {
	body := map[string]any{}
	if prior == nil {
		body["isCustomProfile"] = true
		body["trafficForwardingType"] = plan.TrafficForwardingType.ValueString()
	}
	if prior == nil || !plan.Name.Equal(prior.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsUnknown() && (prior == nil || !plan.Description.Equal(prior.Description)) {
		if plan.Description.IsNull() {
			if prior != nil {
				body["description"] = nil
			}
		} else {
			body["description"] = plan.Description.ValueString()
		}
	}
	if prior == nil || !plan.Priority.Equal(prior.Priority) {
		body["priority"] = plan.Priority.ValueInt32()
	}
	// Device profiles reject agent-only source matching, so it is never sent.
	planned := sharedmodels.ForwardingProfileSettings{
		State:                    plan.State,
		ClientFallbackAction:     plan.ClientFallbackAction,
		ForwardingSourceMatching: types.ObjectNull(sharedmodels.ForwardingSourceMatchingAttrTypes),
		LocalNetworkSettings:     plan.LocalNetworkSettings,
		AssignmentRules:          plan.AssignmentRules,
	}
	var previous *sharedmodels.ForwardingProfileSettings
	if prior != nil {
		previous = &sharedmodels.ForwardingProfileSettings{
			State:                    prior.State,
			ClientFallbackAction:     prior.ClientFallbackAction,
			ForwardingSourceMatching: types.ObjectNull(sharedmodels.ForwardingSourceMatchingAttrTypes),
			LocalNetworkSettings:     prior.LocalNetworkSettings,
			AssignmentRules:          prior.AssignmentRules,
		}
	}
	diags := sharedmodels.SettingsPatch(ctx, body, planned, previous)
	return body, diags
}

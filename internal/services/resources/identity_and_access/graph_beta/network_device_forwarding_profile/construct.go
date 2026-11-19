package graphBetaNetworkDeviceForwardingProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// requestBody builds raw JSON for forwarding profile v2 properties that the SDK does not expose.
func requestBody(
	ctx context.Context,
	plan NetworkDeviceForwardingProfileResourceModel,
	state *NetworkDeviceForwardingProfileResourceModel,
) (map[string]any, diag.Diagnostics) {
	body := map[string]any{}
	if state == nil {
		body["isCustomProfile"] = true
		body["trafficForwardingType"] = plan.TrafficForwardingType.ValueString()
	}
	if state == nil || !plan.Name.Equal(state.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.IsUnknown() &&
		(state == nil || !plan.Description.Equal(state.Description)) {
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
	var prior *sharedmodels.ForwardingProfileSettings
	if state != nil {
		s := state.settings()
		prior = &s
	}
	diags := sharedmodels.SettingsPatch(ctx, body, plan.settings(), prior)
	return body, diags
}

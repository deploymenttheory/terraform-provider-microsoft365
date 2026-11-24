package graphBetaNetworkContentPolicy

import (
	"context"
)

// constructResource builds the content policy body observed from the Entra Global Secure Access content policy blade.
// Create passes a nil state and sends the name, description and settings. Update passes the prior state and sends only
// the fields that changed, matching the portal's PATCH shape while still allowing default_action to be updated. Null
// values are left out of the body.
func constructResource(_ context.Context, plan, state *NetworkContentPolicyResourceModel) map[string]any {
	body := map[string]any{}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := map[string]any{}
		if !plan.DefaultAction.IsNull() {
			settings["defaultAction"] = plan.DefaultAction.ValueString()
		}
		body["settings"] = settings
	}
	if (state == nil || !plan.Name.Equal(state.Name)) && !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		body["name"] = plan.Name.ValueString()
	}
	if (state == nil || !plan.Description.Equal(state.Description)) && !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body["description"] = plan.Description.ValueString()
	}
	return body
}

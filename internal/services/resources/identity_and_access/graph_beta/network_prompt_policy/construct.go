package graphBetaNetworkPromptPolicy

import (
	"context"
)

// constructResource builds the prompt policy body observed from the Entra Global Secure Access prompt policy blade. Create
// passes a nil state and sends the name, description and settings. Update passes the prior state and sends only the
// fields that changed, so an empty body means nothing Graph stores has changed (for example a timeouts-only change).
// A cleared description is sent as JSON null.
func constructResource(_ context.Context, plan, state *NetworkPromptPolicyResourceModel) map[string]any {
	body := map[string]any{}
	if state == nil || !plan.DefaultAction.Equal(state.DefaultAction) {
		settings := map[string]any{}
		if !plan.DefaultAction.IsNull() {
			settings["defaultAction"] = plan.DefaultAction.ValueString()
		}
		body["settings"] = settings
	}
	if state == nil {
		if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
			body["name"] = plan.Name.ValueString()
		}
		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			body["description"] = plan.Description.ValueString()
		}
		return body
	}
	if !plan.Name.Equal(state.Name) && !plan.Name.IsNull() {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Description.Equal(state.Description) {
		body["description"] = plan.Description.ValueStringPointer()
	}
	return body
}

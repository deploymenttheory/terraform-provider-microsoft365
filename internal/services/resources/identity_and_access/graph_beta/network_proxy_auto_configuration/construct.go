package graphBetaNetworkProxyAutoConfiguration

import (
	"errors"
)

var errInvalidPACConfig = errors.New("name, content and is_enabled must be explicitly configured; name and content must be nonempty")

// constructResource builds the raw JSON body for a custom PAC file. With a nil state it builds the POST body;
// otherwise it builds a PATCH body that carries only the changed values, and an empty map when nothing changed.
func constructResource(plan NetworkProxyAutoConfigurationResourceModel, state *NetworkProxyAutoConfigurationResourceModel) (map[string]any, error) {
	if plan.Name.ValueString() == "" || plan.Content.ValueString() == "" || plan.IsEnabled.IsNull() || plan.IsEnabled.IsUnknown() {
		return nil, errInvalidPACConfig
	}
	body := map[string]any{}
	if state == nil || !plan.Name.Equal(state.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if state == nil || !plan.Content.Equal(state.Content) {
		body["content"] = plan.Content.ValueString()
	}
	if state == nil || !plan.IsEnabled.Equal(state.IsEnabled) {
		body["isEnabled"] = plan.IsEnabled.ValueBool()
	}
	return body, nil
}

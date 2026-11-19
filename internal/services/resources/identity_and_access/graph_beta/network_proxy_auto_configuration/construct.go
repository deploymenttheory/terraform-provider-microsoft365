package graphBetaNetworkProxyAutoConfiguration

func pacBody(
	plan NetworkProxyAutoConfigurationResourceModel,
	state *NetworkProxyAutoConfigurationResourceModel,
) (map[string]any, error) {
	if plan.Name.IsNull() || plan.Name.IsUnknown() || plan.Name.ValueString() == "" ||
		plan.Content.IsNull() ||
		plan.Content.IsUnknown() ||
		plan.Content.ValueString() == "" ||
		plan.IsEnabled.IsNull() ||
		plan.IsEnabled.IsUnknown() {
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

package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"
)

// constructResource builds the prompt policy rule body. Create passes a nil state and sends every configured field.
// Update passes the prior state and sends only the fields that changed, so a body holding only @odata.type means nothing
// Graph stores has changed (for example a timeouts-only change). The explicit discriminator is required on PATCH: Graph
// can otherwise silently ignore fields. A cleared description is sent as JSON null, and matchingConditions is always
// sent whole because conversationSchemes and scanResult are replaced together.
func constructResource(ctx context.Context, plan, state *NetworkPromptPolicyRuleResourceModel) (map[string]any, error) {
	body := map[string]any{"@odata.type": "#microsoft.graph.networkaccess.promptRule"}
	// setString omits a nil value, matching the Kiota writer this body replaces.
	setString := func(target map[string]any, key string, value *string) {
		if value != nil {
			target[key] = *value
		}
	}
	if state == nil || !plan.Name.Equal(state.Name) {
		setString(body, "name", plan.Name.ValueStringPointer())
	}
	if state == nil && !plan.Description.IsNull() || state != nil && !plan.Description.Equal(state.Description) {
		body["description"] = plan.Description.ValueStringPointer()
	}
	if state == nil || !plan.Action.Equal(state.Action) {
		setString(body, "action", plan.Action.ValueStringPointer())
	}
	if state == nil || !plan.PromptLogging.Equal(state.PromptLogging) {
		setString(body, "promptLogging", plan.PromptLogging.ValueStringPointer())
	}
	if priority := plan.Priority.ValueInt32Pointer(); priority != nil && (state == nil || !plan.Priority.Equal(state.Priority)) {
		body["priority"] = *priority
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		status := "disabled"
		if plan.Enabled.ValueBool() {
			status = "enabled"
		}
		body["settings"] = map[string]any{"status": status}
	}
	if state != nil && plan.ConversationSchemes.Equal(state.ConversationSchemes) && plan.ScanResult.Equal(state.ScanResult) {
		return body, nil
	}
	var schemes []ConversationSchemeModel
	if diags := plan.ConversationSchemes.ElementsAs(ctx, &schemes, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errInvalidSchemes, diags)
	}
	conversationSchemes := make([]map[string]any, 0, len(schemes))
	for _, scheme := range schemes {
		item := map[string]any{}
		switch scheme.Type.ValueString() {
		case schemeTypeCustom:
			item["@odata.type"] = "#microsoft.graph.networkaccess.customConversationScheme"
			setString(item, "url", scheme.URL.ValueStringPointer())
			item["jsonPath"] = scheme.JSONPath.ValueString()
		case schemeTypePredefined:
			item["@odata.type"] = "#microsoft.graph.networkaccess.predefinedConversationScheme"
			setString(item, "schemeName", scheme.SchemeName.ValueStringPointer())
		default:
			return nil, fmt.Errorf("%w: unsupported type %q", errInvalidSchemes, scheme.Type.ValueString())
		}
		conversationSchemes = append(conversationSchemes, item)
	}
	conditions := map[string]any{"conversationSchemes": conversationSchemes}
	setString(conditions, "scanResult", plan.ScanResult.ValueStringPointer())
	body["matchingConditions"] = conditions
	return body, nil
}

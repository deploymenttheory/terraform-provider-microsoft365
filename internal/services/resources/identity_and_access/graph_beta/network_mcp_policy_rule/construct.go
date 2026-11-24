package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mcpPolicyRuleODataType is the discriminator sent with every write: Graph can otherwise silently ignore fields on PATCH.
const mcpPolicyRuleODataType = "#microsoft.graph.networkaccess.mcpPolicyRule"

// constructResource builds the MCP policy rule body. Create passes a nil state and sends every configured field. Update
// passes the prior state and sends only the fields that changed, so a body holding only @odata.type means nothing Graph
// stores has changed (for example a timeouts-only change). A cleared description or matching_conditions is sent as JSON
// null. Conditions use explicit nulls at each supported level because Graph merges nested PATCH objects.
//
//nolint:gocyclo // The recursive condition encoder is a closure here: the helper rule allows no single-caller helper.
func constructResource(ctx context.Context, plan, state *NetworkMCPPolicyRuleResourceModel) (map[string]any, error) {
	body := map[string]any{"@odata.type": mcpPolicyRuleODataType}
	if (state == nil || !plan.Name.Equal(state.Name)) && !plan.Name.IsNull() {
		body["name"] = plan.Name.ValueString()
	}
	if state == nil && !plan.Description.IsNull() || state != nil && !plan.Description.Equal(state.Description) {
		body["description"] = plan.Description.ValueStringPointer()
	}
	if (state == nil || !plan.Action.Equal(state.Action)) && !plan.Action.IsNull() {
		body["action"] = plan.Action.ValueString()
	}
	if (state == nil || !plan.Priority.Equal(state.Priority)) && !plan.Priority.IsNull() {
		body["priority"] = plan.Priority.ValueInt32()
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		status := "disabled"
		if plan.Enabled.ValueBool() {
			status = "enabled"
		}
		body["settings"] = map[string]any{"status": status}
	}
	if state != nil && plan.MatchingConditions.Equal(state.MatchingConditions) {
		return body, nil
	}
	if plan.MatchingConditions.IsNull() {
		// Create omits absent conditions; Update clears them.
		if state != nil {
			body["matchingConditions"] = nil
		}
		return body, nil
	}
	if plan.MatchingConditions.IsUnknown() {
		return nil, fmt.Errorf("%w: conditions are unknown", errInvalidConditions)
	}

	// conditionFromObject converts a condition object to its Graph JSON, renaming match_type to matchType and sending
	// null attributes as explicit nulls.
	var conditionFromObject func(obj types.Object) (map[string]any, error)
	conditionFromObject = func(obj types.Object) (map[string]any, error) {
		if obj.IsUnknown() {
			return nil, fmt.Errorf("%w: unknown object", errInvalidConditions)
		}
		out := map[string]any{}
		for key, v := range obj.Attributes() {
			apiKey := key
			if key == "match_type" {
				apiKey = "matchType"
			}
			if v.IsUnknown() {
				return nil, fmt.Errorf("%w: %s is unknown", errInvalidConditions, key)
			}
			if v.IsNull() {
				out[apiKey] = nil
				continue
			}
			switch value := v.(type) {
			case types.String:
				out[apiKey] = value.ValueString()
			case types.List:
				values := []string{}
				if d := value.ElementsAs(ctx, &values, false); d.HasError() {
					return nil, fmt.Errorf("%w: %v", errInvalidConditions, d)
				}
				out[apiKey] = values
			case types.Object:
				child, err := conditionFromObject(value)
				if err != nil {
					return nil, err
				}
				out[apiKey] = child
			default:
				return nil, fmt.Errorf("%w: unsupported %s", errInvalidConditions, key)
			}
		}
		return out, nil
	}

	destinations := map[string]any{}
	for tf, api := range mcpConditionKeys {
		v := plan.MatchingConditions.Attributes()[tf]
		if v.IsNull() {
			destinations[api] = nil
			continue
		}
		switch value := v.(type) {
		case types.String:
			if value.IsUnknown() {
				return nil, fmt.Errorf("%w: %s is unknown", errInvalidConditions, tf)
			}
			destinations[api] = value.ValueString()
		case types.Object:
			converted, err := conditionFromObject(value)
			if err != nil {
				return nil, err
			}
			destinations[api] = converted
		default:
			return nil, fmt.Errorf("%w: unsupported condition %s", errInvalidConditions, tf)
		}
	}
	body["matchingConditions"] = map[string]any{"destinations": destinations}
	return body, nil
}

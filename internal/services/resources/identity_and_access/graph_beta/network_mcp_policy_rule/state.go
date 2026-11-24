package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// MapRemoteStateToTerraform validates the response before replacing the last known state. matchingConditions is
// decoded from raw JSON so that any condition the provider does not model is reported instead of silently dropped.
//
//nolint:gocyclo // The recursive condition decoder is a closure here: the helper rule allows no single-caller helper.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkMCPPolicyRuleResourceModel, remote *mcpPolicyRuleResponse) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil || remote.Action == nil ||
		remote.Priority == nil || remote.Settings == nil || remote.Settings.Status == nil ||
		remote.ODataType == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if *remote.ODataType != mcpPolicyRuleODataType {
		return fmt.Errorf("%w: unsupported rule type %q", errInvalidResponse, *remote.ODataType)
	}
	if *remote.Priority < 100 || *remote.Priority > math.MaxInt32 {
		return fmt.Errorf("%w: invalid priority", errInvalidResponse)
	}
	if *remote.Action != "allow" && *remote.Action != "block" {
		return fmt.Errorf("%w: unsupported action %q", errInvalidResponse, *remote.Action)
	}
	status := *remote.Settings.Status
	if status != "enabled" && status != "disabled" {
		return fmt.Errorf("%w: unsupported status %q; enabled cannot be determined", errInvalidResponse, status)
	}
	if len(remote.Conditions) == 0 {
		return fmt.Errorf("%w: missing matchingConditions", errInvalidResponse)
	}

	nonNull := func(raw json.RawMessage) bool {
		return len(raw) > 0 && strings.TrimSpace(string(raw)) != "null"
	}
	// objectFromJSON decodes one condition object. keys maps attribute names to Graph properties at the destinations
	// level; below it attribute names are used as they are, except match_type (matchType). Values outside the schema's
	// allowed sets and non-null properties the schema does not model are errors.
	var objectFromJSON func(raw json.RawMessage, typ types.ObjectType, keys map[string]string) (types.Object, error)
	objectFromJSON = func(raw json.RawMessage, typ types.ObjectType, keys map[string]string) (types.Object, error) {
		var values map[string]json.RawMessage
		if err := json.Unmarshal(raw, &values); err != nil {
			return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: %w", errInvalidResponse, err)
		}
		if _, matching := typ.AttrTypes["match_type"]; matching && (!nonNull(values["values"]) || !nonNull(values["matchType"])) {
			return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: match condition requires values and matchType", errInvalidResponse)
		}
		known := map[string]bool{}
		result := map[string]attr.Value{}
		for tf, t := range typ.AttrTypes {
			api := tf
			if keys != nil {
				api = keys[tf]
			} else if tf == "match_type" {
				api = "matchType"
			}
			known[api] = true
			value := values[api]
			switch at := t.(type) {
			case basetypes.StringType:
				var v *string
				if nonNull(value) {
					if err := json.Unmarshal(value, &v); err != nil {
						return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: %s: %w", errInvalidResponse, api, err)
					}
				}
				allowed := map[string][]string{
					"match_type":          {"exactMatch", "contains", "notExactMatch", "notContains"},
					"insecure_connection": {"required", "excluded"},
					"missing_prm":         {"required", "excluded"},
				}[tf]
				if v != nil && allowed != nil && !slices.Contains(allowed, *v) {
					return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: unsupported %s %q", errInvalidResponse, api, *v)
				}
				result[tf] = types.StringPointerValue(v)
			case types.ListType:
				if !nonNull(value) {
					result[tf] = types.ListNull(at.ElemType)
					continue
				}
				var items []*string
				if err := json.Unmarshal(value, &items); err != nil {
					return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: %s: %w", errInvalidResponse, api, err)
				}
				elems := make([]attr.Value, 0, len(items))
				for _, item := range items {
					if item == nil {
						return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: null list element", errInvalidResponse)
					}
					elems = append(elems, types.StringValue(*item))
				}
				result[tf] = types.ListValueMust(at.ElemType, elems)
			case types.ObjectType:
				if !nonNull(value) {
					result[tf] = types.ObjectNull(at.AttrTypes)
					continue
				}
				child, err := objectFromJSON(value, at, nil)
				if err != nil {
					return child, err
				}
				result[tf] = child
			}
		}
		for key, value := range values {
			if !known[key] && nonNull(value) {
				return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: unsupported non-null MCP condition %s", errInvalidResponse, key)
			}
		}
		obj, diags := types.ObjectValue(typ.AttrTypes, result)
		if diags.HasError() {
			return obj, fmt.Errorf("%w: %v", errInvalidResponse, diags)
		}
		return obj, nil
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(remote.Conditions, &envelope); err != nil {
		return fmt.Errorf("%w: %w", errInvalidResponse, err)
	}
	for key, value := range envelope {
		if key != "destinations" && nonNull(value) {
			return fmt.Errorf("%w: unsupported matchingConditions.%s", errInvalidResponse, key)
		}
	}
	conditions := types.ObjectNull(mcpConditionObjectType.AttrTypes)
	if nonNull(envelope["destinations"]) {
		var err error
		if conditions, err = objectFromJSON(envelope["destinations"], mcpConditionObjectType, mcpConditionKeys); err != nil {
			return err
		}
	}
	data.ID = types.StringValue(*remote.ID)
	data.Name = types.StringValue(*remote.Name)
	data.Description = types.StringPointerValue(remote.Description)
	data.Action = types.StringValue(*remote.Action)
	data.Priority = types.Int32Value(int32(*remote.Priority))
	data.Status = types.StringValue(status)
	data.Enabled = types.BoolValue(status == "enabled")
	data.MatchingConditions = conditions
	return nil
}

package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

var conditionKeys = map[string]string{
	"server_urls":         "serverUrls",
	"protocol_versions":   "protocolVersions",
	"insecure_connection": "insecureConnection",
	"missing_prm":         "missingPrm",
	"tool_matching":       "toolMatching",
	"resource_matching":   "resourceMatching",
	"prompt_matching":     "promptMatching",
}

func matchSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListAttribute{
			MarkdownDescription: "Values in API order. Empty arrays are preserved and are distinct from an omitted condition. Pattern evaluation is performed by the service; regex or wildcard behavior is not guaranteed by this resource.",
			Required:            true,
			ElementType:         types.StringType,
		},
		"match_type": schema.StringAttribute{
			MarkdownDescription: "Match operator: exactMatch, contains, notExactMatch, or notContains. Must be explicitly configured.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("exactMatch", "contains", "notExactMatch", "notContains"),
			},
		},
	}
}

func primitiveSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"names": schema.SingleNestedAttribute{
			MarkdownDescription: "Primitive name matching.",
			Optional:            true,
			Attributes:          matchSchema(),
		},
		"methods": schema.StringAttribute{
			MarkdownDescription: "API method flags, for example call or list,call for tools, read for resources, or get for prompts. Values are sent unchanged; unsupported combinations return API diagnostics.",
			Optional:            true,
			Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
		},
	}
}

func conditionSchema() map[string]schema.Attribute {
	result := map[string]schema.Attribute{}
	for name := range conditionKeys {
		switch name {
		case "insecure_connection", "missing_prm":
			result[name] = schema.StringAttribute{
				MarkdownDescription: "Whether this property is required or excluded by the condition. Omission removes this constraint.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("required", "excluded"),
				},
			}
		case "tool_matching", "resource_matching", "prompt_matching":
			result[name] = schema.SingleNestedAttribute{
				MarkdownDescription: "Primitive conditions applied to the rule. Multiple primitive types may be stored by the API; their combined traffic evaluation is not verified.",
				Optional:            true,
				Attributes:          primitiveSchema(),
			}
		default:
			result[name] = schema.SingleNestedAttribute{
				MarkdownDescription: "String matching condition. Values are preserved without URL or case normalization.",
				Optional:            true,
				Attributes:          matchSchema(),
			}
		}
	}
	return result
}

func conditionObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: conditionAttributeTypes()}
}

func conditionAttributeTypes() map[string]attr.Type {
	out := map[string]attr.Type{}
	for k, v := range conditionSchema() {
		out[k] = v.GetType()
	}
	return out
}

// Conditions use explicit nulls at each supported level because Graph merges nested PATCH objects.
func ruleConditions(
	ctx context.Context,
	data *NetworkMCPPolicyRuleResourceModel,
) (*conditionPayload, error) {
	if data.MatchingConditions.IsNull() {
		return nil, nil //nolint:nilnil // A nil payload represents an intentionally omitted condition.
	}
	if data.MatchingConditions.IsUnknown() {
		return nil, fmt.Errorf("%w: conditions are unknown", errInvalidConditions)
	}
	dest := conditionPayload{}
	for tf, api := range conditionKeys {
		v := data.MatchingConditions.Attributes()[tf]
		if v.IsNull() {
			dest[api] = nil
			continue
		}
		switch value := v.(type) {
		case types.String:
			if value.IsUnknown() {
				return nil, fmt.Errorf("%w: %s is unknown", errInvalidConditions, tf)
			}
			dest[api] = value.ValueString()
		case types.Object:
			converted, err := conditionFromObject(ctx, value)
			if err != nil {
				return nil, err
			}
			dest[api] = converted
		default:
			return nil, fmt.Errorf("%w: unsupported condition %s", errInvalidConditions, tf)
		}
	}
	body := conditionPayload{"destinations": dest}
	return &body, nil
}

func conditionFromObject(ctx context.Context, obj types.Object) (conditionPayload, error) {
	if obj.IsUnknown() {
		return nil, fmt.Errorf("%w: unknown object", errInvalidConditions)
	}
	out := conditionPayload{}
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
			child, err := conditionFromObject(ctx, value)
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

type conditionPayload map[string]any

func (b *conditionPayload) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *conditionPayload) Serialize(w s.SerializationWriter) error {
	for key, value := range *b {
		var err error
		switch v := value.(type) {
		case nil:
			err = w.WriteNullValue(key)
		case string:
			err = w.WriteStringValue(key, &v)
		case []string:
			err = w.WriteCollectionOfStringValues(key, v)
		case conditionPayload:
			err = w.WriteObjectValue(key, &v)
		default:
			return fmt.Errorf("%w: unsupported value %T", errInvalidConditions, value)
		}
		if err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

func conditionsToState(ctx context.Context, raw json.RawMessage) (types.Object, error) {
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return types.ObjectNull(
				conditionAttributeTypes(),
			), fmt.Errorf(
				"%w: %w",
				errInvalidResponse,
				err,
			)
	}
	for key, value := range envelope {
		if key != "destinations" && nonNull(value) {
			return types.ObjectNull(
					conditionAttributeTypes(),
				), fmt.Errorf(
					"%w: unsupported matchingConditions.%s",
					errInvalidResponse,
					key,
				)
		}
	}
	if !nonNull(envelope["destinations"]) {
		return types.ObjectNull(conditionAttributeTypes()), nil
	}
	return objectFromJSON(ctx, envelope["destinations"], conditionObjectType(), conditionKeys)
}

func nonNull(raw json.RawMessage) bool {
	return len(raw) > 0 && strings.TrimSpace(string(raw)) != "null"
}

func objectFromJSON(
	ctx context.Context,
	raw json.RawMessage,
	typ types.ObjectType,
	keys map[string]string,
) (types.Object, error) {
	var values map[string]json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return types.ObjectNull(typ.AttrTypes), fmt.Errorf("%w: %w", errInvalidResponse, err)
	}
	if _, matching := typ.AttrTypes["match_type"]; matching {
		if !nonNull(values["values"]) || !nonNull(values["matchType"]) {
			return types.ObjectNull(
					typ.AttrTypes,
				), fmt.Errorf(
					"%w: match condition requires values and matchType",
					errInvalidResponse,
				)
		}
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
					return types.ObjectNull(
							typ.AttrTypes,
						), fmt.Errorf(
							"%w: %s: %w",
							errInvalidResponse,
							api,
							err,
						)
				}
			}
			if v != nil {
				allowed := []string(nil)
				switch tf {
				case "match_type":
					allowed = []string{"exactMatch", "contains", "notExactMatch", "notContains"}
				case "insecure_connection", "missing_prm":
					allowed = []string{"required", "excluded"}
				}
				if allowed != nil && !slices.Contains(allowed, *v) {
					return types.ObjectNull(
							typ.AttrTypes,
						), fmt.Errorf(
							"%w: unsupported %s %q",
							errInvalidResponse,
							api,
							*v,
						)
				}
			}
			result[tf] = types.StringPointerValue(v)
		case types.ListType:
			if !nonNull(value) {
				result[tf] = types.ListNull(at.ElemType)
				continue
			}
			var items []*string
			if err := json.Unmarshal(value, &items); err != nil {
				return types.ObjectNull(
						typ.AttrTypes,
					), fmt.Errorf(
						"%w: %s: %w",
						errInvalidResponse,
						api,
						err,
					)
			}
			elems := make([]attr.Value, 0, len(items))
			for _, item := range items {
				if item == nil {
					return types.ObjectNull(
							typ.AttrTypes,
						), fmt.Errorf(
							"%w: null list element",
							errInvalidResponse,
						)
				}
				elems = append(elems, types.StringValue(*item))
			}
			result[tf] = types.ListValueMust(at.ElemType, elems)
		case types.ObjectType:
			if !nonNull(value) {
				result[tf] = types.ObjectNull(at.AttrTypes)
				continue
			}
			child, err := objectFromJSON(ctx, value, at, nil)
			if err != nil {
				return child, err
			}
			result[tf] = child
		}
	}
	for key, value := range values {
		if !known[key] && nonNull(value) {
			return types.ObjectNull(
					typ.AttrTypes,
				), fmt.Errorf(
					"%w: unsupported non-null MCP condition %s",
					errInvalidResponse,
					key,
				)
		}
	}
	obj, diags := types.ObjectValue(typ.AttrTypes, result)
	if diags.HasError() {
		return obj, fmt.Errorf("%w: %v", errInvalidResponse, diags)
	}
	return obj, nil
}

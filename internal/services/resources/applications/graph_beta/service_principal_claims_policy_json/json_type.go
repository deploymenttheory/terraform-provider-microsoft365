package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	errPolicyObject   = errors.New("policy_json must contain a JSON object")
	errPolicyShape    = errors.New("policy_json must explicitly include Boolean includeBasicClaimSet and includeApplicationIdInIssuer, and a claims array")
	errPolicyKnown    = errors.New("policy_json must be known before applying")
	errPolicyTrailing = errors.New("unexpected data after the policy JSON object")
	errPolicyResponse = errors.New("empty or invalid claims policy response from Microsoft Graph")
)

type (
	PolicyJSONType  struct{ basetypes.StringType }
	PolicyJSONValue struct{ basetypes.StringValue }
)

var (
	_ basetypes.StringTypable                    = PolicyJSONType{}
	_ basetypes.StringValuableWithSemanticEquals = PolicyJSONValue{}
)

func (t PolicyJSONType) Equal(other attr.Type) bool {
	value, ok := other.(PolicyJSONType)
	return ok && t.StringType.Equal(value.StringType)
}
func (t PolicyJSONType) String() string                       { return "PolicyJSONType" }
func (t PolicyJSONType) ValueType(context.Context) attr.Value { return PolicyJSONValue{} }
func (t PolicyJSONType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return PolicyJSONValue{StringValue: in}, nil
}

func (t PolicyJSONType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	value, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("convert claims policy JSON: %w", err)
	}
	text, ok := value.(basetypes.StringValue)
	if !ok {
		return nil, errPolicyObject
	}
	return PolicyJSONValue{StringValue: text}, nil
}

func (v PolicyJSONValue) Equal(other attr.Value) bool {
	value, ok := other.(PolicyJSONValue)
	return ok && v.StringValue.Equal(value.StringValue)
}
func (v PolicyJSONValue) Type(context.Context) attr.Type { return PolicyJSONType{} }
func NewPolicyJSONValue(value string) PolicyJSONValue {
	return PolicyJSONValue{StringValue: basetypes.NewStringValue(value)}
}

func (v PolicyJSONValue) StringSemanticEquals(_ context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	value, ok := other.(PolicyJSONValue)
	if !ok || v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return false, nil
	}
	left, err := normalisePolicyJSON([]byte(v.ValueString()))
	if err != nil {
		return false, nil
	}
	right, err := normalisePolicyJSON([]byte(value.ValueString()))
	return err == nil && bytes.Equal(left, right), nil
}

type policyJSONValidator struct{}

func (v policyJSONValidator) Description(context.Context) string {
	return "A complete custom claims policy with explicit root flags and claims array"
}

func (v policyJSONValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v policyJSONValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := decodePolicyJSON([]byte(req.ConfigValue.ValueString())); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid claims policy JSON", err.Error())
	}
}

func decodePolicyJSON(content []byte) (map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode claims policy JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errPolicyTrailing
	}
	object, ok := value.(map[string]any)
	if !ok {
		return nil, errPolicyObject
	}
	_, basic := object["includeBasicClaimSet"].(bool)
	_, issuer := object["includeApplicationIdInIssuer"].(bool)
	_, claims := object["claims"].([]any)
	if !basic || !issuer || !claims {
		return nil, errPolicyShape
	}
	return object, nil
}

func stripServerPolicyFields(object map[string]any) {
	for _, key := range []string{"id", "@odata.context", "@odata.etag", "@odata.type"} {
		delete(object, key)
	}
}

func removeNull(object map[string]any, key string) {
	if object[key] == nil {
		delete(object, key)
	}
}

func normalisePolicyJSON(content []byte) ([]byte, error) {
	object, err := decodePolicyJSON(content)
	if err != nil {
		return nil, err
	}
	stripServerPolicyFields(object)
	removeNull(object, "audienceOverride")
	removeNull(object, "groupFilter")
	claims, _ := object["claims"].([]any)
	for _, entry := range claims {
		claim, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		kind, _ := claim["@odata.type"].(string)
		if kind != "#microsoft.graph.customClaim" && kind != "#microsoft.graph.samlNameIdClaim" {
			continue
		}
		if claim["@odata.type"] == "#microsoft.graph.samlNameIdClaim" {
			removeNull(claim, "serviceProviderNameQualifier")
		}
		configurations, _ := claim["configurations"].([]any)
		for _, entry := range configurations {
			configuration, ok := entry.(map[string]any)
			if !ok {
				continue
			}
			removeNull(configuration, "condition")
			removeNull(configuration, "attribute")
			if configuration["transformations"] == nil {
				configuration["transformations"] = []any{}
			}
			transformations, _ := configuration["transformations"].([]any)
			for _, entry := range transformations {
				transformation, ok := entry.(map[string]any)
				if !ok {
					continue
				}
				kind, _ := transformation["@odata.type"].(string)
				if strings.HasSuffix(kind, "Transformation") {
					removeNull(transformation, "input")
				}
			}
		}
	}
	result, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("encode normalized claims policy: %w", err)
	}
	return result, nil
}

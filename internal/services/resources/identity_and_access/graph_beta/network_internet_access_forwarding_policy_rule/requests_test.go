package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	store "github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewInternetAccessForwardingRuleRequestInformationSerializesObservedFQDNPayload(t *testing.T) {
	payload := serializedRulePayload(t, context.Background(), &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		Name:     types.StringValue("Custom Acquire policy internet rule"),
		Action:   types.StringValue("forward"),
		RuleType: types.StringValue(ruleTypeFQDN),
		Ports:    stringSetForTest(t, []string{"80", "443"}),
		Protocol: types.StringValue("udp"),
		Destinations: destinationsForTest(t, []map[string]string{
			{"type": ruleTypeFQDN, "value": "example.com"},
		}),
	}, false)

	expected := map[string]any{
		"@odata.type": internetAccessForwardingRuleODataType,
		"name":        "Custom Acquire policy internet rule",
		"action":      "forward",
		"ruleType":    "fqdn",
		"ports":       []any{"80", "443"},
		"protocol":    "udp",
		"destinations": []any{
			map[string]any{
				"@odata.type": fqdnODataType,
				"value":       "example.com",
			},
		},
	}

	assertRuleJSONMapEqual(t, expected, payload)
}

func TestInternetAccessForwardingRuleSerializesSupportedDestinationTypes(t *testing.T) {
	tests := []struct {
		name        string
		ruleType    string
		destination map[string]string
		expected    map[string]any
	}{
		{
			name:        "ip address",
			ruleType:    ruleTypeIPAddress,
			destination: map[string]string{"type": ruleTypeIPAddress, "value": "192.0.2.10"},
			expected:    map[string]any{"@odata.type": ipAddressODataType, "value": "192.0.2.10"},
		},
		{
			name:        "ip range",
			ruleType:    ruleTypeIPRange,
			destination: map[string]string{"type": ruleTypeIPRange, "begin_address": "192.0.2.10", "end_address": "192.0.2.20"},
			expected:    map[string]any{"@odata.type": ipRangeODataType, "beginAddress": "192.0.2.10", "endAddress": "192.0.2.20"},
		},
		{
			name:        "cidr subnet",
			ruleType:    ruleTypeIPSubnet,
			destination: map[string]string{"type": ruleTypeIPSubnet, "value": "192.0.2.0/24"},
			expected:    map[string]any{"@odata.type": ipSubnetODataType, "value": "192.0.2.0/24"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := serializedRulePayload(t, context.Background(), &NetworkInternetAccessForwardingPolicyRuleResourceModel{
				Name:     types.StringValue("sample"),
				Action:   types.StringValue("bypass"),
				RuleType: types.StringValue(tt.ruleType),
				Ports:    stringSetForTest(t, []string{"443"}),
				Protocol: types.StringValue("tcp"),
				Destinations: destinationsForTest(t, []map[string]string{
					tt.destination,
				}),
			}, false)

			destinations := payload["destinations"].([]any)
			assertRuleJSONMapEqual(t, tt.expected, destinations[0].(map[string]any))
			if payload["ruleType"] == tt.ruleType {
				t.Fatalf("ruleType = %q, expected Graph casing", payload["ruleType"])
			}
		})
	}
}

func TestInternetAccessForwardingRuleUpdateIncludesID(t *testing.T) {
	payload := serializedRulePayload(t, context.Background(), &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		ID:       types.StringValue("66f1c2ee-9f17-4681-a7bc-c324e7dff554"),
		Name:     types.StringValue("Custom Acquire policy internet rule"),
		Action:   types.StringValue("forward"),
		RuleType: types.StringValue(ruleTypeFQDN),
		Ports:    stringSetForTest(t, []string{"80", "443"}),
		Protocol: types.StringValue("tcp"),
		Destinations: destinationsForTest(t, []map[string]string{
			{"type": ruleTypeFQDN, "value": "example.com"},
		}),
	}, true)

	if payload["id"] != "66f1c2ee-9f17-4681-a7bc-c324e7dff554" {
		t.Fatalf("id = %#v, expected update payload to include rule id", payload["id"])
	}
	if payload["protocol"] != "tcp" {
		t.Fatalf("protocol = %#v, expected tcp", payload["protocol"])
	}
	if _, ok := payload["name"]; ok {
		t.Fatalf("update payload included name, but Graph rejects name updates")
	}
	if _, ok := payload["action"]; ok {
		t.Fatalf("update payload included action, but action changes require replacement")
	}
}

func TestConstructResourceValidatesRequiredDestinationShape(t *testing.T) {
	_, err := constructResource(context.Background(), &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		Name:     types.StringValue("bad-range"),
		Action:   types.StringValue("forward"),
		RuleType: types.StringValue(ruleTypeIPRange),
		Ports:    stringSetForTest(t, []string{"443"}),
		Protocol: types.StringValue("tcp"),
		Destinations: destinationsForTest(t, []map[string]string{
			{"type": ruleTypeIPRange, "begin_address": "192.0.2.10"},
		}),
	}, false)
	if err == nil {
		t.Fatal("constructResource returned nil error, expected destination validation error")
	}
	if !strings.Contains(err.Error(), "begin_address and end_address") {
		t.Fatalf("constructResource error = %q, expected ip range validation error", err.Error())
	}
}

func TestConstructResourceValidatesRuleTypeMatchesDestinations(t *testing.T) {
	_, err := constructResource(context.Background(), &NetworkInternetAccessForwardingPolicyRuleResourceModel{
		Name:     types.StringValue("mismatched"),
		Action:   types.StringValue("forward"),
		RuleType: types.StringValue(ruleTypeFQDN),
		Ports:    stringSetForTest(t, []string{"443"}),
		Protocol: types.StringValue("tcp"),
		Destinations: destinationsForTest(t, []map[string]string{
			{"type": ruleTypeIPAddress, "value": "192.0.2.10"},
		}),
	}, false)
	if err == nil {
		t.Fatal("constructResource returned nil error, expected rule_type mismatch validation error")
	}
	if !strings.Contains(err.Error(), "must match rule_type") {
		t.Fatalf("constructResource error = %q, expected rule_type mismatch validation error", err.Error())
	}
}

func serializedRulePayload(t *testing.T, ctx context.Context, model *NetworkInternetAccessForwardingPolicyRuleResourceModel, includeID bool) map[string]any {
	t.Helper()

	body, err := constructResource(ctx, model, includeID)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	requestInfo, err := newInternetAccessForwardingRuleRequestInformation(
		ctx,
		internetAccessForwardingPolicyRuleTestRequestAdapter{},
		abstractions.POST,
		"dad2a411-e330-440d-a7c7-2c830dce5991",
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newInternetAccessForwardingRuleRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	return payload
}

func destinationsForTest(t *testing.T, destinations []map[string]string) types.List {
	t.Helper()

	values := make([]attr.Value, 0, len(destinations))
	for _, destination := range destinations {
		values = append(values, types.ObjectValueMust(destinationObjectType().AttrTypes, map[string]attr.Value{
			"type":          types.StringValue(destination["type"]),
			"value":         optionalStringValue(destination["value"]),
			"begin_address": optionalStringValue(destination["begin_address"]),
			"end_address":   optionalStringValue(destination["end_address"]),
		}))
	}
	return types.ListValueMust(destinationObjectType(), values)
}

func optionalStringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func stringSetForTest(t *testing.T, values []string) types.Set {
	t.Helper()

	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %s", diags.Errors()[0].Detail())
	}
	return set
}

func assertRuleJSONMapEqual(t *testing.T, expected, actual map[string]any) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal expected: %v", err)
	}
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("marshal actual: %v", err)
	}
	if string(expectedJSON) != string(actualJSON) {
		t.Fatalf("unexpected JSON\nexpected: %s\nactual:   %s", expectedJSON, actualJSON)
	}
}

type internetAccessForwardingPolicyRuleTestRequestAdapter struct{}

func (internetAccessForwardingPolicyRuleTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) SetBaseUrl(string) {}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (internetAccessForwardingPolicyRuleTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

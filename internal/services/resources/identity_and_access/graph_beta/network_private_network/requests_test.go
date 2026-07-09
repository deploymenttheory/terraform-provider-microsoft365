package graphBetaNetworkPrivateNetwork

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestPrivateNetworkRequestInformationSerializesPortalPayload(t *testing.T) {
	body := serializedPayload(t, context.Background(), &NetworkPrivateNetworkResourceModel{
		Name:   testString("sample-private-network"),
		AppIDs: testStringSet(t, []string{"81ff8e33-3181-475b-ade6-e27147234bb3"}),
		DNSResolutionIdentification: &DNSResolutionIdentificationModel{
			DNSServers:    testStringSet(t, []string{"8.8.8.8"}),
			FQDNToResolve: testString("example.com"),
			ExpectedIPResolutions: testExpectedIPResolutionSet(t, []ExpectedIPResolutionModel{
				{Type: testString(expectedIPResolutionTypeIPAddress), Value: testString("192.168.1.11")},
				{Type: testString(expectedIPResolutionTypeIPSubnet), Value: testString("192.168.1.1/16")},
				{Type: testString(expectedIPResolutionTypeIPRange), BeginAddress: testString("192.168.1.1"), EndAddress: testString("192.168.1.2")},
			}),
		},
	}, true)

	if body["@odata.type"] != privateNetworkODataType {
		t.Fatalf("@odata.type = %#v, expected %#v", body["@odata.type"], privateNetworkODataType)
	}
	if body["name"] != "sample-private-network" {
		t.Fatalf("name = %#v, expected sample-private-network", body["name"])
	}

	appIDs := body["appIds"].([]any)
	if len(appIDs) != 1 || appIDs[0] != "81ff8e33-3181-475b-ade6-e27147234bb3" {
		t.Fatalf("appIds = %#v", appIDs)
	}

	networkIdentifications := body["networkIdentifications"].([]any)
	dnsIdentification := networkIdentifications[0].(map[string]any)
	expected := dnsIdentification["expectedIpResolutions"].([]any)
	if len(expected) != 3 {
		t.Fatalf("expectedIpResolutions length = %d, expected 3", len(expected))
	}

	expectedODataTypes := map[string]bool{}
	for _, item := range expected {
		resolution := item.(map[string]any)
		odataType, ok := resolution["@odata.type"].(string)
		if !ok {
			t.Fatalf("expectedIpResolutions item missing @odata.type: %#v", resolution)
		}
		expectedODataTypes[odataType] = true
	}
	for _, odataType := range []string{ipAddressODataType, ipSubnetODataType, ipRangeODataType} {
		if !expectedODataTypes[odataType] {
			t.Fatalf("expectedIpResolutions missing @odata.type %q: %#v", odataType, expectedODataTypes)
		}
	}
}

func serializedPayload(t *testing.T, ctx context.Context, model *NetworkPrivateNetworkResourceModel, includeODataType bool) map[string]any {
	t.Helper()
	body, err := constructResource(ctx, model, includeODataType)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	requestInfo, err := newPrivateNetworkRequestInformation(ctx, privateNetworkTestRequestAdapter{}, abstractions.POST, "", body)
	if err != nil {
		t.Fatalf("newPrivateNetworkRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	return payload
}

type privateNetworkTestRequestAdapter struct{}

func (privateNetworkTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (privateNetworkTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (privateNetworkTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (privateNetworkTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (privateNetworkTestRequestAdapter) SetBaseUrl(string) {}

func (privateNetworkTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}

func (privateNetworkTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

func testString(value string) types.String {
	return types.StringValue(value)
}

func testStringSet(t *testing.T, values []string) types.Set {
	t.Helper()
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		t.Fatalf("failed to create string set: %s", diags.Errors()[0].Detail())
	}
	return set
}

func testExpectedIPResolutionSet(t *testing.T, values []ExpectedIPResolutionModel) types.Set {
	t.Helper()
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		attrs := map[string]attr.Value{
			"type":          value.Type,
			"value":         types.StringNull(),
			"begin_address": types.StringNull(),
			"end_address":   types.StringNull(),
		}
		if !value.Value.IsNull() {
			attrs["value"] = value.Value
		}
		if !value.BeginAddress.IsNull() {
			attrs["begin_address"] = value.BeginAddress
		}
		if !value.EndAddress.IsNull() {
			attrs["end_address"] = value.EndAddress
		}

		objectValue, diags := types.ObjectValue(expectedIPResolutionAttrTypes(), attrs)
		if diags.HasError() {
			t.Fatalf("failed to create expected IP resolution object: %s", diags.Errors()[0].Detail())
		}
		elements = append(elements, objectValue)
	}

	set, diags := types.SetValue(expectedIPResolutionObjectType(), elements)
	if diags.HasError() {
		t.Fatalf("failed to create expected IP resolution set: %s", diags.Errors()[0].Detail())
	}
	return set
}

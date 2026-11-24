package graphBetaNetworkPrivateNetwork

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructResourceSerializesPortalPayload(t *testing.T) {
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

	content, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	return payload
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

		objectValue, diags := types.ObjectValue(expectedIPResolutionAttrTypes, attrs)
		if diags.HasError() {
			t.Fatalf("failed to create expected IP resolution object: %s", diags.Errors()[0].Detail())
		}
		elements = append(elements, objectValue)
	}

	set, diags := types.SetValue(types.ObjectType{AttrTypes: expectedIPResolutionAttrTypes}, elements)
	if diags.HasError() {
		t.Fatalf("failed to create expected IP resolution set: %s", diags.Errors()[0].Detail())
	}
	return set
}

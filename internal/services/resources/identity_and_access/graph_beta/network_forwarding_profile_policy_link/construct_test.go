package graphBetaNetworkForwardingProfilePolicyLink

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestConstructResourceSerializesObservedPatchPayload(t *testing.T) {
	body, err := constructResource(&NetworkForwardingProfilePolicyLinkResourceModel{
		State: types.StringValue("enabled"),
	})
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("WriteObjectValue returned error: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("GetSerializedContent returned error: %v", err)
	}
	var actual map[string]any
	if err := json.Unmarshal(content, &actual); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	expected := map[string]any{
		"@odata.type": forwardingPolicyLinkODataType,
		"state":       "enabled",
	}

	assertPolicyLinkJSONMapEqual(t, expected, actual)
}

func TestConstructResourceRejectsInvalidState(t *testing.T) {
	if _, err := constructResource(&NetworkForwardingProfilePolicyLinkResourceModel{State: types.StringValue("invalid")}); err == nil {
		t.Fatal("constructResource returned nil error for invalid state")
	}
}

func assertPolicyLinkJSONMapEqual(t *testing.T, expected, actual map[string]any) {
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

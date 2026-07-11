package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	store "github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestConstructResourceSerializesObservedPatchPayload(t *testing.T) {
	body := constructResource(&NetworkForwardingProfilePolicyLinkResourceModel{
		State: types.StringValue("enabled"),
	})

	actual := serializePolicyLinkParsableForTest(t, body)
	expected := map[string]any{
		"@odata.type": forwardingPolicyLinkODataType,
		"state":       "enabled",
	}

	assertPolicyLinkJSONMapEqual(t, expected, actual)
}

func TestNewForwardingPolicyLinkRequestInformation(t *testing.T) {
	requestInfo, err := newForwardingPolicyLinkRequestInformation(
		context.Background(),
		forwardingPolicyLinkTestRequestAdapter{},
		abstractions.PATCH,
		"72661c0d-027e-4dff-8c76-af103f200903",
		"09837256-2cba-4dde-a121-4d6a129f13db",
		&forwardingPolicyLinkStateRequestBody{ODataType: forwardingPolicyLinkODataType, State: "disabled"},
	)
	if err != nil {
		t.Fatalf("newForwardingPolicyLinkRequestInformation returned error: %v", err)
	}

	if requestInfo.Method != abstractions.PATCH {
		t.Fatalf("method = %s, want PATCH", requestInfo.Method)
	}
	if requestInfo.UrlTemplate != forwardingProfilePolicyLinkURLTemplate {
		t.Fatalf("UrlTemplate = %q, want %q", requestInfo.UrlTemplate, forwardingProfilePolicyLinkURLTemplate)
	}
	if requestInfo.PathParameters["forwardingProfileId"] != "72661c0d-027e-4dff-8c76-af103f200903" {
		t.Fatalf("forwardingProfileId path parameter not set")
	}
	if requestInfo.PathParameters["policyLinkId"] != "09837256-2cba-4dde-a121-4d6a129f13db" {
		t.Fatalf("policyLinkId path parameter not set")
	}
}

func serializePolicyLinkParsableForTest(t *testing.T, body s.Parsable) map[string]any {
	t.Helper()

	requestInfo := abstractions.NewRequestInformation()
	if err := requestInfo.SetContentFromParsable(context.Background(), forwardingPolicyLinkTestRequestAdapter{}, "application/json", body); err != nil {
		t.Fatalf("SetContentFromParsable returned error: %v", err)
	}

	var actual map[string]any
	if err := json.Unmarshal(requestInfo.Content, &actual); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}
	return actual
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

type forwardingPolicyLinkTestRequestAdapter struct{}

func (forwardingPolicyLinkTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (forwardingPolicyLinkTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (forwardingPolicyLinkTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (forwardingPolicyLinkTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}
func (forwardingPolicyLinkTestRequestAdapter) SetBaseUrl(string)                            {}
func (forwardingPolicyLinkTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (forwardingPolicyLinkTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

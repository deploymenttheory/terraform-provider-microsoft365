package graphBetaNetworkWebFilteringPolicy

import (
	"context"
	"encoding/json"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewWebFilteringPolicyRequestInformationSerializesPortalPayload(t *testing.T) {
	name := "sample-for-codex"
	description := "sample for codex"
	body := &webFilteringPolicyRequestBody{
		name:        &name,
		description: &description,
		settings: &webFilteringPolicySettingsRequestBody{
			defaultAction: &webFilteringPolicyDefaultActionRequestBody{
				odataType: defaultActionAllowODataType,
			},
		},
		includePolicyRules: true,
	}

	requestInfo, err := newWebFilteringPolicyRequestInformation(
		context.Background(),
		webFilteringPolicyTestRequestAdapter{},
		abstractions.POST,
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newWebFilteringPolicyRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	if payload["name"] != "sample-for-codex" {
		t.Fatalf("name = %#v, expected sample-for-codex", payload["name"])
	}
	if payload["description"] != "sample for codex" {
		t.Fatalf("description = %#v, expected sample for codex", payload["description"])
	}
	settings := payload["settings"].(map[string]any)
	defaultAction := settings["defaultAction"].(map[string]any)
	if defaultAction["@odata.type"] != defaultActionAllowODataType {
		t.Fatalf("defaultAction @odata.type = %#v, expected %#v", defaultAction["@odata.type"], defaultActionAllowODataType)
	}
	policyRules := payload["policyRules"].([]any)
	if len(policyRules) != 0 {
		t.Fatalf("policyRules length = %d, expected 0", len(policyRules))
	}
}

type webFilteringPolicyTestRequestAdapter struct{}

func (webFilteringPolicyTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webFilteringPolicyTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (webFilteringPolicyTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (webFilteringPolicyTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (webFilteringPolicyTestRequestAdapter) SetBaseUrl(string) {}

func (webFilteringPolicyTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}

func (webFilteringPolicyTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

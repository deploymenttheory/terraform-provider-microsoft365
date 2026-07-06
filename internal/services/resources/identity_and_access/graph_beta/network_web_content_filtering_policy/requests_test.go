package graphBetaNetworkWebContentFilteringPolicy

import (
	"context"
	"encoding/json"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewWebContentFilteringPolicyRequestInformationSerializesPortalPayload(t *testing.T) {
	name := "sample-for-codex"
	description := "sample for codex"
	body := &webContentFilteringPolicyRequestBody{
		name:        &name,
		description: &description,
		settings: &webContentFilteringPolicySettingsRequestBody{
			defaultAction: &webContentFilteringPolicyDefaultActionRequestBody{
				odataType: defaultActionAllowODataType,
			},
		},
		includePolicyRules: true,
	}

	requestInfo, err := newWebContentFilteringPolicyRequestInformation(
		context.Background(),
		webContentFilteringPolicyTestRequestAdapter{},
		abstractions.POST,
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newWebContentFilteringPolicyRequestInformation returned error: %v", err)
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

type webContentFilteringPolicyTestRequestAdapter struct{}

func (webContentFilteringPolicyTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (webContentFilteringPolicyTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (webContentFilteringPolicyTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (webContentFilteringPolicyTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (webContentFilteringPolicyTestRequestAdapter) SetBaseUrl(string) {}

func (webContentFilteringPolicyTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}

func (webContentFilteringPolicyTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

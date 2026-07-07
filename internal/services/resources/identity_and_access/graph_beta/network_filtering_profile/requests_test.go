package graphBetaNetworkFilteringProfile

import (
	"context"
	"encoding/json"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewFilteringProfileRequestInformationSerializesObservedGraphPayload(t *testing.T) {
	name := "test-profile"
	description := "Managed by Terraform"
	priority := int64(100)
	state := "enabled"
	body := &filteringProfileRequestBody{
		name:        &name,
		description: &description,
		priority:    &priority,
		state:       &state,
	}

	requestInfo, err := newFilteringProfileRequestInformation(
		context.Background(),
		filteringProfileTestRequestAdapter{},
		abstractions.POST,
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newFilteringProfileRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	expected := map[string]any{
		"name":        "test-profile",
		"description": "Managed by Terraform",
		"priority":    float64(100),
		"state":       "enabled",
	}
	if len(payload) != len(expected) {
		t.Fatalf("payload keys = %#v, expected only %#v", payload, expected)
	}
	for key, value := range expected {
		if payload[key] != value {
			t.Fatalf("%s = %#v, expected %#v", key, payload[key], value)
		}
	}
}

type filteringProfileTestRequestAdapter struct{}

func (filteringProfileTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (filteringProfileTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (filteringProfileTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (filteringProfileTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (filteringProfileTestRequestAdapter) SetBaseUrl(string) {}

func (filteringProfileTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}

func (filteringProfileTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

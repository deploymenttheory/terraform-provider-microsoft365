package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewConnectorGroupRequestInformationBuildsApplicationProxyURL(t *testing.T) {
	requestInfo, err := newConnectorGroupRequestInformation(
		context.Background(),
		testRequestAdapter{baseURL: defaultGraphBetaBaseURL},
		abstractions.GET,
		"connector-group-id",
		nil,
	)
	if err != nil {
		t.Fatalf("newConnectorGroupRequestInformation returned error: %v", err)
	}

	uri, err := requestInfo.GetUri()
	if err != nil {
		t.Fatalf("GetUri returned error: %v", err)
	}

	expected := "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/connector-group-id"
	if uri.String() != expected {
		t.Fatalf("uri = %q, expected %q", uri.String(), expected)
	}

	if requestInfo.Method != abstractions.GET {
		t.Fatalf("method = %v, expected %v", requestInfo.Method, abstractions.GET)
	}
	if values := requestInfo.Headers.Get("Accept"); len(values) != 1 || values[0] != "application/json" {
		t.Fatalf("Accept header = %#v, expected [application/json]", values)
	}
}

func TestNewConnectorGroupRequestInformationUsesConfiguredBaseURL(t *testing.T) {
	requestInfo, err := newConnectorGroupRequestInformation(
		context.Background(),
		testRequestAdapter{baseURL: "https://graph.microsoft.us/beta/"},
		abstractions.GET,
		"connector-group-id",
		nil,
	)
	if err != nil {
		t.Fatalf("newConnectorGroupRequestInformation returned error: %v", err)
	}

	uri, err := requestInfo.GetUri()
	if err != nil {
		t.Fatalf("GetUri returned error: %v", err)
	}

	expected := "https://graph.microsoft.us/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/connector-group-id"
	if uri.String() != expected {
		t.Fatalf("uri = %q, expected %q", uri.String(), expected)
	}
}

func TestNewConnectorGroupRequestInformationSerializesCreatePayload(t *testing.T) {
	name := "unit-test-connector-group"
	region := "nam"
	body := &connectorGroupRequestBody{
		name:   &name,
		region: &region,
	}

	requestInfo, err := newConnectorGroupRequestInformation(
		context.Background(),
		testRequestAdapter{baseURL: defaultGraphBetaBaseURL},
		abstractions.POST,
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newConnectorGroupRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	if payload["name"] != "unit-test-connector-group" {
		t.Fatalf("name = %#v, expected unit-test-connector-group", payload["name"])
	}
	if payload["region"] != "nam" {
		t.Fatalf("region = %#v, expected nam", payload["region"])
	}
	if _, ok := payload["connectorGroupType"]; ok {
		t.Fatalf("connectorGroupType should not be serialized: %#v", payload)
	}
	if _, ok := payload["isDefault"]; ok {
		t.Fatalf("isDefault should not be serialized: %#v", payload)
	}
}

func TestConstructUpdateResourceOmitsUnchangedRegion(t *testing.T) {
	plan := &OnPremisesConnectorGroupResourceModel{
		Name:   mustString("unit-test-connector-group-renamed"),
		Region: mustString("japan"),
	}
	state := &OnPremisesConnectorGroupResourceModel{
		Name:   mustString("unit-test-connector-group"),
		Region: mustString("japan"),
	}

	body, err := constructUpdateResource(context.Background(), plan, state)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}

	requestInfo, err := newConnectorGroupRequestInformation(
		context.Background(),
		testRequestAdapter{baseURL: defaultGraphBetaBaseURL},
		abstractions.PATCH,
		"connector-group-id",
		body,
	)
	if err != nil {
		t.Fatalf("newConnectorGroupRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	if payload["name"] != "unit-test-connector-group-renamed" {
		t.Fatalf("name = %#v, expected unit-test-connector-group-renamed", payload["name"])
	}
	if _, ok := payload["region"]; ok {
		t.Fatalf("region should be omitted when unchanged: %#v", payload)
	}
}

func mustString(value string) types.String {
	return types.StringValue(value)
}

type testRequestAdapter struct {
	baseURL string
}

func (testRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}

func (testRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (testRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}

func (testRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (testRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}

func (testRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}

func (testRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}

func (testRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}

func (testRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (a testRequestAdapter) SetBaseUrl(string) {}

func (a testRequestAdapter) GetBaseUrl() string {
	return a.baseURL
}

func (testRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

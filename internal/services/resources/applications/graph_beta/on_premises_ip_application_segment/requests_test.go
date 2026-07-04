package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"encoding/json"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewIpApplicationSegmentRequestInformationBuildsApplicationScopedURL(t *testing.T) {
	requestInfo, err := newIpApplicationSegmentRequestInformation(
		context.Background(),
		testRequestAdapter{},
		abstractions.GET,
		"application-object-id",
		"segment-id",
		nil,
	)
	if err != nil {
		t.Fatalf("newIpApplicationSegmentRequestInformation returned error: %v", err)
	}

	uri, err := requestInfo.GetUri()
	if err != nil {
		t.Fatalf("GetUri returned error: %v", err)
	}

	expected := "https://graph.microsoft.com/beta/applications/application-object-id/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/segment-id"
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

func TestNewIpApplicationSegmentRequestInformationSerializesObservedGraphPayload(t *testing.T) {
	body := &ipApplicationSegmentRequestBody{
		destinationHost: "10.10.10.10",
		destinationType: "ip",
		port:            0,
		ports:           []string{"443-443"},
		protocol:        "tcp",
	}

	requestInfo, err := newIpApplicationSegmentRequestInformation(
		context.Background(),
		testRequestAdapter{},
		abstractions.POST,
		"application-object-id",
		"",
		body,
	)
	if err != nil {
		t.Fatalf("newIpApplicationSegmentRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}

	if payload["destinationHost"] != "10.10.10.10" {
		t.Fatalf("destinationHost = %#v, expected 10.10.10.10", payload["destinationHost"])
	}
	if payload["destinationType"] != "ip" {
		t.Fatalf("destinationType = %#v, expected ip", payload["destinationType"])
	}
	if payload["port"] != float64(0) {
		t.Fatalf("port = %#v, expected 0", payload["port"])
	}
	if payload["protocol"] != "tcp" {
		t.Fatalf("protocol = %#v, expected tcp", payload["protocol"])
	}

	ports, ok := payload["ports"].([]any)
	if !ok || len(ports) != 1 || ports[0] != "443-443" {
		t.Fatalf("ports = %#v, expected [443-443]", payload["ports"])
	}
}

type testRequestAdapter struct{}

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

func (testRequestAdapter) SetBaseUrl(string) {}

func (testRequestAdapter) GetBaseUrl() string {
	return graphBetaBaseURL
}

func (testRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

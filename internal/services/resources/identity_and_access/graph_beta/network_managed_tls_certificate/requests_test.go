package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestCreateRequestSerializesPortalPayloadWithoutStatus(t *testing.T) {
	model := &NetworkManagedTLSCertificateResourceModel{
		Name:             types.StringValue("M-TLSi-test1"),
		CommonName:       types.StringValue("Microsoft Entra TLS Inspection Root CA"),
		OrganizationName: types.StringValue("Microsoft"),
		ValidityMonths:   types.Int32Value(120),
		Enabled:          types.BoolValue(true),
	}

	body, err := constructResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	requestInfo, err := newManagedTLSCertificateRequestInformation(context.Background(), managedTLSCertificateTestRequestAdapter{}, abstractions.POST, "", body)
	if err != nil {
		t.Fatalf("newManagedTLSCertificateRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	if payload["name"] != "M-TLSi-test1" || payload["commonName"] != "Microsoft Entra TLS Inspection Root CA" || payload["organizationName"] != "Microsoft" || payload["validityMonths"] != float64(120) {
		t.Fatalf("unexpected create payload: %#v", payload)
	}
	if _, exists := payload["status"]; exists {
		t.Fatalf("create payload unexpectedly contains status: %#v", payload)
	}
}

func TestStatusUpdateSerializesOnlyStatus(t *testing.T) {
	for _, test := range []struct {
		name     string
		enabled  bool
		expected string
	}{
		{name: "enable", enabled: true, expected: "enabled"},
		{name: "disable", enabled: false, expected: "disabled"},
	} {
		t.Run(test.name, func(t *testing.T) {
			body, err := constructStatusUpdate(context.Background(), test.enabled)
			if err != nil {
				t.Fatalf("constructStatusUpdate returned error: %v", err)
			}
			requestInfo, err := newManagedTLSCertificateRequestInformation(context.Background(), managedTLSCertificateTestRequestAdapter{}, abstractions.PATCH, "certificate-id", body)
			if err != nil {
				t.Fatalf("newManagedTLSCertificateRequestInformation returned error: %v", err)
			}

			var payload map[string]any
			if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
				t.Fatalf("failed to unmarshal request content: %v", err)
			}
			if len(payload) != 1 || payload["status"] != test.expected {
				t.Fatalf("unexpected status update payload: %#v", payload)
			}
		})
	}
}

func TestGeneratedNameMatchesPortalFormat(t *testing.T) {
	pattern := regexp.MustCompile(`^M-TLSi-[0-9a-z]{5}$`)
	for range 10 {
		name, err := generateManagedTLSCertificateName()
		if err != nil {
			t.Fatalf("generateManagedTLSCertificateName returned error: %v", err)
		}
		if !pattern.MatchString(name) {
			t.Fatalf("generated name %q does not match portal format", name)
		}
	}
}

type managedTLSCertificateTestRequestAdapter struct{}

func (managedTLSCertificateTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (managedTLSCertificateTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (managedTLSCertificateTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (managedTLSCertificateTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}
func (managedTLSCertificateTestRequestAdapter) SetBaseUrl(string)                            {}
func (managedTLSCertificateTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (managedTLSCertificateTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}

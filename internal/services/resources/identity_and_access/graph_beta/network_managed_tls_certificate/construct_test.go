package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
	content, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(content, &payload); err != nil {
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
	status := "enabled"
	content, err := json.Marshal(managedTLSCertificateRequest{Status: &status})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	if string(content) != `{"status":"enabled"}` {
		t.Fatalf("unexpected status update payload: %s", content)
	}
}

func TestGeneratedNameMatchesPortalFormat(t *testing.T) {
	pattern := regexp.MustCompile(`^M-TLSi-[0-9a-z]{5}$`)
	for range 10 {
		model := &NetworkManagedTLSCertificateResourceModel{Name: types.StringNull()}
		body, err := constructResource(context.Background(), model)
		if err != nil {
			t.Fatalf("constructResource returned error: %v", err)
		}
		if !pattern.MatchString(model.Name.ValueString()) || body.Name == nil || *body.Name != model.Name.ValueString() {
			t.Fatalf("generated name %q (request %v) does not match portal format", model.Name.ValueString(), body.Name)
		}
	}
}

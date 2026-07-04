package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"os"
	"testing"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestIpApplicationSegmentResponseParsesObservedGraphResponse(t *testing.T) {
	responseJSON, err := os.ReadFile("tests/responses/validate_create/post_ip_application_segment_success.json")
	if err != nil {
		t.Fatalf("failed to read response fixture: %v", err)
	}

	parseNode, err := jsonserialization.NewJsonParseNodeFactory().GetRootParseNode("application/json", responseJSON)
	if err != nil {
		t.Fatalf("GetRootParseNode returned error: %v", err)
	}

	parsed, err := parseNode.GetObjectValue(createIpApplicationSegmentResponseFromDiscriminatorValue)
	if err != nil {
		t.Fatalf("GetObjectValue returned error: %v", err)
	}

	response, ok := parsed.(*ipApplicationSegmentResponse)
	if !ok {
		t.Fatalf("parsed response is %T, expected *ipApplicationSegmentResponse", parsed)
	}

	if response.id == nil || *response.id != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("id = %#v, expected fixture id", response.id)
	}
	if response.destinationHost == nil || *response.destinationHost != "192.168.1.100" {
		t.Fatalf("destinationHost = %#v, expected 192.168.1.100", response.destinationHost)
	}
	if response.destinationType == nil || *response.destinationType != "ip" {
		t.Fatalf("destinationType = %#v, expected ip", response.destinationType)
	}
	if len(response.ports) != 2 || response.ports[0] != "80-80" || response.ports[1] != "443-443" {
		t.Fatalf("ports = %#v, expected [80-80 443-443]", response.ports)
	}
	if response.protocol == nil || *response.protocol != "tcp" {
		t.Fatalf("protocol = %#v, expected tcp", response.protocol)
	}
}

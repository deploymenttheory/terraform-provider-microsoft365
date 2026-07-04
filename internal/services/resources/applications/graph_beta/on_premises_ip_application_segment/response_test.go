package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"testing"

	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestIpApplicationSegmentResponseParsesObservedGraphResponse(t *testing.T) {
	responseJSON := []byte(`{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#applications('application-object-id')/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/$entity",
		"action": "tunnel",
		"destinationHost": "10.10.10.10",
		"destinationType": "ip",
		"exclusions": null,
		"id": "segment-id",
		"inclusions": null,
		"port": 0,
		"ports": [
			"443-443"
		],
		"protocol": "tcp"
	}`)

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

	if response.id == nil || *response.id != "segment-id" {
		t.Fatalf("id = %#v, expected segment-id", response.id)
	}
	if response.destinationHost == nil || *response.destinationHost != "10.10.10.10" {
		t.Fatalf("destinationHost = %#v, expected 10.10.10.10", response.destinationHost)
	}
	if response.destinationType == nil || *response.destinationType != "ip" {
		t.Fatalf("destinationType = %#v, expected ip", response.destinationType)
	}
	if len(response.ports) != 1 || response.ports[0] != "443-443" {
		t.Fatalf("ports = %#v, expected [443-443]", response.ports)
	}
	if response.protocol == nil || *response.protocol != "tcp" {
		t.Fatalf("protocol = %#v, expected tcp", response.protocol)
	}
}

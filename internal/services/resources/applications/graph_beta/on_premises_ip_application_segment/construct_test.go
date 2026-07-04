package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGraphDestinationType(t *testing.T) {
	tests := map[string]string{
		"ipAddress":   "ip",
		"ipRangeCidr": "ipRangeCidr",
		"fqdn":        "fqdn",
	}

	for input, expected := range tests {
		if got := graphDestinationType(input); got != expected {
			t.Fatalf("graphDestinationType(%q) = %q, expected %q", input, got, expected)
		}
	}
}

func TestTerraformDestinationType(t *testing.T) {
	tests := map[string]string{
		"ip":          "ipAddress",
		"ipRangeCidr": "ipRangeCidr",
		"fqdn":        "fqdn",
	}

	for input, expected := range tests {
		if got := terraformDestinationType(input); got != expected {
			t.Fatalf("terraformDestinationType(%q) = %q, expected %q", input, got, expected)
		}
	}
}

func TestConstructResourceMapsIpAddressForGraph(t *testing.T) {
	ports, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("443-443")})
	if diags.HasError() {
		t.Fatalf("failed to build ports set: %v", diags)
	}

	body, err := constructResource(context.Background(), &OnPremisesIpApplicationSegmentResourceModel{
		DestinationHost: types.StringValue("10.10.10.10"),
		DestinationType: types.StringValue("ipAddress"),
		Ports:           ports,
		Protocol:        types.StringValue("tcp"),
	})
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	segmentBody, ok := body.(*ipApplicationSegmentRequestBody)
	if !ok {
		t.Fatalf("constructResource returned %T, expected *ipApplicationSegmentRequestBody", body)
	}

	if segmentBody.destinationType != "ip" {
		t.Fatalf("destinationType = %q, expected %q", segmentBody.destinationType, "ip")
	}
	if segmentBody.port != 0 {
		t.Fatalf("port = %d, expected 0", segmentBody.port)
	}
	if len(segmentBody.ports) != 1 || segmentBody.ports[0] != "443-443" {
		t.Fatalf("ports = %#v, expected [443-443]", segmentBody.ports)
	}
}

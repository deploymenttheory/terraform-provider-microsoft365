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
		"ipRange":     "ipRange",
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
		"ipRange":     "ipRange",
		"ipRangeCidr": "ipRangeCidr",
		"fqdn":        "fqdn",
	}

	for input, expected := range tests {
		if got := terraformDestinationType(input); got != expected {
			t.Fatalf("terraformDestinationType(%q) = %q, expected %q", input, got, expected)
		}
	}
}

func TestValidateIpRangeHost(t *testing.T) {
	valid := []string{
		"192.168.1.1..192.168.1.10",
		"192.168.1.5..192.168.1.5",
		"10.0.0.0..10.255.255.255",
	}
	for _, host := range valid {
		if err := validateIpRangeHost(host); err != nil {
			t.Fatalf("validateIpRangeHost(%q) = %v, expected nil", host, err)
		}
	}

	invalid := []string{
		"192.168.1.0/24",            // CIDR is stored as ipRangeCidr by Graph
		"192.168.1.5",               // single IP is stored as ip by Graph
		"host.example.test",         // FQDN is stored as fqdn by Graph
		"192.168.1.1-192.168.1.10",  // hyphenated ranges return 400
		"192.168.1.10..192.168.1.1", // reversed ranges return 500
		"2001:db8::1..2001:db8::5",  // IPv6 ranges return 400
		"192.168.1.1..192.168.1.10..192.168.1.20",
		"192.168.1.1..",
		"..192.168.1.10",
		"",
	}
	for _, host := range invalid {
		if err := validateIpRangeHost(host); err == nil {
			t.Fatalf("validateIpRangeHost(%q) = nil, expected error", host)
		}
	}
}

func TestConstructResourceMapsIpAddressForGraph(t *testing.T) {
	ports, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("443-443")})
	if diags.HasError() {
		t.Fatalf("failed to build ports set: %v", diags)
	}
	protocols, diags := types.SetValue(types.StringType, []attr.Value{
		types.StringValue("udp"),
		types.StringValue("tcp"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build protocol set: %v", diags)
	}

	body, err := constructResource(context.Background(), &OnPremisesIpApplicationSegmentResourceModel{
		DestinationHost: types.StringValue("10.10.10.10"),
		DestinationType: types.StringValue("ipAddress"),
		Ports:           ports,
		Protocol:        protocols,
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
	if segmentBody.protocol != "tcp,udp" {
		t.Fatalf("protocol = %q, expected %q", segmentBody.protocol, "tcp,udp")
	}
}

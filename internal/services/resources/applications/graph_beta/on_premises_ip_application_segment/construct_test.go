package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestConstructResourceSendsGraphDestinationType(t *testing.T) {
	tests := map[string]string{
		"ipAddress":   "ip",
		"ipRange":     "ipRange",
		"ipRangeCidr": "ipRangeCidr",
		"fqdn":        "fqdn",
	}

	for input, expected := range tests {
		body, err := constructResource(context.Background(), &OnPremisesIpApplicationSegmentResourceModel{DestinationType: types.StringValue(input), Ports: types.SetNull(types.StringType), Protocol: types.SetNull(types.StringType)})
		if err != nil {
			t.Fatalf("constructResource(%q) returned error: %v", input, err)
		}
		if body.DestinationType != expected {
			t.Fatalf("destinationType for %q = %q, expected %q", input, body.DestinationType, expected)
		}
	}
}

func TestMapRemoteStateToTerraformMapsTerraformDestinationType(t *testing.T) {
	tests := map[string]string{
		"ip":          "ipAddress",
		"ipRange":     "ipRange",
		"ipRangeCidr": "ipRangeCidr",
		"fqdn":        "fqdn",
	}

	for input, expected := range tests {
		data := &OnPremisesIpApplicationSegmentResourceModel{}
		MapRemoteStateToTerraform(context.Background(), data, &ipApplicationSegmentResponse{DestinationType: &input})
		if data.DestinationType.ValueString() != expected {
			t.Fatalf("destination_type for %q = %q, expected %q", input, data.DestinationType.ValueString(), expected)
		}
	}
}

func TestMapRemoteStateToTerraformParsesObservedGraphResponse(t *testing.T) {
	responseJSON, err := os.ReadFile("tests/responses/validate_get/get_ip_application_segment_maximal.json")
	if err != nil {
		t.Fatalf("failed to read response fixture: %v", err)
	}

	var remote ipApplicationSegmentResponse
	if err := json.Unmarshal(responseJSON, &remote); err != nil {
		t.Fatalf("failed to decode response fixture: %v", err)
	}

	data := &OnPremisesIpApplicationSegmentResourceModel{}
	MapRemoteStateToTerraform(context.Background(), data, &remote)

	if data.ID.IsNull() || data.ID.ValueString() == "" {
		t.Fatal("id was not mapped")
	}
	if data.DestinationHost.IsNull() || data.DestinationType.IsNull() {
		t.Fatalf("destination was not mapped: host %s, type %s", data.DestinationHost, data.DestinationType)
	}
	if len(data.Ports.Elements()) == 0 {
		t.Fatal("ports were not mapped")
	}
	if len(data.Protocol.Elements()) == 0 {
		t.Fatal("protocol was not mapped")
	}
}

func TestValidateConfigRequiresIpRangeStartEndHost(t *testing.T) {
	valid := []string{
		"192.168.1.1..192.168.1.10",
		"192.168.1.5..192.168.1.5",
		"10.0.0.0..10.255.255.255",
	}
	for _, host := range valid {
		if diags := validateIpRangeConfig(t, host); diags.HasError() {
			t.Fatalf("ValidateConfig(%q) = %v, expected no error", host, diags)
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
		if diags := validateIpRangeConfig(t, host); !diags.HasError() {
			t.Fatalf("ValidateConfig(%q) returned no error, expected one", host)
		}
	}
}

// validateIpRangeConfig runs ValidateConfig for an ipRange segment with destinationHost host.
func validateIpRangeConfig(t *testing.T, host string) diag.Diagnostics {
	t.Helper()
	ctx := context.Background()
	r := &OnPremisesIpApplicationSegmentResource{}

	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("failed to obtain schema: %v", schemaResponse.Diagnostics)
	}

	stringSet := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("80-80")})
	protocolSet := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("tcp")})
	config := tfsdk.Config{Schema: schemaResponse.Schema, Raw: tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil)}
	model := OnPremisesIpApplicationSegmentResourceModel{
		ID:                  types.StringNull(),
		ApplicationObjectID: types.StringValue("12345678-1234-1234-1234-123456789012"),
		DestinationHost:     types.StringValue(host),
		DestinationType:     types.StringValue("ipRange"),
		Ports:               stringSet,
		Protocol:            protocolSet,
		Timeouts:            timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})},
	}
	state := tfsdk.State{Schema: schemaResponse.Schema, Raw: config.Raw}
	if setDiags := state.Set(ctx, &model); setDiags.HasError() {
		t.Fatalf("failed to build config: %v", setDiags)
	}
	config.Raw = state.Raw

	var response resource.ValidateConfigResponse
	r.ValidateConfig(ctx, resource.ValidateConfigRequest{Config: config}, &response)
	return response.Diagnostics
}

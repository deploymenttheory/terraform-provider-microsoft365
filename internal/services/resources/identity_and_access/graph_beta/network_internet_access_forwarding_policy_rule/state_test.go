package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRemoteStateToTerraformInternetAccessForwardingRule(t *testing.T) {
	id := "66f1c2ee-9f17-4681-a7bc-c324e7dff554"
	name := "Custom Acquire policy internet rule"
	ruleType := "fqdn"
	action := "forward"
	clientFallbackAction := "block"
	protocol := "udp"
	odataType := fqdnODataType
	value := "example.com"

	data := &NetworkInternetAccessForwardingPolicyRuleResourceModel{}
	MapRemoteStateToTerraform(context.Background(), data, &internetAccessForwardingRuleResponse{
		id:                   &id,
		name:                 &name,
		ruleType:             &ruleType,
		action:               &action,
		clientFallbackAction: &clientFallbackAction,
		ports:                []string{"80", "443"},
		protocol:             &protocol,
		destinations: []ruleDestinationResponse{
			{odataType: &odataType, value: &value},
		},
	})

	if data.ID.ValueString() != id {
		t.Fatalf("id = %q, want %q", data.ID.ValueString(), id)
	}
	if data.ClientFallbackAction.ValueString() != "block" {
		t.Fatalf("client_fallback_action = %q, want block", data.ClientFallbackAction.ValueString())
	}
	if data.RuleType.ValueString() != ruleType {
		t.Fatalf("rule_type = %q, want %q", data.RuleType.ValueString(), ruleType)
	}
	if data.Ports.Elements()[0].String() != "\"80\"" || data.Ports.Elements()[1].String() != "\"443\"" {
		t.Fatalf("ports = %#v, want 80 and 443", data.Ports.Elements())
	}
	destinations := data.Destinations.Elements()
	if len(destinations) != 1 {
		t.Fatalf("destinations length = %d, want 1", len(destinations))
	}
	destination := destinations[0].(types.Object)
	attrs := destination.Attributes()
	if attrs["type"].(types.String).ValueString() != ruleTypeFQDN {
		t.Fatalf("destination type = %q, want fqdn", attrs["type"].(types.String).ValueString())
	}
	if attrs["value"].(types.String).ValueString() != value {
		t.Fatalf("destination value = %q, want %q", attrs["value"].(types.String).ValueString(), value)
	}
}

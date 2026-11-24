package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func TestMapRemoteStateToTerraformInternetAccessForwardingRule(t *testing.T) {
	id := "66f1c2ee-9f17-4681-a7bc-c324e7dff554"
	name := "Custom Acquire policy internet rule"
	ruleType := networkaccess.FQDN_NETWORKDESTINATIONTYPE
	action := networkaccess.FORWARD_FORWARDINGRULEACTION
	clientFallbackAction := networkaccess.BLOCK_CLIENTFALLBACKACTION
	protocol := networkaccess.UDP_NETWORKINGPROTOCOL
	value := "example.com"

	fqdn := networkaccess.NewFqdn()
	fqdn.SetValue(&value)
	remote := networkaccess.NewInternetAccessForwardingRule()
	remote.SetId(&id)
	remote.SetName(&name)
	remote.SetRuleType(&ruleType)
	remote.SetAction(&action)
	remote.SetClientFallbackAction(&clientFallbackAction)
	remote.SetPorts([]string{"80", "443"})
	remote.SetProtocol(&protocol)
	remote.SetDestinations([]networkaccess.RuleDestinationable{fqdn})

	data := &NetworkInternetAccessForwardingPolicyRuleResourceModel{}
	MapRemoteStateToTerraform(context.Background(), data, remote)

	if data.ID.ValueString() != id {
		t.Fatalf("id = %q, want %q", data.ID.ValueString(), id)
	}
	if data.ClientFallbackAction.ValueString() != "block" {
		t.Fatalf("client_fallback_action = %q, want block", data.ClientFallbackAction.ValueString())
	}
	if data.RuleType.ValueString() != ruleTypeFQDN {
		t.Fatalf("rule_type = %q, want %q", data.RuleType.ValueString(), ruleTypeFQDN)
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

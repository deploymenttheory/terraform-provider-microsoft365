package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"testing"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func TestMapRemoteStateToTerraformForwardingPolicyLink(t *testing.T) {
	linkID := "09837256-2cba-4dde-a121-4d6a129f13db"
	version := "1.0.0"
	priority := int64(4)
	state := models.ENABLED_STATUS
	policyID := "f0474b3e-307a-4230-bc1c-cd8ac2f1a2cf"
	policyName := "Default Acquire"
	policyDescription := "Internet Access default acquire forwarding policy"
	trafficForwardingType := models.INTERNET_TRAFFICFORWARDINGTYPE

	policy := models.NewForwardingPolicy()
	policy.SetId(&policyID)
	policy.SetName(&policyName)
	policy.SetDescription(&policyDescription)
	policy.SetTrafficForwardingType(&trafficForwardingType)

	link := models.NewForwardingPolicyLink()
	link.SetId(&linkID)
	link.SetVersion(&version)
	link.SetPriority(&priority)
	link.SetState(&state)
	link.SetPolicy(policy)

	data := &NetworkForwardingProfilePolicyLinkResourceModel{}
	MapRemoteStateToTerraform(context.Background(), data, "72661c0d-027e-4dff-8c76-af103f200903", link)

	if data.ID.ValueString() != "72661c0d-027e-4dff-8c76-af103f200903/09837256-2cba-4dde-a121-4d6a129f13db" {
		t.Fatalf("id = %q", data.ID.ValueString())
	}
	if data.State.ValueString() != "enabled" {
		t.Fatalf("state = %q, want enabled", data.State.ValueString())
	}
	if data.PolicyID.ValueString() != policyID {
		t.Fatalf("policy_id = %q, want %q", data.PolicyID.ValueString(), policyID)
	}
	if data.PolicyLinkID.ValueString() == data.PolicyID.ValueString() {
		t.Fatalf("policy_link_id and policy_id must remain distinct")
	}
	if data.Priority.ValueInt64() != priority {
		t.Fatalf("priority = %d, want %d", data.Priority.ValueInt64(), priority)
	}
	if data.TrafficForwardingType.ValueString() != "internet" {
		t.Fatalf("traffic_forwarding_type = %q, want internet", data.TrafficForwardingType.ValueString())
	}
}

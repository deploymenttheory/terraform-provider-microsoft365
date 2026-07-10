package graphBetaNetworkForwardingProfile

import (
	"testing"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func TestMapRemoteStateToDataSourceForwardingProfile(t *testing.T) {
	profileID := "72661c0d-027e-4dff-8c76-af103f200903"
	profileName := "Internet traffic forwarding profile"
	profileDescription := "Default Internet Access forwarding profile"
	profileVersion := "1.0.0"
	profileState := models.DISABLED_STATUS
	profilePriority := int32(300)
	isCustomProfile := false
	trafficForwardingType := models.INTERNET_TRAFFICFORWARDINGTYPE

	policyLinkID := "09837256-2cba-4dde-a121-4d6a129f13db"
	policyLinkState := models.ENABLED_STATUS
	policyPriority := int64(4)
	policyID := "f0474b3e-307a-4230-bc1c-cd8ac2f1a2cf"
	policyName := "Default Acquire"
	policyDescription := "Default acquire policy"

	policy := models.NewForwardingPolicy()
	policy.SetId(&policyID)
	policy.SetName(&policyName)
	policy.SetDescription(&policyDescription)
	policy.SetTrafficForwardingType(&trafficForwardingType)

	link := models.NewForwardingPolicyLink()
	link.SetId(&policyLinkID)
	link.SetState(&policyLinkState)
	link.SetPriority(&policyPriority)
	link.SetPolicy(policy)

	profile := models.NewForwardingProfile()
	profile.SetId(&profileID)
	profile.SetName(&profileName)
	profile.SetDescription(&profileDescription)
	profile.SetVersion(&profileVersion)
	profile.SetState(&profileState)
	profile.SetPriority(&profilePriority)
	profile.SetIsCustomProfile(&isCustomProfile)
	profile.SetTrafficForwardingType(&trafficForwardingType)
	profile.SetPolicies([]models.PolicyLinkable{link})
	profile.GetAdditionalData()["clientFallbackAction"] = "bypass"

	mapped := MapRemoteStateToDataSource(profile)
	if mapped.ID.ValueString() != profileID {
		t.Fatalf("id = %q, want %q", mapped.ID.ValueString(), profileID)
	}
	if mapped.Name.ValueString() != profileName {
		t.Fatalf("name = %q, want %q", mapped.Name.ValueString(), profileName)
	}
	if mapped.State.ValueString() != "disabled" {
		t.Fatalf("state = %q, want disabled", mapped.State.ValueString())
	}
	if mapped.TrafficForwardingType.ValueString() != "internet" {
		t.Fatalf("traffic_forwarding_type = %q, want internet", mapped.TrafficForwardingType.ValueString())
	}
	if mapped.ClientFallbackAction.ValueString() != "bypass" {
		t.Fatalf("client_fallback_action = %q, want bypass", mapped.ClientFallbackAction.ValueString())
	}
	if len(mapped.Policies) != 1 {
		t.Fatalf("policies length = %d, want 1", len(mapped.Policies))
	}
	if mapped.Policies[0].PolicyLinkID.ValueString() == mapped.Policies[0].PolicyID.ValueString() {
		t.Fatalf("policy_link_id and policy_id must remain distinct")
	}
	if mapped.Policies[0].Priority.ValueInt64() != policyPriority {
		t.Fatalf("policy priority = %d, want %d", mapped.Policies[0].Priority.ValueInt64(), policyPriority)
	}
}

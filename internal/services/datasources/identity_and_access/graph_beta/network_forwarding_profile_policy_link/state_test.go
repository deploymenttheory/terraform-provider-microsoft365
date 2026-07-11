package graphBetaNetworkForwardingProfilePolicyLink

import (
	"testing"

	profileds "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/network_forwarding_profile"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPolicyLinkMatchToModel(t *testing.T) {
	match := policyLinkMatch{
		profile: profileds.ForwardingProfileModel{
			ID:                    types.StringValue("72661c0d-027e-4dff-8c76-af103f200903"),
			Name:                  types.StringValue("Internet traffic forwarding profile"),
			TrafficForwardingType: types.StringValue("internet"),
		},
		link: profileds.ForwardingProfilePolicyLinkModel{
			PolicyLinkID:          types.StringValue("f576d498-0067-4cc8-960b-b6e3ebf571ea"),
			PolicyID:              types.StringValue("dad2a411-e330-440d-a7c7-2c830dce5991"),
			PolicyName:            types.StringValue("Custom Acquire"),
			State:                 types.StringValue("enabled"),
			TrafficForwardingType: types.StringValue("internet"),
		},
	}

	model := match.toModel()
	if model.ID.ValueString() != "72661c0d-027e-4dff-8c76-af103f200903/f576d498-0067-4cc8-960b-b6e3ebf571ea" {
		t.Fatalf("id = %q", model.ID.ValueString())
	}
	if model.ForwardingProfileID.ValueString() != match.profile.ID.ValueString() {
		t.Fatalf("forwarding_profile_id = %q", model.ForwardingProfileID.ValueString())
	}
	if model.PolicyLinkID.ValueString() == model.PolicyID.ValueString() {
		t.Fatalf("policy_link_id and policy_id must remain distinct")
	}
}

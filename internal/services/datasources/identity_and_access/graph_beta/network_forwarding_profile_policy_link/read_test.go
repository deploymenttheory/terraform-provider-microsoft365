package graphBetaNetworkForwardingProfilePolicyLink

import (
	"testing"

	profileds "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/network_forwarding_profile"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDetermineLookupMethod(t *testing.T) {
	tests := []struct {
		name     string
		model    NetworkForwardingProfilePolicyLinkDataSourceModel
		expected lookupMethod
	}{
		{
			name: "by forwarding profile id",
			model: NetworkForwardingProfilePolicyLinkDataSourceModel{
				ForwardingProfileID: types.StringValue("72661c0d-027e-4dff-8c76-af103f200903"),
			},
			expected: lookupByForwardingProfileID,
		},
		{
			name: "by forwarding profile name",
			model: NetworkForwardingProfilePolicyLinkDataSourceModel{
				ForwardingProfileName: types.StringValue("Internet traffic forwarding profile"),
			},
			expected: lookupByForwardingProfileName,
		},
		{
			name: "by traffic forwarding type",
			model: NetworkForwardingProfilePolicyLinkDataSourceModel{
				TrafficForwardingType: types.StringValue("internet"),
			},
			expected: lookupByTrafficForwardingType,
		},
		{
			name:     "unset",
			model:    NetworkForwardingProfilePolicyLinkDataSourceModel{},
			expected: lookupUnset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := determineLookupMethod(tt.model); actual != tt.expected {
				t.Fatalf("determineLookupMethod() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestFindPolicyLinkMatches(t *testing.T) {
	profiles := []profileds.ForwardingProfileModel{
		{
			ID:                    types.StringValue("72661c0d-027e-4dff-8c76-af103f200903"),
			Name:                  types.StringValue("Internet traffic forwarding profile"),
			TrafficForwardingType: types.StringValue("internet"),
			Policies: []profileds.ForwardingProfilePolicyLinkModel{
				{
					PolicyLinkID: types.StringValue("f576d498-0067-4cc8-960b-b6e3ebf571ea"),
					PolicyID:     types.StringValue("dad2a411-e330-440d-a7c7-2c830dce5991"),
					PolicyName:   types.StringValue("Custom Acquire"),
				},
			},
		},
	}

	matches := findPolicyLinkMatches(profiles, "custom acquire")
	if len(matches) != 1 {
		t.Fatalf("matches length = %d, want 1", len(matches))
	}
	if matches[0].link.PolicyLinkID.ValueString() != "f576d498-0067-4cc8-960b-b6e3ebf571ea" {
		t.Fatalf("policy link id = %q", matches[0].link.PolicyLinkID.ValueString())
	}
}

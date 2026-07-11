package graphBetaNetworkForwardingProfile

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDetermineLookupMethod(t *testing.T) {
	tests := []struct {
		name     string
		model    NetworkForwardingProfileDataSourceModel
		expected lookupMethod
	}{
		{
			name: "by forwarding profile id",
			model: NetworkForwardingProfileDataSourceModel{
				ForwardingProfileID: types.StringValue("72661c0d-027e-4dff-8c76-af103f200903"),
			},
			expected: lookupByForwardingProfileID,
		},
		{
			name: "by name",
			model: NetworkForwardingProfileDataSourceModel{
				Name: types.StringValue("Internet traffic forwarding profile"),
			},
			expected: lookupByName,
		},
		{
			name: "by traffic forwarding type",
			model: NetworkForwardingProfileDataSourceModel{
				TrafficForwardingType: types.StringValue("internet"),
			},
			expected: lookupByTrafficForwardingType,
		},
		{
			name: "list all",
			model: NetworkForwardingProfileDataSourceModel{
				ListAll: types.BoolValue(true),
			},
			expected: lookupListAll,
		},
		{
			name:     "unset",
			model:    NetworkForwardingProfileDataSourceModel{},
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

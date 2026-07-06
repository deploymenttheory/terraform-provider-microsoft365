package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"testing"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func TestMapRemoteStateToTerraformKnownPolicyTypes(t *testing.T) {
	tests := []struct {
		name               string
		linkODataType      string
		policyODataType    string
		expectedPolicyType string
	}{
		{
			name:               "legacy filtering",
			linkODataType:      filteringPolicyLinkODataType,
			policyODataType:    filteringPolicyODataType,
			expectedPolicyType: policyTypeFiltering,
		},
		{
			name:               "web filtering",
			linkODataType:      webFilteringPolicyLinkODataType,
			policyODataType:    webFilteringPolicyODataType,
			expectedPolicyType: policyTypeWebFiltering,
		},
		{
			name:               "cloud firewall",
			linkODataType:      cloudFirewallPolicyLinkODataType,
			policyODataType:    cloudFirewallPolicyODataType,
			expectedPolicyType: policyTypeCloudFirewall,
		},
		{
			name:               "threat intelligence",
			linkODataType:      threatIntelligencePolicyLinkODataType,
			policyODataType:    threatIntelligencePolicyODataType,
			expectedPolicyType: policyTypeThreatIntelligence,
		},
		{
			name:               "tls inspection",
			linkODataType:      tlsInspectionPolicyLinkODataType,
			policyODataType:    tlsInspectionPolicyODataType,
			expectedPolicyType: policyTypeTlsInspection,
		},
		{
			name:               "custom",
			linkODataType:      "#microsoft.graph.networkaccess.examplePolicyLink",
			policyODataType:    "#microsoft.graph.networkaccess.examplePolicy",
			expectedPolicyType: policyTypeCustom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := models.NewPolicyLink()
			linkID := "link-id"
			link.SetId(&linkID)
			link.SetOdataType(&tt.linkODataType)

			policy := models.NewPolicy()
			policyID := "policy-id"
			policy.SetId(&policyID)
			policy.SetOdataType(&tt.policyODataType)
			link.SetPolicy(policy)

			data := &NetworkFilteringProfilePolicyLinkResourceModel{}
			MapRemoteStateToTerraform(context.Background(), data, "profile-id", link)

			if data.PolicyType.ValueString() != tt.expectedPolicyType {
				t.Fatalf("policy_type = %q, want %q", data.PolicyType.ValueString(), tt.expectedPolicyType)
			}
			if data.PolicyLinkODataType.ValueString() != tt.linkODataType {
				t.Fatalf("policy_link_odata_type = %q, want %q", data.PolicyLinkODataType.ValueString(), tt.linkODataType)
			}
			if data.PolicyODataType.ValueString() != tt.policyODataType {
				t.Fatalf("policy_odata_type = %q, want %q", data.PolicyODataType.ValueString(), tt.policyODataType)
			}
			if data.ID.ValueString() != "profile-id/link-id" {
				t.Fatalf("id = %q, want profile-id/link-id", data.ID.ValueString())
			}
			if !data.Priority.IsNull() {
				t.Fatalf("priority = %#v, want null for base policy link", data.Priority)
			}
		})
	}
}

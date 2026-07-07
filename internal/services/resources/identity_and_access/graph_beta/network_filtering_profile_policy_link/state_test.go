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
			name:               "prompt",
			linkODataType:      promptPolicyLinkODataType,
			policyODataType:    promptPolicyODataType,
			expectedPolicyType: policyTypePrompt,
		},
		{
			name:               "content",
			linkODataType:      filePolicyLinkODataType,
			policyODataType:    filePolicyODataType,
			expectedPolicyType: policyTypeContent,
		},
		{
			name:               "netskope dlp",
			linkODataType:      securityProviderPolicyLinkODataType,
			policyODataType:    securityProviderPolicyODataType,
			expectedPolicyType: policyTypeNetskopeDlp,
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
			if data.ID.ValueString() != "profile-id/link-id" {
				t.Fatalf("id = %q, want profile-id/link-id", data.ID.ValueString())
			}
			if !data.Priority.IsNull() {
				t.Fatalf("priority = %#v, want null for base policy link", data.Priority)
			}
			if !data.CreatedDateTime.IsNull() {
				t.Fatalf("created_date_time = %#v, want null for base policy link", data.CreatedDateTime)
			}
			if !data.LastModifiedDateTime.IsNull() {
				t.Fatalf("last_modified_date_time = %#v, want null for base policy link", data.LastModifiedDateTime)
			}
		})
	}
}

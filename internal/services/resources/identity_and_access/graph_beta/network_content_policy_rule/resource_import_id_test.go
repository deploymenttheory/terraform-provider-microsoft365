package graphBetaNetworkContentPolicyRule_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkContentPolicyRule "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_content_policy_rule"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceNetworkContentPolicyRule_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceNetworkContentPolicyRule_ImportRejectsInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()

	tests := []struct {
		name     string
		importID string
	}{
		{name: "content policy ID not a GUID", importID: "not-a-guid/00000000-0000-0000-0000-000000000302"},
		{name: "rule ID with query string", importID: "00000000-0000-0000-0000-000000000301/00000000-0000-0000-0000-000000000302?x=1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:        `resource "` + graphBetaNetworkContentPolicyRule.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaNetworkContentPolicyRule.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}

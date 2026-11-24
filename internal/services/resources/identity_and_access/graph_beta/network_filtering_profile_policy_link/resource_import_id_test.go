package graphBetaNetworkFilteringProfilePolicyLink_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkFilteringProfilePolicyLink "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_profile_policy_link"
)

// TestUnitResourceNetworkFilteringProfilePolicyLink_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceNetworkFilteringProfilePolicyLink_ImportRejectsInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()

	tests := []struct {
		name     string
		importID string
	}{
		{name: "filtering profile ID not a GUID", importID: "not-a-guid/00000000-0000-0000-0000-000000000002"},
		{name: "policy link ID with query string", importID: "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002?x=1"},
		{name: "policy link ID with fragment", importID: "00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002#x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:        `resource "` + graphBetaNetworkFilteringProfilePolicyLink.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaNetworkFilteringProfilePolicyLink.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}

package graphBetaAgentIdentityBlueprintFederatedIdentityCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentityBlueprintFederatedIdentityCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_federated_identity_credential"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceAgentIdentityBlueprintFederatedIdentityCredential_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceAgentIdentityBlueprintFederatedIdentityCredential_ImportRejectsInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()

	tests := []struct {
		name     string
		importID string
	}{
		{name: "blueprint ID not a GUID", importID: "not-a-guid/11111111-1111-1111-1111-111111111111"},
		{name: "credential ID with query string", importID: "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111?x=1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:        `resource "` + graphBetaAgentIdentityBlueprintFederatedIdentityCredential.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaAgentIdentityBlueprintFederatedIdentityCredential.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}

package graphBetaAgentIdentityBlueprintIdentifierUri_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentityBlueprintIdentifierUri "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_identifier_uri"
	identifierUriMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_identifier_uri/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentityBlueprintIdentifierUri.ResourceName

	// testResource is the test resource implementation
	testResource = graphBetaAgentIdentityBlueprintIdentifierUri.AgentIdentityBlueprintIdentifierUriTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *identifierUriMocks.AgentIdentityBlueprintIdentifierUriMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	identifierUriMock := &identifierUriMocks.AgentIdentityBlueprintIdentifierUriMock{}
	identifierUriMock.RegisterMocks()

	return mockClient, identifierUriMock
}

func TestUnitResourceAgentIdentityBlueprintIdentifierUri_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, identifierUriMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer identifierUriMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("identifier_uri").MatchesRegex(regexp.MustCompile(`^api://`)),
					check.That(resourceType+".test_minimal").Key("scope.value").HasValue("access_agent"),
				),
			},
			{
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					blueprintID := rs.Primary.Attributes["blueprint_id"]
					identifierUri := rs.Primary.Attributes["identifier_uri"]
					return fmt.Sprintf("%s/%s", blueprintID, identifierUri), nil
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "blueprint_id",
			},
		},
	})
}

func testConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

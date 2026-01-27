package graphBetaApplicationsAgentIdentityBlueprint_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationsAgentIdentityBlueprint "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint"
	agentIdentityBlueprintMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationsAgentIdentityBlueprint.ResourceName

	// testResource is the test resource implementation for agent identity blueprints
	testResource = graphBetaApplicationsAgentIdentityBlueprint.AgentIdentityBlueprintTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintMocks.AgentIdentityBlueprintMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintMock := &agentIdentityBlueprintMocks.AgentIdentityBlueprintMock{}
	agentIdentityBlueprintMock.RegisterMocks()
	return mockClient, agentIdentityBlueprintMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintMocks.AgentIdentityBlueprintMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintMock := &agentIdentityBlueprintMocks.AgentIdentityBlueprintMock{}
	agentIdentityBlueprintMock.RegisterErrorMocks()
	return mockClient, agentIdentityBlueprintMock
}

func TestUnitResourceAgentIdentityBlueprint_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("unit-test-agent-identity-blueprint-minimal"),
					check.That(resourceType+".test_minimal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
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
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceAgentIdentityBlueprint_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").HasValue("unit-test-agent-identity-blueprint-maximal"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("This is a test agent identity blueprint with all optional fields configured"),
					check.That(resourceType+".test_maximal").Key("sign_in_audience").HasValue("AzureADMyOrg"),

					// Sponsors
					check.That(resourceType+".test_maximal").Key("sponsor_user_ids.#").HasValue("2"),

					// Owners
					check.That(resourceType+".test_maximal").Key("owner_user_ids.#").HasValue("2"),

					// Tags
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("terraform"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("test"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("agent"),
				),
			},
			{
				ResourceName: resourceType + ".test_maximal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_maximal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_maximal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceAgentIdentityBlueprint_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("display_name").HasValue("unit-test-agent-identity-blueprint-minimal"),
				),
			},
			{
				Config: testConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("unit-test-agent-identity-blueprint-minimal-updated"),
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description"),
				),
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

func testConfigMaximal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testConfigMinimalUpdated() string {
	return `
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_minimal" {
  display_name     = "unit-test-agent-identity-blueprint-minimal-updated"
  description      = "Updated description"
  sponsor_user_ids = ["11111111-1111-1111-1111-111111111111"]
  owner_user_ids   = ["11111111-1111-1111-1111-111111111111"]
}
`
}

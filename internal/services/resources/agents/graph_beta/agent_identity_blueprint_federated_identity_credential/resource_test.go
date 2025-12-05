package graphBetaAgentIdentityBlueprintFederatedIdentityCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentityBlueprintFederatedIdentityCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_federated_identity_credential"
	agentIdentityBlueprintFederatedIdentityCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_federated_identity_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentityBlueprintFederatedIdentityCredential.ResourceName

	// testResource is the test resource implementation for federated identity credentials
	testResource = graphBetaAgentIdentityBlueprintFederatedIdentityCredential.AgentIdentityBlueprintFederatedIdentityCredentialTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintFederatedIdentityCredentialMocks.AgentIdentityBlueprintFederatedIdentityCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintFederatedIdentityCredentialMock := &agentIdentityBlueprintFederatedIdentityCredentialMocks.AgentIdentityBlueprintFederatedIdentityCredentialMock{}
	agentIdentityBlueprintFederatedIdentityCredentialMock.RegisterMocks()
	return mockClient, agentIdentityBlueprintFederatedIdentityCredentialMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintFederatedIdentityCredentialMocks.AgentIdentityBlueprintFederatedIdentityCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintFederatedIdentityCredentialMock := &agentIdentityBlueprintFederatedIdentityCredentialMocks.AgentIdentityBlueprintFederatedIdentityCredentialMock{}
	agentIdentityBlueprintFederatedIdentityCredentialMock.RegisterErrorMocks()
	return mockClient, agentIdentityBlueprintFederatedIdentityCredentialMock
}

func TestAgentIdentityBlueprintFederatedIdentityCredentialResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("name").HasValue("unit-test-fic-minimal"),
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_minimal").Key("subject").HasValue("repo:octo-org/octo-repo:environment:Production"),
					check.That(resourceType+".test_minimal").Key("audiences.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_minimal"),
			},
		},
	})
}

func TestAgentIdentityBlueprintFederatedIdentityCredentialResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_maximal").Key("name").HasValue("unit-test-fic-maximal"),
					check.That(resourceType+".test_maximal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_maximal").Key("subject").HasValue("repo:octo-org/octo-repo:environment:Production"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("This is a test federated identity credential with all optional fields configured"),
					check.That(resourceType+".test_maximal").Key("audiences.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_maximal"),
			},
		},
	})
}

func TestAgentIdentityBlueprintFederatedIdentityCredentialResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
				),
			},
			{
				Config: testConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description for unit test"),
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
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "test_minimal" {
  blueprint_id = "11111111-1111-1111-1111-111111111111"
  name         = "unit-test-fic-minimal"
  issuer       = "https://token.actions.githubusercontent.com"
  subject      = "repo:octo-org/octo-repo:environment:Production"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Updated description for unit test"
}
`
}

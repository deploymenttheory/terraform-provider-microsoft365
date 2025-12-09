package graphBetaAgentIdentity_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentity "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity"
	agentIdentityMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentity.ResourceName

	// testResource is the test resource implementation for agent identities
	testResource = graphBetaAgentIdentity.AgentIdentityTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityMocks.AgentIdentityMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityMock := &agentIdentityMocks.AgentIdentityMock{}
	agentIdentityMock.RegisterMocks()
	return mockClient, agentIdentityMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentIdentityMocks.AgentIdentityMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityMock := &agentIdentityMocks.AgentIdentityMock{}
	agentIdentityMock.RegisterErrorMocks()
	return mockClient, agentIdentityMock
}

func TestUnitAgentIdentityResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("Unit Test Agent Identity"),
					check.That(resourceType+".test_minimal").Key("agent_identity_blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("service_principal_type").HasValue("ServiceIdentity"),
					check.That(resourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("sponsor_ids.#").HasValue("1"),
					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
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
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					return fmt.Sprintf("%s/%s", id, blueprintId), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitAgentIdentityResource_WithTags(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithTags(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_with_tags").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_with_tags").Key("display_name").HasValue("Unit Test Agent Identity With Tags"),
					check.That(resourceType+".test_with_tags").Key("agent_identity_blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_with_tags").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_with_tags").Key("sponsor_ids.#").HasValue("1"),
					check.That(resourceType+".test_with_tags").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName: resourceType + ".test_with_tags",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_with_tags"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_with_tags")
					}
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					return fmt.Sprintf("%s/%s", id, blueprintId), nil
				},
				ImportStateVerify: true,
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

func testConfigWithTags() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_tags.tf")
	if err != nil {
		panic(err)
	}
	return content
}

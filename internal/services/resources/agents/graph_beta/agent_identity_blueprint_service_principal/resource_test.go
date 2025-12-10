package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationsAgentIdentityBlueprintServicePrincipal "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_service_principal"
	agentIdentityBlueprintServicePrincipalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_service_principal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationsAgentIdentityBlueprintServicePrincipal.ResourceName

	// testResource is the test resource implementation for agent identity blueprint service principals
	testResource = graphBetaApplicationsAgentIdentityBlueprintServicePrincipal.AgentIdentityBlueprintServicePrincipalTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintServicePrincipalMocks.AgentIdentityBlueprintServicePrincipalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintServicePrincipalMock := &agentIdentityBlueprintServicePrincipalMocks.AgentIdentityBlueprintServicePrincipalMock{}
	agentIdentityBlueprintServicePrincipalMock.RegisterMocks()
	return mockClient, agentIdentityBlueprintServicePrincipalMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintServicePrincipalMocks.AgentIdentityBlueprintServicePrincipalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentIdentityBlueprintServicePrincipalMock := &agentIdentityBlueprintServicePrincipalMocks.AgentIdentityBlueprintServicePrincipalMock{}
	agentIdentityBlueprintServicePrincipalMock.RegisterErrorMocks()
	return mockClient, agentIdentityBlueprintServicePrincipalMock
}

func TestAgentIdentityBlueprintServicePrincipalResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentIdentityBlueprintServicePrincipalMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentIdentityBlueprintServicePrincipalMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").HasValue("11111111-1111-1111-1111-111111111111"),
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

func testConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

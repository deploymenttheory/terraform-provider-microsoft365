package graphBetaAgentInstance_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentInstance "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_instance"
	agentInstanceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_instance/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaAgentInstance.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *agentInstanceMocks.AgentInstanceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentInstanceMock := &agentInstanceMocks.AgentInstanceMock{}
	agentInstanceMock.RegisterMocks()
	return mockClient, agentInstanceMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentInstanceMocks.AgentInstanceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentInstanceMock := &agentInstanceMocks.AgentInstanceMock{}
	agentInstanceMock.RegisterErrorMocks()
	return mockClient, agentInstanceMock
}

// TestUnitResourceAgentInstance_01_Minimal tests creating an agent instance with minimal configuration
func TestUnitResourceAgentInstance_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentInstanceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentInstanceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("Unit Test Agent Instance Minimal"),
					check.That(resourceType+".test_minimal").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.display_name").HasValue("Unit Test Agent Card Minimal"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.description").HasValue("Minimal unit test agent card manifest description"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.protocol_version").HasValue("1.0"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.version").HasValue("1.0.0"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.supports_authenticated_extended_card").HasValue("false"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.streaming").HasValue("false"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.push_notifications").HasValue("false"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.state_transition_history").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitResourceAgentInstance_02_Maximal tests creating an agent instance with maximal configuration
func TestUnitResourceAgentInstance_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentInstanceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentInstanceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").HasValue("IT Service Desk Agent Maximal"),
					check.That(resourceType+".test_maximal").Key("originating_store").HasValue("Deployment Theory"),
					check.That(resourceType+".test_maximal").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("url").HasValue("https://servicedesk.deploymenttheory.com/api"),
					check.That(resourceType+".test_maximal").Key("preferred_transport").HasValue("HTTP+JSON"),
					check.That(resourceType+".test_maximal").Key("additional_interfaces.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.display_name").HasValue("IT Service Desk Agent"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.description").Exists(),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.protocol_version").HasValue("1.0"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.version").HasValue("2.0.0"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.icon_url").HasValue("https://servicedesk.example.com/assets/agent-icon.png"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.documentation_url").HasValue("https://docs.example.com/servicedesk-agent"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.supports_authenticated_extended_card").HasValue("false"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.default_input_modes.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.default_output_modes.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.provider.organization").HasValue("Deployment Theory"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.provider.url").HasValue("https://www.deploymenttheory.com"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.push_notifications").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.state_transition_history").HasValue("false"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.extensions.#").HasValue("1"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitResourceAgentInstance_03_UpdateMinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitResourceAgentInstance_03_UpdateMinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentInstanceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentInstanceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with minimal configuration
			{
				Config: testConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Update Test Agent Minimal"),
					check.That(resourceType+".test_update").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.display_name").HasValue("Update Test Agent Card Minimal"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("1.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("false"),
				),
			},
			// Step 2: Update to maximal configuration
			{
				Config: testConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Update Test Agent Maximal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("url").HasValue("https://updated-agent.example.com/api"),
					check.That(resourceType+".test_update").Key("preferred_transport").HasValue("HTTP+JSON"),
					check.That(resourceType+".test_update").Key("additional_interfaces.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.display_name").HasValue("Update Test Agent Card Maximal"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("2.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.push_notifications").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
		},
	})
}

// TestUnitResourceAgentInstance_04_UpdateMaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitResourceAgentInstance_04_UpdateMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentInstanceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentInstanceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with maximal configuration
			{
				Config: testConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Update Test Agent Maximal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("url").HasValue("https://updated-agent.example.com/api"),
					check.That(resourceType+".test_update").Key("additional_interfaces.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
			// Step 2: Update to minimal configuration
			{
				Config: testConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Update Test Agent Minimal"),
					check.That(resourceType+".test_update").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.display_name").HasValue("Update Test Agent Card Minimal"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("1.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("false"),
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

func testConfigUpdateMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_update_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testConfigUpdateMaximal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_update_maximal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

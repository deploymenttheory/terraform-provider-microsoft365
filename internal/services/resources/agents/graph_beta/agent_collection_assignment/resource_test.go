package graphBetaAgentsAgentCollectionAssignment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentCollectionAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_collection_assignment"
	agentCollectionAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_collection_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaAgentCollectionAssignment.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *agentCollectionAssignmentMocks.AgentCollectionAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentCollectionAssignmentMock := &agentCollectionAssignmentMocks.AgentCollectionAssignmentMock{}
	agentCollectionAssignmentMock.RegisterMocks()
	return mockClient, agentCollectionAssignmentMock
}

// TestUnitResourceAgentCollectionAssignment_01_Minimal tests creating an agent collection assignment with minimal configuration
func TestUnitResourceAgentCollectionAssignment_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentCollectionAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentCollectionAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("agent_instance_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("agent_collection_id").HasValue("22222222-2222-2222-2222-222222222222"),
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

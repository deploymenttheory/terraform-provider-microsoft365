package graphBetaAgentsAgentCollection_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentCollection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_collection"
	agentCollectionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_collection/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaAgentCollection.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *agentCollectionMocks.AgentCollectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentCollectionMock := &agentCollectionMocks.AgentCollectionMock{}
	agentCollectionMock.RegisterMocks()
	return mockClient, agentCollectionMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentCollectionMocks.AgentCollectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentCollectionMock := &agentCollectionMocks.AgentCollectionMock{}
	agentCollectionMock.RegisterErrorMocks()
	return mockClient, agentCollectionMock
}

// TestUnitResourceAgentCollection_01_Minimal tests creating an agent collection with minimal configuration
func TestUnitResourceAgentCollection_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentCollectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentCollectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("Unit Test Agent Collection Minimal"),
					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_minimal").Key("created_by").Exists(),
					check.That(resourceType+".test_minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".test_minimal").Key("last_modified_date_time").Exists(),
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

// TestUnitResourceAgentCollection_02_Maximal tests creating an agent collection with maximal configuration
func TestUnitResourceAgentCollection_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentCollectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentCollectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").HasValue("Unit Test Agent Collection Maximal"),
					check.That(resourceType+".test_maximal").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("A comprehensive test agent collection with all fields configured"),
					check.That(resourceType+".test_maximal").Key("managed_by").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test_maximal").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_maximal").Key("created_by").Exists(),
					check.That(resourceType+".test_maximal").Key("created_date_time").Exists(),
					check.That(resourceType+".test_maximal").Key("last_modified_date_time").Exists(),
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

// TestUnitResourceAgentCollection_03_UpdateMinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitResourceAgentCollection_03_UpdateMinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentCollectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentCollectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Unit Test Agent Collection Update Minimal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				Config: testConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Unit Test Agent Collection Update Maximal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("description").HasValue("Updated agent collection with all fields configured"),
					check.That(resourceType+".test_update").Key("managed_by").HasValue("33333333-3333-3333-3333-333333333333"),
				),
			},
		},
	})
}

// TestUnitResourceAgentCollection_04_UpdateMaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitResourceAgentCollection_04_UpdateMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentCollectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentCollectionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Unit Test Agent Collection Update Maximal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("description").HasValue("Updated agent collection with all fields configured"),
					check.That(resourceType+".test_update").Key("managed_by").HasValue("33333333-3333-3333-3333-333333333333"),
				),
			},
			{
				Config: testConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("display_name").HasValue("Unit Test Agent Collection Update Minimal"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
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

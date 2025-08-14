package graphBetaGroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	return mockClient, groupMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterErrorMocks()

	return mockClient, groupMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// TestGroupResource_Schema validates the resource schema
func TestGroupResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "display_name", "Test Minimal Group - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_nickname", "testminimalgroup"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_enabled", "false"),
				),
			},
		},
	})
}

// TestGroupResource_Minimal tests basic CRUD operations
func TestGroupResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "display_name", "Test Minimal Group - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_groups_group.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestGroupResource_Maximal tests maximal configuration
func TestGroupResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "display_name", "Test Maximal Group - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "description", "Maximal group for testing with all features"),
				),
			},
		},
	})
}

// TestGroupResource_Update tests update operations
func TestGroupResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
				),
			},
			// Update to maximal
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.maximal"),
				),
			},
		},
	})
}

// TestGroupResource_ErrorHandling tests error scenarios
func TestGroupResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Validation error: Invalid display name"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}
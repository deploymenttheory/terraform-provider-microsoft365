package graphBetaGroupPolicyConfigurations_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	groupPolicyConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_configurations/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupPolicyConfigurationMocks.GroupPolicyConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	groupPolicyConfigurationMock := &groupPolicyConfigurationMocks.GroupPolicyConfigurationMock{}
	groupPolicyConfigurationMock.RegisterMocks()

	return mockClient, groupPolicyConfigurationMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *groupPolicyConfigurationMocks.GroupPolicyConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	groupPolicyConfigurationMock := &groupPolicyConfigurationMocks.GroupPolicyConfigurationMock{}
	groupPolicyConfigurationMock.RegisterErrorMocks()

	return mockClient, groupPolicyConfigurationMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigSimple returns the simple configuration for testing
func testConfigSimple() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_simple.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestGroupPolicyConfigurationResource_Schema validates the resource schema
func TestGroupPolicyConfigurationResource_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.minimal", "display_name", "Test Minimal Group Policy Configuration"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.minimal", "last_modified_date_time"),
				),
			},
		},
	})
}

// TestGroupPolicyConfigurationResource_CRUD tests the full CRUD lifecycle
func TestGroupPolicyConfigurationResource_CRUD(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_group_policy_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.minimal", "display_name", "Test Minimal Group Policy Configuration"),
				),
			},
			// Update to maximal configuration
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_group_policy_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "display_name", "Test Maximal Group Policy Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "description", "Comprehensive test description for group policy configuration with all fields"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "role_scope_tag_ids.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "definition_values.#", "2"),
				),
			},
		},
	})
}

// TestGroupPolicyConfigurationResource_ImportState tests the import functionality
func TestGroupPolicyConfigurationResource_ImportState(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_group_policy_configuration.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestGroupPolicyConfigurationResource_ErrorHandling tests error scenarios
func TestGroupPolicyConfigurationResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Bad Request|BadRequest"),
			},
		},
	})
}

// TestGroupPolicyConfigurationResource_DefinitionValues tests definition values functionality
func TestGroupPolicyConfigurationResource_DefinitionValues(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_group_policy_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "definition_values.#", "2"),
					// Check first definition value
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "definition_values.*", map[string]string{
						"enabled":       "true",
						"definition_id": "157dca4c-91f0-4857-b9de-8db11c0944ee",
					}),
					// Check second definition value
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_group_policy_configuration.maximal", "definition_values.*", map[string]string{
						"enabled":       "true",
						"definition_id": "98d69f26-2201-4aed-8927-d20c29b24ed5",
					}),
				),
			},
		},
	})
}

// TestGroupPolicyConfigurationResource_Simple tests simple definition values without presentation values
func TestGroupPolicyConfigurationResource_Simple(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSimple(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_group_policy_configuration.simple"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.simple", "display_name", "Test Simple Group Policy Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.simple", "description", "Simple test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.simple", "definition_values.#", "2"),
					// Check definition values
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_group_policy_configuration.simple", "definition_values.*", map[string]string{
						"enabled":       "true",
						"definition_id": "157dca4c-91f0-4857-b9de-8db11c0944ee",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_group_policy_configuration.simple", "definition_values.*", map[string]string{
						"enabled":       "true",
						"definition_id": "98d69f26-2201-4aed-8927-d20c29b24ed5",
					}),
				),
			},
		},
	})
}

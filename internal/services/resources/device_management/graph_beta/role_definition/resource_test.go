package graphBetaRoleDefinition_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleDefinitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_definition/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *roleDefinitionMocks.RoleDefinitionMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	roleDefinitionMock := &roleDefinitionMocks.RoleDefinitionMock{}
	roleDefinitionMock.RegisterMocks()

	return mockClient, roleDefinitionMock
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

// testConfigMaximal returns the maximal custom configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestRoleDefinitionResource_Schema validates the resource schema
func TestRoleDefinitionResource_Schema(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "display_name", "Test Minimal Role Definition - Unique"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in_role_definition"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "role_permissions.#", "1"),
				),
			},
		},
	})
}

// TestRoleDefinitionResource_Minimal tests basic CRUD operations
func TestRoleDefinitionResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_definition.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "display_name", "Test Minimal Role Definition - Unique"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in_role_definition"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_definition.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_definition.maximal_custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "display_name", "Test Maximal Custom Role Definition - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "description", "Comprehensive custom role definition for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestRoleDefinitionResource_UpdateInPlace tests in-place updates
func TestRoleDefinitionResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_definition.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "display_name", "Test Minimal Role Definition - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_definition.maximal_custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "display_name", "Test Maximal Custom Role Definition - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "description", "Comprehensive custom role definition for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_custom", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestRoleDefinitionResource_RequiredFields tests required field validation
func TestRoleDefinitionResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "invalid_prefix_validation",
			config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name = "Test Role Definition"
  description  = "Test description"  
  role_permissions = [
    {
      allowed_resource_actions = [
        "InvalidPrefix_Permission"
      ]
    }
  ]
}
`,
			expectedError: `must start with 'Microsoft.Intune_'`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

// TestRoleDefinitionResource_ErrorHandling tests error scenarios
func TestRoleDefinitionResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name = "Test Role Definition for Error Handling"
  description  = "Test description"
  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_Invalid_Permission_Name"
      ]
    }
  ]
}
`,
				ExpectError: regexp.MustCompile(`invalid resource operation ID`),
			},
		},
	})
}

// TestRoleDefinitionResource_RoleScopeTagIds tests role scope tag IDs handling
func TestRoleDefinitionResource_RoleScopeTagIds(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name       = "Test Role Definition with Role Scope Tags"
  description        = "Test description"
  role_scope_tag_ids = ["0", "1", "2"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

// TestRoleDefinitionResource_RolePermissions tests role permissions handling
func TestRoleDefinitionResource_RolePermissions(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name = "Test Role Definition with Permissions"
  description  = "Test description"

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_ManagedDevices_Read",
        "Microsoft.Intune_ManagedDevices_Update",
        "Microsoft.Intune_DeviceConfigurations_Read"
      ]
    }
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "display_name", "Test Role Definition with Permissions"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_permissions.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_permissions.0.allowed_resource_actions.#", "3"),
				),
			},
		},
	})
}

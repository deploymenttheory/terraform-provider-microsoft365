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

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

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

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *roleDefinitionMocks.RoleDefinitionMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	roleDefinitionMock := &roleDefinitionMocks.RoleDefinitionMock{}
	roleDefinitionMock.RegisterErrorMocks()

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

// testConfigMaximalBuiltIn returns the maximal built-in configuration for testing
func testConfigMaximalBuiltIn() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal_builtin.tf"))
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in_role_definition", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in", "false"),

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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in_role_definition", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.minimal", "is_built_in", "false"),
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
			name: "missing is_built_in_role_definition",
			config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name = "Test Role Definition"
  description  = "Test description"
  is_built_in  = false
}
`,
			expectedError: `The argument "is_built_in_role_definition" is required`,
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
	_, roleDefinitionMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name                = "Test Role Definition"
  description                 = "Test description"
  is_built_in_role_definition = true
  is_built_in                 = true
}
`,
				ExpectError: regexp.MustCompile(`Invalid role definition data|BadRequest`),
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
  display_name                = "Test Role Definition with Role Scope Tags"
  description                 = "Test description"
  is_built_in_role_definition = false
  is_built_in                 = false
  role_scope_tag_ids          = ["0", "1", "2"]
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

// TestRoleDefinitionResource_BuiltInRole tests built-in role definition handling
func TestRoleDefinitionResource_BuiltInRole(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximalBuiltIn(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_definition.maximal_builtin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "display_name", "Test Maximal Built-in Role Definition - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "description", "Comprehensive built-in role definition for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "is_built_in", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "is_built_in_role_definition", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "built_in_role_name", "Endpoint Security Manager"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.maximal_builtin", "role_permissions.#", "1"),
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
  display_name                = "Test Role Definition with Permissions"
  description                 = "Test description"
  is_built_in_role_definition = false
  is_built_in                 = false

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.management/managedDevices/read",
        "microsoft.management/managedDevices/write",
        "microsoft.management/deviceConfigurations/read"
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
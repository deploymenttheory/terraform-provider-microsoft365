package graphBetaRoleDefinitionAssignment_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *roleAssignmentMocks.RoleAssignmentMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	roleAssignmentMock := &roleAssignmentMocks.RoleAssignmentMock{}
	roleAssignmentMock.RegisterMocks()

	return mockClient, roleAssignmentMock
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

// TestRoleAssignmentResource_Schema validates the resource schema
func TestRoleAssignmentResource_Schema(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "display_name", "Test Minimal Role Assignment - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "description", "Minimal role assignment for testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "role_definition_id", "0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "members.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "scope_configuration.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "scope_configuration.0.type", "AllLicensedUsers"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_Minimal tests basic CRUD operations
func TestRoleAssignmentResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "display_name", "Test Minimal Role Assignment - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "description", "Minimal role assignment for testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "role_definition_id", "0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "scope_configuration.0.type", "AllLicensedUsers"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_assignment.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["microsoft365_graph_beta_device_management_role_assignment.minimal"]
					if !ok {
						return "", fmt.Errorf("not found: microsoft365_graph_beta_device_management_role_assignment.minimal")
					}
					id := rs.Primary.ID
					roleDefId := rs.Primary.Attributes["role_definition_id"]
					compositeId := fmt.Sprintf("%s/%s", id, roleDefId)
					fmt.Printf("DEBUG: ImportStateIdFunc - id: %s, roleDefId: %s, compositeId: %s\n", id, roleDefId, compositeId)
					return compositeId, nil
				},
			},
		},
	})
}

// TestRoleAssignmentResource_Maximal tests maximal configuration
func TestRoleAssignmentResource_Maximal(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_assignment.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "display_name", "Test Maximal Role Assignment - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "description", "Comprehensive role assignment for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "role_definition_id", "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "members.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "scope_configuration.0.type", "ResourceScopes"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "scope_configuration.0.resource_scopes.#", "3"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_UpdateInPlace tests in-place updates
func TestRoleAssignmentResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "display_name", "Test Minimal Role Assignment - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.minimal", "scope_configuration.0.type", "AllLicensedUsers"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_assignment.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "display_name", "Test Maximal Role Assignment - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.maximal", "scope_configuration.0.type", "ResourceScopes"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_ScopeValidation tests scope configuration validation
func TestRoleAssignmentResource_ScopeValidation(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing_scope_configuration",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
}
`,
			expectedError: `Missing Scope Configuration`,
		},
		{
			name: "resource_scopes_without_scopes",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  
  scope_configuration {
    type = "ResourceScopes"
  }
}
`,
			expectedError: `Missing Resource Scopes`,
		},
		{
			name: "all_users_with_scopes",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  
  scope_configuration {
    type = "AllLicensedUsers"
    resource_scopes = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  }
}
`,
			expectedError: `Invalid Resource Scopes`,
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

// TestRoleAssignmentResource_AllDevicesScope tests AllDevices scope type
func TestRoleAssignmentResource_AllDevicesScope(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test All Devices Role Assignment"
  description        = "Role assignment for all devices scope"
  role_definition_id = "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"
  
  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
    "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
  ]
  
  scope_configuration {
    type = "AllDevices"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "display_name", "Test All Devices Role Assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "scope_configuration.0.type", "AllDevices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.#", "2"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_Members tests members handling
func TestRoleAssignmentResource_Members(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Members Role Assignment"
  description        = "Minimal role assignment for testing"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  
  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
    "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2",
    "35d09841-af73-43e6-a59f-024fef1b6b95"
  ]
  
  scope_configuration {
    type = "AllLicensedUsers"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.*", "ea8e2fb8-e909-44e6-bae7-56757cf6f347"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.*", "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.*", "35d09841-af73-43e6-a59f-024fef1b6b95"),
				),
			},
		},
	})
}

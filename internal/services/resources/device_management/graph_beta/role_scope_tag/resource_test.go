package graphBetaRoleScopeTag_test

import (
	"os"
	"path/filepath"
	"regexp"

	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleScopeTagMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_scope_tag/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *roleScopeTagMocks.RoleScopeTagMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	roleScopeTagMock := &roleScopeTagMocks.RoleScopeTagMock{}
	roleScopeTagMock.RegisterMocks()

	return mockClient, roleScopeTagMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *roleScopeTagMocks.RoleScopeTagMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	roleScopeTagMock := &roleScopeTagMocks.RoleScopeTagMock{}
	roleScopeTagMock.RegisterErrorMocks()

	return mockClient, roleScopeTagMock
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

// TestRoleScopeTagResource_Schema validates the resource schema
func TestRoleScopeTagResource_Schema(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "display_name", "Test Minimal Role Scope Tag - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "is_built_in", "false"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_Minimal tests basic CRUD operations
func TestRoleScopeTagResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_scope_tag.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "display_name", "Test Minimal Role Scope Tag - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "is_built_in", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_scope_tag.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_scope_tag.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "display_name", "Test Maximal Role Scope Tag - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "description", "Maximal role scope tag for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "is_built_in", "false"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_UpdateInPlace tests in-place updates
func TestRoleScopeTagResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_scope_tag.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "display_name", "Test Minimal Role Scope Tag - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_scope_tag.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "display_name", "Test Maximal Role Scope Tag - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "description", "Maximal role scope tag for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "assignments.#", "2"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_RequiredFields tests required field validation
func TestRoleScopeTagResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  description = "Test Role Scope Tag"
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestRoleScopeTagResource_DisplayNameUniqueness tests display name uniqueness validation
func TestRoleScopeTagResource_DisplayNameUniqueness(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create first role scope tag
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "first" {
  display_name = "Unique Test Role Scope Tag"
}
`,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_role_scope_tag.first"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.first", "display_name", "Unique Test Role Scope Tag"),
				),
			},
			// Try to create second role scope tag with same display name
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "first" {
  display_name = "Unique Test Role Scope Tag"
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "second" {
  display_name = "Unique Test Role Scope Tag"
}
`,
				ExpectError: regexp.MustCompile(`role scope tag with display name 'Unique Test Role Scope Tag' already exists|Display names must be unique`),
			},
		},
	})
}

// TestRoleScopeTagResource_ErrorHandling tests error scenarios
func TestRoleScopeTagResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag"
  description  = "Test role scope tag for error testing"
}
`,
				ExpectError: regexp.MustCompile(`failed to retrieve existing role scope tags for validation|Invalid role scope tag data|BadRequest|Internal server error`),
			},
		},
	})
}

// TestRoleScopeTagResource_Assignments tests assignments handling
func TestRoleScopeTagResource_Assignments(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag with Assignments"
  description  = "Role scope tag with group assignments"
  
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    }
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_EmptyDescription tests handling of empty description
func TestRoleScopeTagResource_EmptyDescription(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag Empty Description"
  description  = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Role Scope Tag Empty Description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", ""),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_NoDescription tests handling when description is not provided
func TestRoleScopeTagResource_NoDescription(t *testing.T) {
	setupUnitTestEnvironment()
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag No Description"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Role Scope Tag No Description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "is_built_in", "false"),
				),
			},
		},
	})
}

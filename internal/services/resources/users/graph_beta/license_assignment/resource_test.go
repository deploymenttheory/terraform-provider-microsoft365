package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Helper functions to return the test configurations by reading from files
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMinimalToMaximal() string {
	// For minimal to maximal test, we need to use the maximal config
	// but with the minimal resource name and user_id to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	// Replace the user_id to match the minimal one
	updatedMaximal = strings.Replace(updatedMaximal, "00000000-0000-0000-0000-000000000003", "00000000-0000-0000-0000-000000000002", 1)

	return updatedMaximal
}

func testConfigError() string {
	// Create an error configuration with invalid user ID
	return `
resource "microsoft365_graph_beta_users_user_license_assignment" "error" {
  user_id = "invalid-user-id"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}
`
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.UserLicenseAssignmentMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	licenseAssignmentMock := &localMocks.UserLicenseAssignmentMock{}
	licenseAssignmentMock.RegisterMocks()

	return mockClient, licenseAssignmentMock
}

// Helper function to check if a resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// Helper function to get maximal config with a custom resource name
func testConfigMaximalWithResourceName(resourceName string) string {
	// Read the maximal config
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name
	updated := strings.Replace(string(content), "maximal", resourceName, 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_users_user_license_assignment" "%s" {
  user_id = "00000000-0000-0000-0000-000000000003"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}`, resourceName)
}

// TestUnitUserLicenseAssignmentResource_Create_Minimal tests the creation of a license assignment with minimal configuration
func TestUnitUserLicenseAssignmentResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.disabled_plans.#", "0"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Create_Maximal tests the creation of a license assignment with maximal configuration
func TestUnitUserLicenseAssignmentResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "user_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.sku_id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.disabled_plans.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.1.sku_id", "77777777-7777-7777-7777-777777777777"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "remove_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "remove_licenses.0", "88888888-8888-8888-8888-888888888888"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitUserLicenseAssignmentResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.disabled_plans.#", "0"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.disabled_plans.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.1.sku_id", "77777777-7777-7777-7777-777777777777"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "remove_licenses.#", "1"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitUserLicenseAssignmentResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "remove_licenses.#", "1"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),
					// Don't check for absence of attributes as they may appear as computed
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Delete_Minimal tests deleting a license assignment with minimal configuration
func TestUnitUserLicenseAssignmentResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_users_user_license_assignment.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Delete_Maximal tests deleting a license assignment with maximal configuration
func TestUnitUserLicenseAssignmentResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_users_user_license_assignment.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Import tests importing a resource
func TestUnitUserLicenseAssignmentResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_users_user_license_assignment.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"remove_licenses", // This is expected to be ignored on import
					"add_licenses",    // The import only sets assigned_licenses, not add_licenses
				},
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Error tests error handling
func TestUnitUserLicenseAssignmentResource_Error(t *testing.T) {
	// Set up mock environment
	_, licenseAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	licenseAssignmentMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Attribute user_id Must be a valid UUID format"),
			},
		},
	})
}

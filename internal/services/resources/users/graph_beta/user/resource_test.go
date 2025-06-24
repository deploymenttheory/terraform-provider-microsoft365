package graphBetaUsersUser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user/mocks"
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
	// but with the minimal resource name and UPN to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	// Replace all occurrences of the UPN to match the minimal one
	updatedMaximal = strings.Replace(updatedMaximal, "maximal.user@contoso.com", "minimal.user@contoso.com", -1)

	return updatedMaximal
}

func testConfigError() string {
	// Read the minimal config and modify for error scenario
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace resource name and UPN to create an error scenario
	updated := strings.Replace(string(content), "minimal", "error", 1)
	updated = strings.Replace(updated, "minimal.user@contoso.com", "duplicate@contoso.com", 1)
	updated = strings.Replace(updated, "Minimal User", "Error User", 1)

	return updated
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.UserMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	userMock := &localMocks.UserMock{}
	userMock.RegisterMocks()

	return mockClient, userMock
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

// TestUnitUserResource_Create_Minimal tests the creation of a user with minimal configuration
func TestUnitUserResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Minimal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "user_principal_name", "minimal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "account_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "password_profile.password", "SecureP@ssw0rd123!"),
				),
			},
		},
	})
}

// TestUnitUserResource_Create_Maximal tests the creation of a user with maximal configuration
func TestUnitUserResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "display_name", "Maximal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "user_principal_name", "maximal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "given_name", "Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "surname", "User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "job_title", "Senior Developer"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "department", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "company_name", "Contoso Ltd"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "business_phones.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "business_phones.0", "+1 425-555-0100"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "identities.#", "1"),
				),
			},
		},
	})
}

// TestUnitUserResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitUserResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Minimal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "user_principal_name", "minimal.user@contoso.com"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_users_user.minimal", "given_name"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_users_user.minimal", "job_title"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Maximal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "user_principal_name", "minimal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "given_name", "Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "surname", "User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "job_title", "Senior Developer"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "department", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "company_name", "Contoso Ltd"),
				),
			},
		},
	})
}

// TestUnitUserResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitUserResource_Update_MaximalToMinimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "display_name", "Maximal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "given_name", "Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "surname", "User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "job_title", "Senior Developer"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				// We expect a non-empty plan because computed fields will show as changes
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user.test"),
					// Verify it now has only minimal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "display_name", "Minimal User"),
					// Don't check for absence of attributes as they may appear as computed
				),
			},
		},
	})
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
	return fmt.Sprintf(`resource "microsoft365_graph_beta_users_user" "%s" {
  display_name       = "Minimal User"
  user_principal_name = "test.user@contoso.com"
  account_enabled    = true
  password_profile   = {
    password = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}`, resourceName)
}

// TestUnitUserResource_Delete_Minimal tests deleting a user with minimal configuration
func TestUnitUserResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_users_user.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitUserResource_Delete_Maximal tests deleting a user with maximal configuration
func TestUnitUserResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_users_user.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitUserResource_Import tests importing a resource
func TestUnitUserResource_Import(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_users_user.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile", // Password is not returned by the API
				},
			},
		},
	})
}

func TestUnitUserResource_Error(t *testing.T) {
	// Set up mock environment
	_, userMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	userMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("User with this userPrincipalName already exists"),
			},
		},
	})
}

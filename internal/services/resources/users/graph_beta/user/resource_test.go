package graphBetaUsersUser_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Minimal configuration with only required attributes
	testConfigMinimalTemplate = `
resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "minimal.user@contoso.com"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}
`

	// Maximal configuration with all possible attributes
	testConfigMaximalTemplate = `
resource "microsoft365_graph_beta_users_user" "maximal" {
  display_name        = "Maximal User"
  user_principal_name = "maximal.user@contoso.com"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "User"
  mail                = "maximal.user@contoso.com"
  mail_nickname       = "maxuser"
  job_title           = "Senior Developer"
  department          = "Engineering"
  company_name        = "Contoso Ltd"
  office_location     = "Building A"
  city                = "Redmond"
  state               = "WA"
  country             = "US"
  postal_code         = "98052"
  usage_location      = "US"
  business_phones     = ["+1 425-555-0100"]
  mobile_phone        = "+1 425-555-0101"
  password_profile = {
    password                          = "SecureP@ssw0rd!"
    force_change_password_next_sign_in = true
  }
  identities = [
    {
      sign_in_type       = "emailAddress"
      issuer             = "contoso.com"
      issuer_assigned_id = "maximal.user@contoso.com"
    }
  ]
  other_mails     = ["maximal.user.other@contoso.com"]
  proxy_addresses = ["SMTP:maximal.user@contoso.com"]
}
`

	// Update configuration for testing changes
	testConfigUpdateTemplate = `
resource "microsoft365_graph_beta_users_user" "test" {
  display_name        = "Updated User"
  user_principal_name = "minimal.user@contoso.com"
  job_title           = "Updated Job Title"
  department          = "Updated Department"
  password_profile = {
    password = "NewSecureP@ssw0rd!"
  }
}
`

	// Configuration with multiple resources for testing lifecycle
	testConfigMultipleResourcesTemplate = `
resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "minimal.user@contoso.com"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}

resource "microsoft365_graph_beta_users_user" "maximal" {
  display_name        = "Maximal User"
  user_principal_name = "maximal.user@contoso.com"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "User"
  mail                = "maximal.user@contoso.com"
  mail_nickname       = "maxuser"
  job_title           = "Senior Developer"
  department          = "Engineering"
  company_name        = "Contoso Ltd"
  office_location     = "Building A"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}
`

	// Configuration for updating minimal to maximal
	testConfigMinimalToMaximalTemplate = `
resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "minimal.user@contoso.com"
  given_name          = "Updated"
  surname             = "User"
  job_title           = "Developer"
  department          = "IT"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}

resource "microsoft365_graph_beta_users_user" "maximal" {
  display_name        = "Maximal User"
  user_principal_name = "maximal.user@contoso.com"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}
`

	// Error configuration with duplicate UPN
	testConfigErrorTemplate = `
resource "microsoft365_graph_beta_users_user" "error" {
  display_name        = "Error User"
  user_principal_name = "duplicate@contoso.com"
  password_profile = {
    password = "SecureP@ssw0rd!"
  }
}
`
)

// Helper functions to return the test configurations
func testConfigMinimal() string {
	return testConfigMinimalTemplate
}

func testConfigMaximal() string {
	return testConfigMaximalTemplate
}

func testConfigUpdate() string {
	return testConfigUpdateTemplate
}

func testConfigMultipleResources() string {
	return testConfigMultipleResourcesTemplate
}

func testConfigMinimalToMaximal() string {
	return testConfigMinimalToMaximalTemplate
}

func testConfigError() string {
	return testConfigErrorTemplate
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

func TestUnitUserResource_Minimal(t *testing.T) {
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "password_profile.password", "SecureP@ssw0rd!"),
				),
			},
		},
	})
}

func TestUnitUserResource_Maximal(t *testing.T) {
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

// TestUnitUserResource_Create tests the creation of multiple resources
func TestUnitUserResource_Create(t *testing.T) {
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
				Config: testConfigMultipleResources(),
				Check: resource.ComposeTestCheckFunc(
					// Check minimal resource
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Minimal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "user_principal_name", "minimal.user@contoso.com"),

					// Check maximal resource
					testCheckExists("microsoft365_graph_beta_users_user.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "display_name", "Maximal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "user_principal_name", "maximal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "given_name", "Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "surname", "User"),
				),
			},
		},
	})
}

// TestUnitUserResource_Update tests updating a resource
func TestUnitUserResource_Update(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Minimal User"),
				),
			},
			// Update with new values
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "display_name", "Updated User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "job_title", "Updated Job Title"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.test", "department", "Updated Department"),
				),
			},
		},
	})
}

// TestUnitUserResource_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitUserResource_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create both resources
			{
				Config: testConfigMultipleResources(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					testCheckExists("microsoft365_graph_beta_users_user.maximal"),
				),
			},
			// Update - transform minimal to maximal and maximal to minimal
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Check former minimal now with more attributes
					testCheckExists("microsoft365_graph_beta_users_user.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "display_name", "Minimal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "given_name", "Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "surname", "User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "job_title", "Developer"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.minimal", "department", "IT"),

					// Check former maximal now with fewer attributes
					testCheckExists("microsoft365_graph_beta_users_user.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "display_name", "Maximal User"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user.maximal", "user_principal_name", "maximal.user@contoso.com"),
				),
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

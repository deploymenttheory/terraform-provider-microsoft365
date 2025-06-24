package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Basic configuration with standard attributes
	testConfigBasicTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "test" {
  user_id = "00000000-0000-0000-0000-000000000001"
  add_licenses = [{
    sku_id = "11111111-1111-1111-1111-111111111111"
    disabled_plans = ["22222222-2222-2222-2222-222222222222"]
  }]
}
`

	// Minimal configuration with only required attributes
	testConfigMinimalTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "00000000-0000-0000-0000-000000000002"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}
`

	// Maximal configuration with all possible attributes
	testConfigMaximalTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = "00000000-0000-0000-0000-000000000003"
  add_licenses = [
    {
      sku_id = "44444444-4444-4444-4444-444444444444"
      disabled_plans = [
        "55555555-5555-5555-5555-555555555555",
        "66666666-6666-6666-6666-666666666666"
      ]
    },
    {
      sku_id = "77777777-7777-7777-7777-777777777777"
    }
  ]
}
`

	// Update configuration for testing changes
	testConfigUpdateTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "test" {
  user_id = "00000000-0000-0000-0000-000000000001"
  add_licenses = [
    {
      sku_id = "11111111-1111-1111-1111-111111111111"
      disabled_plans = []
    },
    {
      sku_id = "88888888-8888-8888-8888-888888888888"
    }
  ]
}
`

	// Configuration with remove_licenses
	testConfigRemoveLicensesTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "test" {
  user_id = "00000000-0000-0000-0000-000000000001"
  add_licenses = [{
    sku_id = "11111111-1111-1111-1111-111111111111"
    disabled_plans = ["22222222-2222-2222-2222-222222222222"]
  }]
  remove_licenses = ["88888888-8888-8888-8888-888888888888"]
}
`

	// Configuration with multiple resources for testing lifecycle
	testConfigMultipleResourcesTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "00000000-0000-0000-0000-000000000002"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}

resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = "00000000-0000-0000-0000-000000000003"
  add_licenses = [
    {
      sku_id = "44444444-4444-4444-4444-444444444444"
      disabled_plans = [
        "55555555-5555-5555-5555-555555555555",
        "66666666-6666-6666-6666-666666666666"
      ]
    },
    {
      sku_id = "77777777-7777-7777-7777-777777777777"
    }
  ]
}
`

	// Configuration for updating minimal to maximal
	testConfigMinimalToMaximalTemplate = `
resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "00000000-0000-0000-0000-000000000002"
  add_licenses = [
    {
      sku_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      sku_id = "44444444-4444-4444-4444-444444444444"
      disabled_plans = [
        "55555555-5555-5555-5555-555555555555"
      ]
    }
  ]
}

resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = "00000000-0000-0000-0000-000000000003"
  add_licenses = [{
    sku_id = "44444444-4444-4444-4444-444444444444"
  }]
}
`
)

// Helper functions to return the test configurations
func testConfigBasic() string {
	return testConfigBasicTemplate
}

func testConfigMinimal() string {
	return testConfigMinimalTemplate
}

func testConfigMaximal() string {
	return testConfigMaximalTemplate
}

func testConfigUpdate() string {
	return testConfigUpdateTemplate
}

func testConfigRemoveLicenses() string {
	return testConfigRemoveLicensesTemplate
}

func testConfigMultipleResources() string {
	return testConfigMultipleResourcesTemplate
}

func testConfigMinimalToMaximal() string {
	return testConfigMinimalToMaximalTemplate
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

func TestUnitUserLicenseAssignmentResource_Minimal(t *testing.T) {
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_principal_name", "minimal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.disabled_plans.#", "0"),
				),
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Maximal(t *testing.T) {
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "user_principal_name", "maximal.user@contoso.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.sku_id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.disabled_plans.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.1.sku_id", "77777777-7777-7777-7777-777777777777"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.1.disabled_plans.#", "0"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Create tests the creation of multiple resources
func TestUnitUserLicenseAssignmentResource_Create(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),

					// Check maximal resource
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "user_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.sku_id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.1.sku_id", "77777777-7777-7777-7777-777777777777"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_Update tests updating a resource
func TestUnitUserLicenseAssignmentResource_Update(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with basic configuration
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "user_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.0.sku_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.0.disabled_plans.#", "1"),
				),
			},
			// Update to add another license and remove disabled plans
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "user_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.0.sku_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.0.disabled_plans.#", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.1.sku_id", "88888888-8888-8888-8888-888888888888"),
				),
			},
		},
	})
}

// TestUnitUserLicenseAssignmentResource_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitUserLicenseAssignmentResource_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.maximal"),
				),
			},
			// Update - transform minimal to maximal and maximal to minimal
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Check former minimal now with more licenses
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "user_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.0.sku_id", "33333333-3333-3333-3333-333333333333"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.minimal", "add_licenses.1.sku_id", "44444444-4444-4444-4444-444444444444"),

					// Check former maximal now with fewer licenses
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "user_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.maximal", "add_licenses.0.sku_id", "44444444-4444-4444-4444-444444444444"),
				),
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
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_users_user_license_assignment.test",
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

func TestUnitUserLicenseAssignmentResource_RemoveLicenses(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with update configuration first to add the license we want to remove
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "add_licenses.1.sku_id", "88888888-8888-8888-8888-888888888888"),
				),
			},
			// Then remove the license
			{
				Config: testConfigRemoveLicenses(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_users_user_license_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "remove_licenses.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_users_user_license_assignment.test", "remove_licenses.0", "88888888-8888-8888-8888-888888888888"),
				),
			},
		},
	})
}

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
				Config: `
resource "microsoft365_graph_beta_users_user_license_assignment" "error" {
  user_id = "99999999-9999-9999-9999-999999999999"
  add_licenses = [{
    sku_id = "11111111-1111-1111-1111-111111111111"
  }]
}
`,
				ExpectError: regexp.MustCompile("Error assigning license"),
			},
		},
	})
}

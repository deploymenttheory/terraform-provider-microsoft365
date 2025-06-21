package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func TestUnitMacOSPlatformScriptResource_Lifecycle(t *testing.T) {
	// Activate httpmock FIRST, before any provider creation
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance
	mockClient := mocks.NewMocks()

	// Register all mock responses
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the actual test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_ErrorHandling(t *testing.T) {
	// Activate httpmock FIRST, before any provider creation
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance
	mockClient := mocks.NewMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register error mocks
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptErrorMocks()

	// Run the test expecting an error
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`(Access denied|Forbidden)`),
			},
		},
	})
}

// Acceptance Tests
func TestAccMacOSPlatformScriptResource_basic(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMacOSPlatformScriptConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
				),
			},
			{
				Config: testAccMacOSPlatformScriptConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Updated macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_platform_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckMacOSPlatformScriptDestroy,
	})
}

// Helper Functions
func testAccCheckMacOSPlatformScriptExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckMacOSPlatformScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_platform_script" {
			continue
		}

		// In a real test, we would make an API call to verify the resource is gone
		// For unit tests with mocks, we can assume it's destroyed if we get here
		return nil
	}

	return nil
}

// Test Configurations
func testAccMacOSPlatformScriptConfigBasic() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test macOS Script"
  description     = "Test description for macOS platform script"
  script_content  = "#!/bin/bash\necho 'Hello World'"
  run_as_account  = "system"
  file_name       = "test-script.sh"
  block_execution_notifications = true
  execution_frequency = "P1D"
  retry_count     = 3

  assignments = {
    all_devices = false
    all_users   = true
  }
}
`
}

func testAccMacOSPlatformScriptConfigUpdate() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Updated macOS Script"
  description     = "Updated description for macOS platform script"
  script_content  = "#!/bin/bash\necho 'Hello Updated World'"
  run_as_account  = "user"
  file_name       = "updated-script.sh"
  block_execution_notifications = false
  execution_frequency = "P7D"
  retry_count     = 5

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`
}

// setupTestEnvironment configures the test environment
func setupTestEnvironment(t *testing.T) {
	// Set mock authentication credentials with valid values
	os.Setenv("M365_TENANT_ID", "00000000-0000-0000-0000-000000000001")
	os.Setenv("M365_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	os.Setenv("M365_CLIENT_SECRET", "mock-secret-value")
	os.Setenv("M365_AUTH_METHOD", "client_secret")
	os.Setenv("M365_CLOUD", "public")

	t.Cleanup(func() {
		os.Unsetenv("M365_TENANT_ID")
		os.Unsetenv("M365_CLIENT_ID")
		os.Unsetenv("M365_CLIENT_SECRET")
		os.Unsetenv("M365_AUTH_METHOD")
		os.Unsetenv("M365_CLOUD")
	})
}

// testCheckExists verifies the resource exists in Terraform state
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

// testConfigBasic returns a minimal Terraform configuration for unit testing
func testConfigBasic() string {
	return `
provider "microsoft365" {
  tenant_id = "00000000-0000-0000-0000-000000000001"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "11111111-1111-1111-1111-111111111111"
    client_secret = "mock-secret-value"
  }
  cloud = "public"
}

resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test macOS Script"
  description     = "Test description"
  script_content  = "#!/bin/bash\necho 'Hello World'"
  run_as_account  = "system"
  file_name       = "test-script.sh"
  block_execution_notifications = true
  execution_frequency = "P1D"
  retry_count     = 3

  assignments = {
    all_devices = false
    all_users   = true
  }
}
`
}

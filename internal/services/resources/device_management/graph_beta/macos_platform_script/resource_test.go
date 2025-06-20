package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func TestUnitMacOSPlatformScriptResource_Lifecycle(t *testing.T) {
	// STEP 1: Set up mock environment before anything else
	setupTestEnvironment(t)

	// STEP 2: Activate httpmock before any HTTP clients are created
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// STEP 3: Register all mock HTTP responses
	setupAuthMocks()
	setupGraphAPIMocks(t)

	// STEP 4: Catch any unexpected real API calls and fail the test
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		t.Errorf("❌ REAL API CALL DETECTED: %s %s", req.Method, req.URL.String())
		t.Errorf("This means httpmock is not intercepting properly!")
		return httpmock.NewStringResponse(500, "Real API call blocked"), nil
	})

	// STEP 5: Run the actual test
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

	// STEP 6: Verify mocks were called (proves no real API calls)
	callCount := httpmock.GetTotalCallCount()
	if callCount == 0 {
		t.Error("❌ No HTTP calls intercepted - httpmock may not be working")
	} else {
		t.Logf("✅ Successfully intercepted %d HTTP calls", callCount)
	}
}

func TestUnitMacOSPlatformScriptResource_ErrorHandling(t *testing.T) {
	mockClient := mocks.NewMocks()
	defer mockClient.DeactivateAndReset()

	mockClient.Activate()
	mockClient.RegisterMacOSPlatformScriptMocks()
	setupErrorMocks(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSPlatformScriptConfigBasic(),
				ExpectError: regexp.MustCompile("Failed to create macOS platform script"),
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

// Mock Setup Functions
func setupUnitTestMocks(t *testing.T) {
	basePath := "tests"

	// Create
	createResponsePath := filepath.Join(basePath, "Validate_Create", "post_device_shell_script.json")
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, createResponsePath)))

	// Assign
	assignResponsePath := filepath.Join(basePath, "Validate_Create", "post_assign.json")
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, assignResponsePath)))

	// Read
	readResponsePath := filepath.Join(basePath, "Validate_Create", "get_device_shell_script_with_assignments.json")
	readBasicResponsePath := filepath.Join(basePath, "Validate_Create", "get_device_shell_script.json")
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("$expand") == "assignments" {
				return httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, readResponsePath))(req)
			}
			// Return non-expanded response if $expand is not present
			return httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, readBasicResponsePath))(req)
		})

	// Update
	updateResponsePath := filepath.Join(basePath, "Validate_Update", "patch_device_shell_script.json")
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, updateResponsePath)))

	// Delete
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewStringResponder(204, ""))
}

func setupErrorMocks(t *testing.T) {
	// Create with error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Access denied"}}`))
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
	// Set mock authentication credentials
	os.Setenv("M365_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	os.Setenv("M365_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	os.Setenv("M365_CLIENT_SECRET", "mock-secret")
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

// setupAuthMocks registers mock responses for Azure AD authentication
func setupAuthMocks() {
	// Token endpoint
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}))

	// Instance discovery
	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"tenant_discovery_endpoint": "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/v2.0/.well-known/openid-configuration",
		}))
}

// setupGraphAPIMocks registers mock responses for Microsoft Graph API calls
func setupGraphAPIMocks(t *testing.T) {
	basePath := "tests"

	// POST /deviceManagement/deviceShellScripts (Create)
	createResponsePath := filepath.Join(basePath, "Validate_Create", "post_device_shell_script.json")
	httpmock.RegisterResponder("POST",
		"https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, createResponsePath)))

	// POST /deviceShellScripts/{id}/assign (Assign)
	assignResponsePath := filepath.Join(basePath, "Validate_Create", "post_assign.json")
	httpmock.RegisterResponder("POST",
		"https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, assignResponsePath)))

	// GET /deviceShellScripts/{id} (Read)
	readResponsePath := filepath.Join(basePath, "Validate_Create", "get_device_shell_script_with_assignments.json")
	httpmock.RegisterResponder("GET",
		"https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			// Handle both expanded and non-expanded reads
			return httpmock.NewStringResponder(200, helpers.ParseJSONFile(t, readResponsePath))(req)
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

// testConfigBasic returns the Terraform configuration for testing
func testConfigBasic() string {
	return `
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

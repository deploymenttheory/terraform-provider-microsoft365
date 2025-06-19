package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Unit Tests
func TestConstructResource(t *testing.T) {
	t.Skip("Skipping test that requires access to unexported functions")
}

func TestConstructAssignment(t *testing.T) {
	t.Skip("Skipping test that requires access to unexported functions")
}

// Acceptance Tests
func TestAccMacOSPlatformScriptResource_basic(t *testing.T) {
	mockClient := mocks.NewMocks()
	defer mockClient.DeactivateAndReset()

	mockClient.Activate()
	setupBasicMocks()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMacOSPlatformScriptConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
				),
			},
		},
	})
}

func TestAccMacOSPlatformScriptResource_update(t *testing.T) {
	mockClient := mocks.NewMocks()
	defer mockClient.DeactivateAndReset()

	mockClient.Activate()
	setupUpdateMocks()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMacOSPlatformScriptConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
				),
			},
			{
				Config: testAccMacOSPlatformScriptConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Updated macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
				),
			},
		},
	})
}

// Mock Setup Functions
func setupBasicMocks() {
	// Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":                          "00000000-0000-0000-0000-000000000001",
			"displayName":                 "Test macOS Script",
			"description":                 "Test description for macOS platform script",
			"runAsAccount":                "system",
			"fileName":                    "test-script.sh",
			"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn",
			"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
			"lastModifiedDateTime":        "2023-11-01T10:30:00.0000000Z",
			"roleScopeTagIds":             []string{"0"},
			"blockExecutionNotifications": true,
			"executionFrequency":          "P1D",
			"retryCount":                  3,
		}))

	// Assign
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"value": "Assignment completed successfully",
		}))

	// Read
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000001",
				"displayName":                 "Test macOS Script",
				"description":                 "Test description for macOS platform script",
				"runAsAccount":                "system",
				"fileName":                    "test-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn",
				"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime":        "2023-11-01T10:30:00.0000000Z",
				"roleScopeTagIds":             []string{"0"},
				"blockExecutionNotifications": true,
				"executionFrequency":          "P1D",
				"retryCount":                  3,
			}

			if req.URL.Query().Get("$expand") == "assignments" {
				response["assignments"] = map[string]interface{}{
					"value": []map[string]interface{}{
						{
							"id": "00000000-0000-0000-0000-000000000001_assignment1",
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
							},
						},
					},
				}
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Delete
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func setupUpdateMocks() {
	setupBasicMocks()

	// Update
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":                          "00000000-0000-0000-0000-000000000001",
			"displayName":                 "Updated macOS Script",
			"description":                 "Updated description for macOS platform script",
			"runAsAccount":                "user",
			"fileName":                    "updated-script.sh",
			"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gVXBkYXRlZCBXb3JsZCc=",
			"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
			"lastModifiedDateTime":        "2023-11-01T11:45:00.0000000Z",
			"roleScopeTagIds":             []string{"0"},
			"blockExecutionNotifications": false,
			"executionFrequency":          "P7D",
			"retryCount":                  5,
		}))
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

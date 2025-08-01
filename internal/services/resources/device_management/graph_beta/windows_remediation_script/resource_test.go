package graphBetaWindowsRemediationScript_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	windowsRemediationScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_remediation_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*windowsRemediationScriptMocks.WindowsRemediationScriptMock, *windowsRemediationScriptMocks.WindowsRemediationScriptMock) {
	httpmock.Activate()
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	errorMock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	return mock, errorMock
}

func setupTestEnvironment(t *testing.T) {
	// Set up any test-specific environment variables or configurations here if needed
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
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

	// Fix the display name to match test expectations
	updated = strings.Replace(updated, "Test Maximal Windows Remediation Script - Unique", "Test Maximal Windows Remediation Script", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_windows_remediation_script" "%s" {
  display_name                = "Test Minimal Windows Remediation Script"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitWindowsRemediationScriptResource_Create_Minimal tests the creation of a Windows remediation script with minimal configuration
func TestUnitWindowsRemediationScriptResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "display_name", "Test Minimal Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "publisher", "Terraform Provider Test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "run_as_32_bit", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "enforce_signature_check", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "detection_script_content"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "remediation_script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "role_scope_tag_ids.0", "0"),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Create_Maximal tests the creation of a Windows remediation script with maximal configuration
func TestUnitWindowsRemediationScriptResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "display_name", "Test Maximal Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "description", "Maximal Windows remediation script for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "publisher", "Terraform Provider Test Suite"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "run_as_32_bit", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "detection_script_content"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "remediation_script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "role_scope_tag_ids.1", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "assignments.#", "3"),
					// Daily schedule assignment
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "assignments.*", map[string]string{
						"type":                    "groupAssignmentTarget",
						"group_id":                "44444444-4444-4444-4444-444444444444",
						"filter_id":               "55555555-5555-5555-5555-555555555555",
						"filter_type":             "include",
						"daily_schedule.interval": "1",
						"daily_schedule.time":     "09:00:00",
						"daily_schedule.use_utc":  "true",
					}),
					// Hourly schedule assignment
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "assignments.*", map[string]string{
						"type":                     "groupAssignmentTarget",
						"group_id":                 "33333333-3333-3333-3333-333333333333",
						"filter_id":                "66666666-6666-6666-6666-666666666666",
						"filter_type":              "exclude",
						"hourly_schedule.interval": "4",
					}),
					// Run once schedule assignment
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "assignments.*", map[string]string{
						"type":                      "allDevicesAssignmentTarget",
						"filter_id":                 "00000000-0000-0000-0000-000000000000",
						"filter_type":               "none",
						"run_once_schedule.date":    "2024-12-31",
						"run_once_schedule.time":    "23:59:00",
						"run_once_schedule.use_utc": "false",
					}),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitWindowsRemediationScriptResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Minimal Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.#", "1"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Maximal Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "description", "Maximal Windows remediation script for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_32_bit", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitWindowsRemediationScriptResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Maximal Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_32_bit", "true"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Minimal Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.#", "1"),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Delete_Minimal tests deleting a Windows remediation script with minimal configuration
func TestUnitWindowsRemediationScriptResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.minimal"),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Delete_Maximal tests deleting a Windows remediation script with maximal configuration
func TestUnitWindowsRemediationScriptResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.maximal"),
				),
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Import tests importing a Windows remediation script
func TestUnitWindowsRemediationScriptResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_remediation_script.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitWindowsRemediationScriptResource_Error tests error handling
func TestUnitWindowsRemediationScriptResource_Error(t *testing.T) {
	// Set up mock environment
	_, errorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the error mocks
	errorMock.RegisterErrorMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Validation error: Invalid display name"),
			},
		},
	})
}

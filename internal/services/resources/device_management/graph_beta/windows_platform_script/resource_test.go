package graphBetaWindowsPlatformScript_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	platformScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*platformScriptMocks.WindowsPlatformScriptMock, *platformScriptMocks.WindowsPlatformScriptMock) {
	httpmock.Activate()
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	errorMock := &platformScriptMocks.WindowsPlatformScriptMock{}
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
	updated = strings.Replace(updated, "Test Maximal Windows Platform Script - Unique", "Test Maximal Windows Platform Script", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_windows_platform_script" "%s" {
  display_name    = "Test Minimal Windows Platform Script"
  file_name       = "test_minimal.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'Hello World'\nExit 0"
  run_as_account  = "system"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitWindowsPlatformScriptResource_Create_Minimal tests the creation of a platform script with minimal configuration
func TestUnitWindowsPlatformScriptResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "display_name", "Test Minimal Windows Platform Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "file_name", "test_minimal.ps1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "run_as_account", "system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "role_scope_tag_ids.0", "0"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Create_Maximal tests the creation of a platform script with maximal configuration
func TestUnitWindowsPlatformScriptResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "display_name", "Test Maximal Windows Platform Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "description", "Maximal platform script for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "file_name", "test_maximal.ps1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "run_as_account", "user"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "role_scope_tag_ids.1", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "run_as_32_bit", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.maximal", "assignments.0.group_id", "44444444-4444-4444-4444-444444444444"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitWindowsPlatformScriptResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Minimal Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "role_scope_tag_ids.#", "1"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Maximal Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_32_bit", "false"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitWindowsPlatformScriptResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Maximal Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "user"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Minimal Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "role_scope_tag_ids.#", "1"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Delete_Minimal tests deleting a platform script with minimal configuration
func TestUnitWindowsPlatformScriptResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.minimal"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Delete_Maximal tests deleting a platform script with maximal configuration
func TestUnitWindowsPlatformScriptResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.maximal"),
				),
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Import tests importing a platform script
func TestUnitWindowsPlatformScriptResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &platformScriptMocks.WindowsPlatformScriptMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_platform_script.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitWindowsPlatformScriptResource_Error tests error handling
func TestUnitWindowsPlatformScriptResource_Error(t *testing.T) {
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
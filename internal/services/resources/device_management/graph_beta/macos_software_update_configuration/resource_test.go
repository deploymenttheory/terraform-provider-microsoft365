package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	softwareUpdateConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_software_update_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock, *softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock) {
	httpmock.Activate()
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	errorMock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
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
	updated = strings.Replace(updated, "Test Maximal macOS Software Update Configuration - Unique", "Test Maximal macOS Software Update Configuration", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "%s" {
  display_name           = "Test Minimal macOS Software Update Configuration"
  update_schedule_type   = "alwaysUpdate"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Minimal tests the creation of a software update configuration with minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "display_name", "Test Minimal macOS Software Update Configuration - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "config_data_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "firmware_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "all_other_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "role_scope_tag_ids.0", "0"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Maximal tests the creation of a software update configuration with maximal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "display_name", "Test Maximal macOS Software Update Configuration - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "description", "Maximal software update configuration for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "update_time_window_utc_offset_in_minutes", "-480"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "max_user_deferrals_count", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "priority", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "role_scope_tag_ids.1", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "custom_update_time_windows.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.0.group_id", "44444444-4444-4444-4444-444444444444"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "role_scope_tag_ids.#", "1"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Maximal software update configuration for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "high"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "role_scope_tag_ids.#", "1"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Minimal tests deleting a software update configuration with minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Maximal tests deleting a software update configuration with maximal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Import tests importing a software update configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Error tests error handling
func TestUnitMacOSSoftwareUpdateConfigurationResource_Error(t *testing.T) {
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
package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_software_update_configuration/mocks"
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
	// but with the minimal resource name to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	return updatedMaximal
}

func testConfigError() string {
	// Read the minimal config and modify for error scenario
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace resource name and display name to create an error scenario
	updated := strings.Replace(string(content), "minimal", "error", 1)
	updated = strings.Replace(updated, "Minimal macOS Software Update Configuration", "Error macOS Software Update Configuration", 1)

	return updated
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.MacOSSoftwareUpdateConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	configMock := &localMocks.MacOSSoftwareUpdateConfigurationMock{}
	configMock.RegisterMocks()

	return mockClient, configMock
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

// TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Minimal tests the creation of a configuration with minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "display_name", "Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "critical_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "config_data_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "firmware_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "all_other_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "update_time_window_utc_offset_in_minutes", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "assignments.all_users", "false"),
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Maximal tests the creation of a configuration with maximal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "display_name", "Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "config_data_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "firmware_update_behavior", "downloadOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "all_other_update_behavior", "installLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "update_time_window_utc_offset_in_minutes", "60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "max_user_deferrals_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "priority", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.all_users", "false"),
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

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "display_name", "Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "critical_update_behavior", "default"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "description"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "max_user_deferrals_count"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "priority"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "display_name", "Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "config_data_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "firmware_update_behavior", "downloadOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "all_other_update_behavior", "installLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "max_user_deferrals_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "priority", "high"),
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
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "%s" {
  display_name                          = "Minimal macOS Software Update Configuration"
  critical_update_behavior              = "default"
  config_data_update_behavior           = "default"
  firmware_update_behavior              = "default"
  all_other_update_behavior             = "default"
  update_schedule_type                  = "alwaysUpdate"
  update_time_window_utc_offset_in_minutes = 0

  assignments = {
    all_devices = false
    all_users   = false
  }
}`, resourceName)
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "max_user_deferrals_count", "3"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				// We expect a non-empty plan because computed fields will show as changes
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					// Verify it now has only minimal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "default"),
					// Don't check for absence of attributes as they may appear as computed
				),
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Minimal tests deleting a configuration with minimal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Maximal tests deleting a configuration with maximal configuration
func TestUnitMacOSSoftwareUpdateConfigurationResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitMacOSSoftwareUpdateConfigurationResource_Import tests importing a resource
func TestUnitMacOSSoftwareUpdateConfigurationResource_Import(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
				),
			},
			// Import
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
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	configMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Configuration with this name already exists"),
			},
		},
	})
}

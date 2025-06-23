package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_software_update_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Basic configuration with standard attributes
	testConfigBasicTemplate = `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name                          = "Test macOS Software Update Configuration"
  description                           = "Test description"
  critical_update_behavior              = "default"
  config_data_update_behavior           = "default"
  firmware_update_behavior              = "default"
  all_other_update_behavior             = "default"
  update_schedule_type                  = "alwaysUpdate"
  update_time_window_utc_offset_in_minutes = 0
  priority                              = "low"

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`

	// Minimal configuration with only required attributes
	testConfigMinimalTemplate = `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "minimal" {
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
}
`

	// Maximal configuration with all possible attributes
	testConfigMaximalTemplate = `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "maximal" {
  display_name                          = "Maximal macOS Software Update Configuration"
  description                           = "This is a comprehensive configuration with all fields populated"
  critical_update_behavior              = "installASAP"
  config_data_update_behavior           = "notifyOnly"
  firmware_update_behavior              = "downloadOnly"
  all_other_update_behavior             = "installLater"
  update_schedule_type                  = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 60
  max_user_deferrals_count              = 3
  priority                              = "high"
  role_scope_tag_ids                    = ["0", "1"]

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`

	// Update configuration for testing changes
	testConfigUpdateTemplate = `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name                          = "Updated macOS Software Update Configuration"
  description                           = "Updated description"
  critical_update_behavior              = "notifyOnly"
  config_data_update_behavior           = "downloadOnly"
  firmware_update_behavior              = "installASAP"
  all_other_update_behavior             = "installLater"
  update_schedule_type                  = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 120
  max_user_deferrals_count              = 5
  priority                              = "high"

  assignments = {
    all_devices = false
    all_users   = true
  }
}
`

	// Group assignments configuration
	testConfigGroupAssignmentsTemplate = `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "group_assigned" {
  display_name                          = "Group Assignment Software Update Configuration"
  description                           = "Configuration with group assignments"
  critical_update_behavior              = "default"
  config_data_update_behavior           = "default"
  firmware_update_behavior              = "default"
  all_other_update_behavior             = "default"
  update_schedule_type                  = "alwaysUpdate"
  update_time_window_utc_offset_in_minutes = 0
  priority                              = "low"

  assignments = {
    all_devices = false
    all_users   = false
    include_group_ids = ["11111111-1111-1111-1111-111111111111"]
    exclude_group_ids = ["22222222-2222-2222-2222-222222222222"]
  }
}
`
)

// Unit test provider configuration
const unitTestProviderConfig = `
provider "microsoft365" {
  tenant_id = "00000000-0000-0000-0000-000000000001"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "11111111-1111-1111-1111-111111111111"
    client_secret = "mock-secret-value"
  }
  cloud = "public"
}
`

// Acceptance test provider configuration
const accTestProviderConfig = `
provider "microsoft365" {
  # Configuration from environment variables
}
`

func TestUnitMacOSSoftwareUpdateConfigurationResource_Basic(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "config_data_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "firmware_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "all_other_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "low"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_users", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_Minimal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "display_name", "Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal", "assignments.all_users", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_Maximal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal", "assignments.all_users", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_GroupAssignments(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "display_name", "Group Assignment Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.all_users", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.include_group_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.include_group_ids.0", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.exclude_group_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.exclude_group_ids.0", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_FullLifecycle(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
				),
			},
			// Update
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Updated macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_time_window_utc_offset_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "max_user_deferrals_count", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_users", "true"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_ErrorHandling(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register error mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`.*Access denied.*`),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_Update(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register local mocks directly
	macOSMock := localMocks.GetMock()
	macOSMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: localMocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Updated macOS Software Update Configuration"),
				),
			},
		},
	})
}

func TestAccMacOSSoftwareUpdateConfigurationResource_Basic(t *testing.T) {
	// Skip acceptance tests unless explicitly enabled
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: localMocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacOSSoftwareUpdateConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSSoftwareUpdateConfigurationExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "config_data_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "firmware_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "all_other_update_behavior", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_devices", "true"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper function to check if the resource exists
func testAccCheckMacOSSoftwareUpdateConfigurationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		// Add code to verify the resource exists in the real API
		// This is only needed for acceptance tests

		return nil
	}
}

// Helper function to check if the resource was destroyed
func testAccCheckMacOSSoftwareUpdateConfigurationDestroy(s *terraform.State) error {
	// Add code to verify the resource was destroyed in the real API
	// This is only needed for acceptance tests

	return nil
}

// Helper functions to generate test configurations
func testConfigBasic() string {
	return unitTestProviderConfig + testConfigBasicTemplate
}

func testConfigMinimal() string {
	return unitTestProviderConfig + testConfigMinimalTemplate
}

func testConfigMaximal() string {
	return unitTestProviderConfig + testConfigMaximalTemplate
}

func testConfigUpdate() string {
	return unitTestProviderConfig + testConfigUpdateTemplate
}

func testConfigGroupAssignments() string {
	return unitTestProviderConfig + testConfigGroupAssignmentsTemplate
}

func testAccConfigBasic() string {
	return accTestProviderConfig + testConfigBasicTemplate
}

func testAccConfigMinimal() string {
	return accTestProviderConfig + testConfigMinimalTemplate
}

func testAccConfigMaximal() string {
	return accTestProviderConfig + testConfigMaximalTemplate
}

func testAccConfigUpdate() string {
	return accTestProviderConfig + testConfigUpdateTemplate
}

// Setup test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "")
	os.Setenv("TF_VAR_tenant_id", "00000000-0000-0000-0000-000000000001")
	os.Setenv("TF_VAR_client_id", "11111111-1111-1111-1111-111111111111")
	os.Setenv("TF_VAR_client_secret", "mock-secret-value")

	// Clean up environment variables after the test
	t.Cleanup(func() {
		os.Unsetenv("TF_ACC")
		os.Unsetenv("TF_VAR_tenant_id")
		os.Unsetenv("TF_VAR_client_id")
		os.Unsetenv("TF_VAR_client_secret")
	})
}

// Helper function to check if the resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
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

// Helper function for acceptance test prechecks
func testAccPreCheck(t *testing.T) {
	// Check for required environment variables for acceptance tests
	requiredEnvVars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
	}

	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", v)
		}
	}
}

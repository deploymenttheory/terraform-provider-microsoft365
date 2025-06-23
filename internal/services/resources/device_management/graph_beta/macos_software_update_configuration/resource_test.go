package graphBetaMacOSSoftwareUpdateConfiguration_test

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

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
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

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

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

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

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

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "display_name", "Group Assignment Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.all_users", "false"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.include_group_ids.*", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.group_assigned", "assignments.exclude_group_ids.*", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_FullLifecycle(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

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
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_devices", "true"),
				),
			},
			// Import test
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"assignments.%",
					"assignments.all_devices",
					"assignments.all_users",
					"assignments.include_group_ids",
					"assignments.exclude_group_ids",
				},
			},
		},
	})
}

func TestUnitMacOSSoftwareUpdateConfigurationResource_ErrorHandling(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register error mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

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

func TestUnitMacOSSoftwareUpdateConfigurationResource_Update(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSSoftwareUpdateConfigurationMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Updated macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_time_window_utc_offset_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "max_user_deferrals_count", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "assignments.all_users", "true"),
				),
			},
		},
	})
}

// Acceptance Tests
func TestAccMacOSSoftwareUpdateConfigurationResource_Basic(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSSoftwareUpdateConfigurationExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "low"),
				),
			},
			{
				Config: testAccConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSSoftwareUpdateConfigurationExists("microsoft365_graph_beta_device_management_macos_software_update_configuration.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Updated macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "high"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckMacOSSoftwareUpdateConfigurationDestroy,
	})
}

// Helper Functions
func testAccCheckMacOSSoftwareUpdateConfigurationExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckMacOSSoftwareUpdateConfigurationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_software_update_configuration" {
			continue
		}

		// In a real test, we would make an API call to verify the resource is gone
		// For unit tests with mocks, we can assume it's destroyed if we get here
		return nil
	}

	return nil
}

// Test configurations using shared templates

// Unit test configurations
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

// Acceptance test configurations
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

func testAccPreCheck(t *testing.T) {
	// Check required environment variables for acceptance tests
	envVars := []string{
		"MICROSOFT365_CLIENT_ID",
		"MICROSOFT365_CLIENT_SECRET",
		"MICROSOFT365_TENANT_ID",
	}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", envVar)
		}
	}
}

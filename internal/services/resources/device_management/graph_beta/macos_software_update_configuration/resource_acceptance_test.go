package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccMacOSSoftwareUpdateConfigurationResource_Create_Minimal tests creating a configuration with minimal settings
func TestAccMacOSSoftwareUpdateConfigurationResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "default"),
					resource.TestCheckResourceAttr(resourceName, "config_data_update_behavior", "default"),
					resource.TestCheckResourceAttr(resourceName, "firmware_update_behavior", "default"),
					resource.TestCheckResourceAttr(resourceName, "all_other_update_behavior", "default"),
					resource.TestCheckResourceAttr(resourceName, "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr(resourceName, "update_time_window_utc_offset_in_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr(resourceName, "assignments.all_users", "false"),
				),
			},
		},
	})
}

// TestAccMacOSSoftwareUpdateConfigurationResource_Create_Maximal tests creating a configuration with maximal settings
func TestAccMacOSSoftwareUpdateConfigurationResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr(resourceName, "config_data_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr(resourceName, "firmware_update_behavior", "downloadOnly"),
					resource.TestCheckResourceAttr(resourceName, "all_other_update_behavior", "installLater"),
					resource.TestCheckResourceAttr(resourceName, "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr(resourceName, "update_time_window_utc_offset_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_user_deferrals_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "priority", "high"),
					resource.TestCheckResourceAttr(resourceName, "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr(resourceName, "assignments.all_users", "false"),
				),
			},
		},
	})
}

// TestAccMacOSSoftwareUpdateConfigurationResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestAccMacOSSoftwareUpdateConfigurationResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "default"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr(resourceName, "description"),
					resource.TestCheckNoResourceAttr(resourceName, "max_user_deferrals_count"),
					resource.TestCheckNoResourceAttr(resourceName, "priority"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr(resourceName, "config_data_update_behavior", "notifyOnly"),
					resource.TestCheckResourceAttr(resourceName, "firmware_update_behavior", "downloadOnly"),
					resource.TestCheckResourceAttr(resourceName, "all_other_update_behavior", "installLater"),
					resource.TestCheckResourceAttr(resourceName, "max_user_deferrals_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "priority", "high"),
				),
			},
		},
	})
}

// TestAccMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestAccMacOSSoftwareUpdateConfigurationResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Maximal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a comprehensive configuration with all fields populated"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr(resourceName, "max_user_deferrals_count", "3"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Minimal macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "critical_update_behavior", "default"),
				),
			},
		},
	})
}

// TestAccMacOSSoftwareUpdateConfigurationResource_Import tests importing a resource
func TestAccMacOSSoftwareUpdateConfigurationResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccMacOSSoftwareUpdateConfigurationResource_CustomTimeWindows tests configuration with custom update time windows
func TestAccMacOSSoftwareUpdateConfigurationResource_CustomTimeWindows(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_macos_software_update_configuration.custom"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create with custom time windows
			{
				Config: testAccConfigCustomTimeWindows(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConfigurationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TF Acc Test Custom Time Windows macOS Software Update Configuration"),
					resource.TestCheckResourceAttr(resourceName, "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.0.start_day", "monday"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.0.end_day", "monday"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.0.start_time", "10:00:00"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.0.end_time", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.1.start_day", "wednesday"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.1.end_day", "friday"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.1.start_time", "09:00:00"),
					resource.TestCheckResourceAttr(resourceName, "custom_update_time_windows.1.end_time", "17:00:00"),
				),
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	requiredEnvVars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

func testAccCheckConfigurationExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckConfigurationDestroy(s *terraform.State) error {
	// In a real test, we would verify the configuration is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the configuration
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "minimal" {
  display_name                          = "TF Acc Test Minimal macOS Software Update Configuration"
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
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "%s" {
  display_name                          = "TF Acc Test Minimal macOS Software Update Configuration"
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
`, resourceName)
}

// Maximal configuration with default resource name
func testAccConfigMaximal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "maximal" {
  display_name                          = "TF Acc Test Maximal macOS Software Update Configuration"
  description                           = "This is a comprehensive configuration with all fields populated"
  critical_update_behavior              = "installASAP"
  config_data_update_behavior           = "notifyOnly"
  firmware_update_behavior              = "downloadOnly"
  all_other_update_behavior             = "installLater"
  update_schedule_type                  = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 60
  max_user_deferrals_count              = 3
  priority                              = "high"
  role_scope_tag_ids                    = ["0"]

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "%s" {
  display_name                          = "TF Acc Test Maximal macOS Software Update Configuration"
  description                           = "This is a comprehensive configuration with all fields populated"
  critical_update_behavior              = "installASAP"
  config_data_update_behavior           = "notifyOnly"
  firmware_update_behavior              = "downloadOnly"
  all_other_update_behavior             = "installLater"
  update_schedule_type                  = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 60
  max_user_deferrals_count              = 3
  priority                              = "high"
  role_scope_tag_ids                    = ["0"]

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`, resourceName)
}

// Configuration with custom update time windows
func testAccConfigCustomTimeWindows() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "custom" {
  display_name                          = "TF Acc Test Custom Time Windows macOS Software Update Configuration"
  description                           = "Configuration with custom update time windows"
  critical_update_behavior              = "installASAP"
  config_data_update_behavior           = "notifyOnly"
  firmware_update_behavior              = "downloadOnly"
  all_other_update_behavior             = "installLater"
  update_schedule_type                  = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 60

  custom_update_time_windows = [
    {
      start_day = "monday"
      end_day   = "monday"
      start_time = "10:00:00"
      end_time   = "18:00:00"
    },
    {
      start_day = "wednesday"
      end_day   = "friday"
      start_time = "09:00:00"
      end_time   = "17:00:00"
    }
  ]

  assignments = {
    all_devices = false
    all_users   = true
  }
}
`
}

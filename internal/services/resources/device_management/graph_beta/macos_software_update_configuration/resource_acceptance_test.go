package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceMacOSSoftwareUpdateConfiguration_01_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccMacOSSoftwareUpdateConfigurationConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Acceptance macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "critical_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "config_data_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "firmware_update_behavior", "installASAP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "all_other_update_behavior", "installASAP"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_software_update_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccMacOSSoftwareUpdateConfigurationConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Acceptance macOS Software Update Configuration - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "updateDuringTimeWindows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "priority", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "max_user_deferrals_count", "3"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccMacOSSoftwareUpdateConfigurationConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "display_name", "Test Acceptance macOS Software Update Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test", "update_schedule_type", "alwaysUpdate"),
				),
			},
		},
	})
}

func TestAccResourceMacOSSoftwareUpdateConfiguration_02_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccMacOSSoftwareUpdateConfigurationConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_software_update_configuration.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test_assignments", "display_name", "Test macOS Software Update Configuration with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_software_update_configuration.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccResourceMacOSSoftwareUpdateConfiguration_02_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingUpdateScheduleType(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingCriticalUpdateBehavior(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingConfigDataUpdateBehavior(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingFirmwareUpdateBehavior(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_missingAllOtherUpdateBehavior(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccResourceMacOSSoftwareUpdateConfiguration_04_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_invalidUpdateScheduleType(),
				ExpectError: regexp.MustCompile("Attribute update_schedule_type value must be one of"),
			},
			{
				Config:      testAccMacOSSoftwareUpdateConfigurationConfig_invalidPriority(),
				ExpectError: regexp.MustCompile("Attribute priority value must be one of"),
			},
		},
	})
}

func testAccMacOSSoftwareUpdateConfigurationConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name           = "Test Acceptance macOS Software Update Configuration"
  update_schedule_type   = "alwaysUpdate"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name           = "Test Acceptance macOS Software Update Configuration - Updated"
  description           = "Updated description for acceptance testing"
  update_schedule_type   = "updateDuringTimeWindows"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = -480
  max_user_deferrals_count = 3
  priority = "high"
  role_scope_tag_ids = ["0", "1"]
  
  custom_update_time_windows = [
    {
      start_day = "monday"
      end_day = "friday"
      start_time = "02:00:00"
      end_time = "06:00:00"
    }
  ]
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_assignments" {
  display_name           = "Test macOS Software Update Configuration with Assignments"
  update_schedule_type   = "alwaysUpdate"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`)
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  update_schedule_type   = "alwaysUpdate"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingUpdateScheduleType() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingCriticalUpdateBehavior() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "alwaysUpdate"
  config_data_update_behavior = "downloadInstallRestart"
  firmware_update_behavior = "downloadInstallRestart"
  all_other_update_behavior = "downloadInstallRestart"
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingConfigDataUpdateBehavior() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "alwaysUpdate"
  critical_update_behavior = "downloadInstallRestart"
  firmware_update_behavior = "downloadInstallRestart"
  all_other_update_behavior = "downloadInstallRestart"
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingFirmwareUpdateBehavior() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "alwaysUpdate"
  critical_update_behavior = "downloadInstallRestart"
  config_data_update_behavior = "downloadInstallRestart"
  all_other_update_behavior = "downloadInstallRestart"
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_missingAllOtherUpdateBehavior() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "alwaysUpdate"
  critical_update_behavior = "downloadInstallRestart"
  config_data_update_behavior = "downloadInstallRestart"
  firmware_update_behavior = "downloadInstallRestart"
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_invalidUpdateScheduleType() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "invalid"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
}
`
}

func testAccMacOSSoftwareUpdateConfigurationConfig_invalidPriority() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test" {
  display_name = "Test Configuration"
  update_schedule_type = "alwaysUpdate"
  critical_update_behavior = "installASAP"
  config_data_update_behavior = "installASAP"
  firmware_update_behavior = "installASAP"
  all_other_update_behavior = "installASAP"
  update_time_window_utc_offset_in_minutes = 0
  priority = "invalid"
}
`
}

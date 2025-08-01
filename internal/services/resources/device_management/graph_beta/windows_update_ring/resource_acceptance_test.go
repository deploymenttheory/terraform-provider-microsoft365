package graphBetaWindowsUpdateRing_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsUpdateRingResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccWindowsUpdateRingConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_update_ring.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Acceptance Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "microsoft_update_service_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "drivers_excluded", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "quality_updates_deferral_period_in_days", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "feature_updates_deferral_period_in_days", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "allow_windows11_upgrade", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "skip_checks_before_restart", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "userDefined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "feature_updates_rollback_window_in_days", "10"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_update_ring.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccWindowsUpdateRingConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_update_ring.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Acceptance Windows Update Ring - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "business_ready_updates_only", "businessReadyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "quality_updates_deferral_period_in_days", "7"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "feature_updates_deferral_period_in_days", "14"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccWindowsUpdateRingConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_update_ring.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Acceptance Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "userDefined"),
				),
			},
		},
	})
}

func TestAccWindowsUpdateRingResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccWindowsUpdateRingConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_update_ring.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test_assignments", "display_name", "Test Windows Update Ring with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccWindowsUpdateRingResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsUpdateRingConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingMicrosoftUpdateServiceAllowed(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingDriversExcluded(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingQualityUpdatesDeferralPeriod(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingFeatureUpdatesDeferralPeriod(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingAllowWindows11Upgrade(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingSkipChecksBeforeRestart(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingAutomaticUpdateMode(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_missingFeatureUpdatesRollbackWindow(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccWindowsUpdateRingResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsUpdateRingConfig_invalidAutomaticUpdateMode(),
				ExpectError: regexp.MustCompile("Attribute automatic_update_mode value must be one of"),
			},
			{
				Config:      testAccWindowsUpdateRingConfig_invalidBusinessReadyUpdatesOnly(),
				ExpectError: regexp.MustCompile("Attribute business_ready_updates_only value must be one of"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("ARM_TENANT_ID") == "" {
		t.Skip("ARM_TENANT_ID must be set for acceptance tests")
	}
	if os.Getenv("ARM_CLIENT_ID") == "" {
		t.Skip("ARM_CLIENT_ID must be set for acceptance tests")
	}
	if os.Getenv("ARM_CLIENT_SECRET") == "" {
		t.Skip("ARM_CLIENT_SECRET must be set for acceptance tests")
	}
}

func testAccWindowsUpdateRingConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Acceptance Windows Update Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Acceptance Windows Update Ring - Updated"
  description                              = "Updated description for acceptance testing"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 7
  feature_updates_deferral_period_in_days  = 14
  allow_windows11_upgrade                  = false
  skip_checks_before_restart               = true
  automatic_update_mode                    = "autoInstallAndRebootAtScheduledTime"
  business_ready_updates_only              = "businessReadyOnly"
  delivery_optimization_mode               = "httpWithPeeringNat"
  prerelease_features                      = "settingsOnly"
  update_weeks                             = "firstWeek"
  active_hours_start                       = "09:00:00"
  active_hours_end                         = "17:00:00"
  user_pause_access                        = "disabled"
  feature_updates_rollback_window_in_days  = 10
  engaged_restart_deadline_in_days         = 3
  role_scope_tag_ids                       = ["0", "1"]
  
  deadline_settings = {
    deadline_for_feature_updates_in_days  = 7
    deadline_for_quality_updates_in_days  = 2
    deadline_grace_period_in_days         = 1
    postpone_reboot_until_after_deadline  = true
  }
}
`
}

func testAccWindowsUpdateRingConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "test_assignments" {
  display_name                             = "Test Windows Update Ring with Assignments"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`)
}

func testAccWindowsUpdateRingConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingMicrosoftUpdateServiceAllowed() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingDriversExcluded() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingQualityUpdatesDeferralPeriod() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingFeatureUpdatesDeferralPeriod() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingAllowWindows11Upgrade() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingSkipChecksBeforeRestart() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingAutomaticUpdateMode() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_missingFeatureUpdatesRollbackWindow() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
}
`
}

func testAccWindowsUpdateRingConfig_invalidAutomaticUpdateMode() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "invalid"
  feature_updates_rollback_window_in_days  = 10
}
`
}

func testAccWindowsUpdateRingConfig_invalidBusinessReadyUpdatesOnly() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                             = "Test Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  business_ready_updates_only              = "invalid"
  feature_updates_rollback_window_in_days  = 10
}
`
}
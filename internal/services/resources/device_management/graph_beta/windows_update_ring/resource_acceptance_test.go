package graphBetaWindowsUpdateRing_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsUpdateRingResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsUpdateRingDestroy,
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
		},
	})
}

func TestAccWindowsUpdateRingResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsUpdateRingDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
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
		},
	})
}

func TestAccWindowsUpdateRingResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsUpdateRingDestroy,
		Steps: []resource.TestStep{
			// Create with all assignment types
			{
				Config: testAccWindowsUpdateRingConfig_comprehensiveAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "display_name", "Test All Assignment Types Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					// Verify role scope tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.assignments", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccWindowsUpdateRingResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsUpdateRingDestroy,
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
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsUpdateRingDestroy,
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

func testAccWindowsUpdateRingConfig_minimal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsUpdateRingConfig_maximal() string {
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsUpdateRingConfig_comprehensiveAssignments() string {
	groups := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
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

// testAccCheckWindowsUpdateRingDestroy verifies that Windows update rings have been destroyed
func testAccCheckWindowsUpdateRingDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_update_ring" {
			continue
		}

		// Attempt to get the Windows update ring by ID
		_, err := graphClient.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if Windows update ring %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("Windows update ring %s still exists", rs.Primary.ID)
	}

	return nil
}

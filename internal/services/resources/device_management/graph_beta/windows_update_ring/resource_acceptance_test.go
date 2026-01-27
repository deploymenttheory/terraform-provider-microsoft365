package graphBetaWindowsUpdateRing_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateRing "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_ring"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testResource = graphBetaWindowsUpdateRing.WindowsUpdateRingTestResource{}

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// Test 001: Scenario 1 - Notify Download
func TestAccResourceWindowsUpdateRing_01_Scenario_NotifyDownload(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 1: Notify Download")
				},
				Config: loadAcceptanceTestTerraform("scenario_001_notify_download.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-001-notify-download-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("automatic_update_mode").HasValue("notifyDownload"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("description").HasValue("Scenario 1: Notify Download"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 1")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 002: Scenario 2 - Auto Install at Maintenance Time
func TestAccResourceWindowsUpdateRing_02_Scenario_AutoInstallMaintenance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 2: Auto Install at Maintenance Time")
				},
				Config: loadAcceptanceTestTerraform("scenario_002_auto_install_maintenance.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-002-auto-install-maintenance-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("automatic_update_mode").HasValue("autoInstallAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("active_hours_start").HasValue("08:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("active_hours_end").HasValue("17:00:00"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 2")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 003: Scenario 3 - Auto Install and Reboot at Maintenance Time
func TestAccResourceWindowsUpdateRing_03_Scenario_AutoRebootMaintenance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 3: Auto Install and Reboot at Maintenance Time")
				},
				Config: loadAcceptanceTestTerraform("scenario_003_auto_reboot_maintenance.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-003-auto-reboot-maintenance-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("active_hours_start").HasValue("08:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("active_hours_end").HasValue("17:00:00"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 3")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 004: Scenario 4 - Auto Install and Restart at Scheduled Time
func TestAccResourceWindowsUpdateRing_04_Scenario_ScheduledInstall(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 4: Auto Install and Restart at Scheduled Time")
				},
				Config: loadAcceptanceTestTerraform("scenario_004_scheduled_install.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-004-scheduled-install-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtScheduledTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("scheduled_install_day").HasValue("everyday"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("scheduled_install_time").HasValue("03:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("update_weeks").HasValue("everyWeek"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 4")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_004",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 005: Scenario 5 - Auto Install and Reboot Without End User Control
func TestAccResourceWindowsUpdateRing_05_Scenario_NoEndUserControl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 5: Auto Install and Reboot Without End User Control")
				},
				Config: loadAcceptanceTestTerraform("scenario_005_no_end_user_control.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-005-no-end-user-control-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("automatic_update_mode").HasValue("autoInstallAndRebootWithoutEndUserControl"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 5")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_005",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 006: Scenario 6 - Windows Default (Reset)
func TestAccResourceWindowsUpdateRing_06_Scenario_WindowsDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating Scenario 6: Windows Default (Reset)")
				},
				Config: loadAcceptanceTestTerraform("scenario_006_windows_default.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-update-ring-006-windows-default-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("automatic_update_mode").HasValue("windowsDefault"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("update_notification_level").HasValue("disableAllNotifications"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing Scenario 6")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_006",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 007: Full Lifecycle Through All Scenarios
func TestAccResourceWindowsUpdateRing_07_FullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Creating with Scenario 1")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_1_scenario_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("notifyDownload"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Updating to Scenario 2")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_2_scenario_002.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("active_hours_start").HasValue("08:00:00"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Updating to Scenario 3")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_3_scenario_003.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtMaintenanceTime"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Updating to Scenario 4")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_4_scenario_004.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtScheduledTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("scheduled_install_day").HasValue("everyday"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Updating to Scenario 5")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_5_scenario_005.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootWithoutEndUserControl"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Updating to Scenario 6")
				},
				Config: loadAcceptanceTestTerraform("lifecycle_step_6_scenario_006.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("windowsDefault"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("update_notification_level").HasValue("disableAllNotifications"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Lifecycle: Importing final state")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_007",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 008: Minimal Assignments
func TestAccResourceWindowsUpdateRing_08_AssignmentsMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_008").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing minimal assignments")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_008",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 009: Maximal Assignments
func TestAccResourceWindowsUpdateRing_09_AssignmentsMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Creating with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_009").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_009").Key("assignments.#").HasValue("5"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Importing maximal assignments")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_009",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 010: Assignments Lifecycle - Minimal to Maximal
func TestAccResourceWindowsUpdateRing_10_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Lifecycle: Creating with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_010").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_010").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Lifecycle: Updating to maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_010").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_010").Key("assignments.#").HasValue("5"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Lifecycle: Importing final state")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_010",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 011: Assignments Lifecycle - Maximal to Minimal
func TestAccResourceWindowsUpdateRing_11_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateRing.ResourceName,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Downgrade: Creating with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_011").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_011").Key("assignments.#").HasValue("5"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Downgrade: Updating to minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("assignments_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows update ring", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_011").ExistsInGraph(testResource),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_011").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Assignments Downgrade: Importing final state")
				},
				ResourceName:            graphBetaWindowsUpdateRing.ResourceName + ".test_011",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 012: Error and Validation Testing
func TestAccResourceWindowsUpdateRing_12_ValidationErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Testing invalid automatic_update_mode")
				},
				Config:      loadAcceptanceTestTerraform("validation_invalid_update_mode.tf"),
				ExpectError: regexp.MustCompile(`Attribute automatic_update_mode value must be one of`),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Testing invalid business_ready_updates_only")
				},
				Config:      loadAcceptanceTestTerraform("validation_invalid_business_ready.tf"),
				ExpectError: regexp.MustCompile(`Attribute business_ready_updates_only value must be one of`),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsUpdateRing.ResourceName, "Testing missing display_name")
				},
				Config:      loadAcceptanceTestTerraform("validation_missing_display_name.tf"),
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

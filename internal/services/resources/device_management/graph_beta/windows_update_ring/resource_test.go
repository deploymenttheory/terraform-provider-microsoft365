package graphBetaWindowsUpdateRing_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateRing "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_ring"
	windowsUpdateRingMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_ring/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *windowsUpdateRingMocks.WindowsUpdateRingMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	windowsUpdateRingMock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	windowsUpdateRingMock.RegisterMocks()
	return mockClient, windowsUpdateRingMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *windowsUpdateRingMocks.WindowsUpdateRingMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	windowsUpdateRingMock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	windowsUpdateRingMock.RegisterErrorMocks()
	return mockClient, windowsUpdateRingMock
}

// Test 001: Scenario 1 - Notify Download
func TestWindowsUpdateRingResource_001_Scenario_NotifyDownload(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_001_notify_download.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-001-notify-download-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("automatic_update_mode").HasValue("notifyDownload"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_001").Key("description").HasValue("Scenario 1: Notify Download"),
				),
			},
		},
	})
}

// Test 002: Scenario 2 - Auto Install at Maintenance Time
func TestWindowsUpdateRingResource_002_Scenario_AutoInstallMaintenance(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_002_auto_install_maintenance.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-002-auto-install-maintenance-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("automatic_update_mode").HasValue("autoInstallAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("active_hours_start").HasValue("08:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_002").Key("active_hours_end").HasValue("17:00:00"),
				),
			},
		},
	})
}

// Test 003: Scenario 3 - Auto Install and Reboot at Maintenance Time
func TestWindowsUpdateRingResource_003_Scenario_AutoRebootMaintenance(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_003_auto_reboot_maintenance.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-003-auto-reboot-maintenance-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("active_hours_start").HasValue("08:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_003").Key("active_hours_end").HasValue("17:00:00"),
				),
			},
		},
	})
}

// Test 004: Scenario 4 - Auto Install and Restart at Scheduled Time
func TestWindowsUpdateRingResource_004_Scenario_ScheduledInstall(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_004_scheduled_install.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-004-scheduled-install-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtScheduledTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("scheduled_install_day").HasValue("everyday"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("scheduled_install_time").HasValue("03:00:00"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_004").Key("update_weeks").HasValue("everyWeek"),
				),
			},
		},
	})
}

// Test 005: Scenario 5 - Auto Install and Reboot Without End User Control
func TestWindowsUpdateRingResource_005_Scenario_NoEndUserControl(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_005_no_end_user_control.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-005-no-end-user-control-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_005").Key("automatic_update_mode").HasValue("autoInstallAndRebootWithoutEndUserControl"),
				),
			},
		},
	})
}

// Test 006: Scenario 6 - Windows Default (Reset)
func TestWindowsUpdateRingResource_006_Scenario_WindowsDefault(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_006_windows_default.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-windows-update-ring-006-windows-default-[a-z0-9]{8}$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("automatic_update_mode").HasValue("windowsDefault"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_006").Key("update_notification_level").HasValue("disableAllNotifications"),
				),
			},
		},
	})
}

// Test 007: Full Lifecycle Through All Scenarios
func TestWindowsUpdateRingResource_007_FullLifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("lifecycle_step_1_scenario_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_007").Key("automatic_update_mode").HasValue("notifyDownload"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_2_scenario_002.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAtMaintenanceTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("active_hours_start").HasValue("08:00:00"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_3_scenario_003.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtMaintenanceTime"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_4_scenario_004.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootAtScheduledTime"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("scheduled_install_day").HasValue("everyday"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_5_scenario_005.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_007").Key("automatic_update_mode").HasValue("autoInstallAndRebootWithoutEndUserControl"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_6_scenario_006.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("automatic_update_mode").HasValue("windowsDefault"),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_007").Key("update_notification_level").HasValue("disableAllNotifications"),
				),
			},
		},
	})
}

// Test 008: Minimal Assignments
func TestWindowsUpdateRingResource_008_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 009: Maximal Assignments
func TestWindowsUpdateRingResource_009_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_009").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsUpdateRing.ResourceName+".test_009").Key("assignments.#").HasValue("5"),
				),
			},
		},
	})
}

// Test 010: Assignments Lifecycle - Minimal to Maximal
func TestWindowsUpdateRingResource_010_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_010").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("assignments_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_010").Key("assignments.#").HasValue("5"),
				),
			},
		},
	})
}

// Test 011: Assignments Lifecycle - Maximal to Minimal
func TestWindowsUpdateRingResource_011_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_011").Key("assignments.#").HasValue("5"),
				),
			},
			{
				Config: loadUnitTestTerraform("assignments_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdateRing.ResourceName + ".test_011").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 012: Error and Validation Testing
func TestWindowsUpdateRingResource_012_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("validation_invalid_update_mode.tf"),
				ExpectError: regexp.MustCompile(`Attribute automatic_update_mode value must be one of`),
			},
			{
				Config:      loadUnitTestTerraform("validation_invalid_business_ready.tf"),
				ExpectError: regexp.MustCompile(`Attribute business_ready_updates_only value must be one of`),
			},
			{
				Config:      loadUnitTestTerraform("validation_missing_display_name.tf"),
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

package graphBetaDeviceAndAppManagementAppAssignment_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_assignment"
	mobileAppAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *mobileAppAssignmentMocks.MobileAppAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	mobileAppAssignmentMock := &mobileAppAssignmentMocks.MobileAppAssignmentMock{}
	mobileAppAssignmentMock.RegisterMocks()
	return mockClient, mobileAppAssignmentMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// diagnosticRegexp builds a whitespace tolerant pattern for a diagnostic phrase, since
// Terraform wraps diagnostic text across lines when rendering it.
func diagnosticRegexp(phrase string) *regexp.Regexp {
	fields := strings.Fields(phrase)
	for i, field := range fields {
		fields[i] = regexp.QuoteMeta(field)
	}
	return regexp.MustCompile(`(?s)` + strings.Join(fields, `\s+`))
}

// checkSettingSent asserts on the settings payload the provider actually sent for the given
// settings block, rather than only on what ended up in Terraform state. Which fields leave the
// provider is the whole point of these tests, and state alone cannot show it.
func checkSettingSent(t *testing.T, mobileAppAssignmentMock *mobileAppAssignmentMocks.MobileAppAssignmentMock, block, field string, wantPresent bool, wantValue any) resource.TestCheckFunc {
	t.Helper()

	return func(*terraform.State) error {
		settings, ok := mobileAppAssignmentMock.SettingsSent(block)
		if !ok {
			return fmt.Errorf("no assignment request was captured for settings block %s", block)
		}

		value, present := settings[field]
		if present != wantPresent {
			return fmt.Errorf("%s present in request = %v, want %v (payload: %v)", field, present, wantPresent, settings)
		}
		if wantPresent && value != wantValue {
			return fmt.Errorf("%s sent as %v, want %v", field, value, wantValue)
		}
		return nil
	}
}

// Test 001: Scenario 1 - available intent with is_removable omitted
//
// The Intune service rejects isRemovable for any intent other than required, so the provider
// must not send the field at all. This is the configuration reported in issue #3692.
func TestUnitResourceMobileAppAssignment_01_Scenario_Available_IosStore_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_available_ios_store_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("source").HasValue("direct"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("target.target_type").HasValue("groupAssignment"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("settings.ios_store.is_removable").DoesNotExist(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("settings.ios_store.uninstall_on_device_removal").DoesNotExist(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_001").Key("settings.ios_store.prevent_managed_app_backup").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosStore", "isRemovable", false, nil),
				),
			},
		},
	})
}

// Test 002: Scenario 2 - available intent for the ios_vpp settings block
func TestUnitResourceMobileAppAssignment_02_Scenario_Available_IosVpp(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_available_ios_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_002").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_002").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_002").Key("settings.ios_vpp.use_device_licensing").HasValue("false"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_002").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", false, nil),
				),
			},
		},
	})
}

// Test 003: Scenario 3 - available intent for the ios_lob settings block
func TestUnitResourceMobileAppAssignment_03_Scenario_Available_IosLob(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_scenario_available_ios_lob.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_003").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_003").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_003").Key("settings.ios_lob.is_removable").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosLob", "isRemovable", false, nil),
				),
			},
		},
	})
}

// Test 004: Scenario 4 - required intent with every ios_vpp setting configured
//
// required is the one intent that accepts is_removable, so it must still reach the API.
func TestUnitResourceMobileAppAssignment_04_Scenario_Required_IosVpp_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_scenario_required_ios_vpp_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("intent").HasValue("required"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("settings.ios_vpp.is_removable").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("settings.ios_vpp.prevent_auto_app_update").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("settings.ios_vpp.prevent_managed_app_backup").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("settings.ios_vpp.uninstall_on_device_removal").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_004").Key("settings.ios_vpp.use_device_licensing").HasValue("true"),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", true, true),
				),
			},
		},
	})
}

// Test 005: Scenario 5 - required intent with is_removable omitted
//
// An omitted attribute is not sent even for the one intent that supports it, so the Intune
// service default applies rather than a value the practitioner never asked for.
func TestUnitResourceMobileAppAssignment_05_Scenario_Required_IsRemovable_Omitted(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_scenario_required_is_removable_omitted.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_005").Key("intent").HasValue("required"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_005").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", false, nil),
				),
			},
		},
	})
}

// Test 006: Scenario 6 - uninstall intent for the ios_vpp settings block
func TestUnitResourceMobileAppAssignment_06_Scenario_Uninstall_IosVpp(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_scenario_uninstall_ios_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_006").Key("intent").HasValue("uninstall"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_006").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", false, nil),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "uninstallOnDeviceRemoval", false, nil),
				),
			},
		},
	})
}

// Test 007: Scenario 7 - uninstall intent with an empty settings block
//
// The service rejects isRemovable, uninstallOnDeviceRemoval and preventManagedAppBackup for an
// uninstall intent, so an empty block must produce an empty payload rather than a defaulted one.
func TestUnitResourceMobileAppAssignment_07_Scenario_Uninstall_IosStore(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_scenario_uninstall_ios_store.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_007").Key("intent").HasValue("uninstall"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_007").Key("settings.ios_store.is_removable").DoesNotExist(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_007").Key("settings.ios_store.uninstall_on_device_removal").DoesNotExist(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_007").Key("settings.ios_store.prevent_managed_app_backup").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosStore", "isRemovable", false, nil),
					checkSettingSent(t, mobileAppAssignmentMock, "iosStore", "uninstallOnDeviceRemoval", false, nil),
					checkSettingSent(t, mobileAppAssignmentMock, "iosStore", "preventManagedAppBackup", false, nil),
				),
			},
		},
	})
}

// Test 008: Scenario 8 - intent unknown at plan time
//
// Because the affected attributes carry no schema default, an omitted is_removable is null in
// the plan whether or not the intent is known, so no unknown value is introduced and no plan
// modification is needed to reconcile the two.
func TestUnitResourceMobileAppAssignment_08_Scenario_Unknown_Intent(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_scenario_unknown_intent.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_008").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_008").Key("settings.ios_vpp.is_removable").DoesNotExist(),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", false, nil),
				),
			},
		},
	})
}

// Test 009: Scenario 9 - settings lifecycle, updated in place
//
// Graph rejects a PATCH carrying intent, source or target with "Cannot patch read-only
// properties", so an update must send settings alone. Settings are also the only part of an
// assignment that does not force replacement.
func TestUnitResourceMobileAppAssignment_09_Lifecycle_Settings_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("009_lifecycle_settings_update_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName + ".test_009").Key("settings.ios_vpp.is_removable").HasValue("true"),
				),
			},
			{
				Config: loadUnitTestTerraform("009_lifecycle_settings_update_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_009").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_009").Key("settings.ios_vpp.is_removable").HasValue("false"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_009").Key("settings.ios_vpp.use_device_licensing").HasValue("true"),
					checkSettingSent(t, mobileAppAssignmentMock, "iosVpp", "isRemovable", true, false),
				),
			},
		},
	})
}

// Test 010: Scenario 10 - import by composite id
//
// An assignment cannot be addressed by its own id alone: mobile_app_id is required to query it
// and is not derivable, so the import id carries both.
func TestUnitResourceMobileAppAssignment_10_Import_CompositeId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("010_scenario_import.tf"),
			},
			{
				ResourceName: graphBetaMobileAppAssignment.ResourceName + ".test_010",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[graphBetaMobileAppAssignment.ResourceName+".test_010"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return rs.Primary.Attributes["mobile_app_id"] + ":" + rs.Primary.ID, nil
				},
				ImportStateVerify: true,
				// timeouts are configuration only and are never returned by the API.
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				ResourceName:  graphBetaMobileAppAssignment.ResourceName + ".test_010",
				ImportState:   true,
				ImportStateId: "not-a-composite-id",
				ExpectError:   diagnosticRegexp("Expected import ID in format 'mobileAppId:assignmentId'"),
			},
		},
	})
}

// Test 011: Scenario 11 - validation errors
//
// The service rejects these combinations with an HTTP 400. Catching them at plan time gives an
// actionable error rather than a failure part-way through an apply.
func TestUnitResourceMobileAppAssignment_11_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("011_validation_is_removable_available.tf"),
				ExpectError: diagnosticRegexp("`is_removable` can only be set when `intent` is `required`"),
			},
			{
				Config:      loadUnitTestTerraform("011_validation_is_removable_uninstall.tf"),
				ExpectError: diagnosticRegexp("`is_removable` can only be set when `intent` is `required`"),
			},
			{
				Config:      loadUnitTestTerraform("011_validation_is_removable_available_without_enrollment.tf"),
				ExpectError: diagnosticRegexp("`is_removable` can only be set when `intent` is `required`"),
			},
			{
				Config:      loadUnitTestTerraform("011_validation_uninstall_on_device_removal.tf"),
				ExpectError: diagnosticRegexp("`uninstall_on_device_removal` cannot be set when `intent` is `uninstall`"),
			},
			{
				Config:      loadUnitTestTerraform("011_validation_prevent_managed_app_backup.tf"),
				ExpectError: diagnosticRegexp("`prevent_managed_app_backup` cannot be set when `intent` is `uninstall`"),
			},
			{
				Config:      loadUnitTestTerraform("011_validation_auto_update_settings_required.tf"),
				ExpectError: diagnosticRegexp("`auto_update_settings` can only be set when `intent` is `available`"),
			},
		},
	})
}

// Test 012: Scenario 12 - win32_catalog settings round trip
//
// Read had no case for Win32CatalogAppAssignmentSettings and fell through to returning nil, so
// a refresh nulled the whole settings block in state. That was invisible while Read discarded
// its result, and became a permanent diff once Read was fixed.
func TestUnitResourceMobileAppAssignment_12_Settings_Win32Catalog(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("012_settings_win32_catalog.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("settings.win32_catalog.notifications").HasValue("showAll"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("settings.win32_catalog.delivery_optimization_priority").HasValue("foreground"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("settings.win32_catalog.auto_update_settings.auto_update_superseded_apps_state").HasValue("enabled"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("settings.win32_catalog.install_time_settings.use_local_time").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_012").Key("settings.win32_catalog.restart_settings.grace_period_in_minutes").HasValue("60"),
				),
			},
		},
	})
}

// Test 013: Scenario 13 - win32_lob settings round trip
//
// auto_update_settings was never mapped back, so a configured value drifted on every read.
func TestUnitResourceMobileAppAssignment_13_Settings_Win32Lob(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("013_settings_win32_lob.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("intent").HasValue("available"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("settings.win32_lob.notifications").HasValue("showReboot"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("settings.win32_lob.delivery_optimization_priority").HasValue("foreground"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("settings.win32_lob.auto_update_settings.auto_update_superseded_apps_state").HasValue("enabled"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("settings.win32_lob.install_time_settings.use_local_time").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_013").Key("settings.win32_lob.restart_settings.grace_period_in_minutes").HasValue("120"),
				),
			},
		},
	})
}

// Test 014: Scenario 14 - win_get settings round trip
func TestUnitResourceMobileAppAssignment_14_Settings_WinGet(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("014_settings_win_get.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_014").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_014").Key("settings.win_get.notifications").HasValue("showAll"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_014").Key("settings.win_get.install_time_settings.use_local_time").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_014").Key("settings.win_get.restart_settings.grace_period_in_minutes").HasValue("90"),
				),
			},
		},
	})
}

// Test 015: Scenario 15 - android_managed_store settings round trip
func TestUnitResourceMobileAppAssignment_15_Settings_AndroidManagedStore(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("015_settings_android_managed_store.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_015").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_015").Key("settings.android_managed_store.auto_update_mode").HasValue("priority"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_015").Key("settings.android_managed_store.android_managed_store_app_track_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 016: Scenario 16 - macos_lob settings round trip
func TestUnitResourceMobileAppAssignment_16_Settings_MacOsLob(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("016_settings_macos_lob.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_016").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_016").Key("settings.macos_lob.uninstall_on_device_removal").HasValue("true"),
				),
			},
		},
	})
}

// Test 017: Scenario 17 - macos_vpp settings round trip
func TestUnitResourceMobileAppAssignment_17_Settings_MacOsVpp(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("017_settings_macos_vpp.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_017").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_017").Key("settings.macos_vpp.prevent_auto_app_update").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_017").Key("settings.macos_vpp.prevent_managed_app_backup").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_017").Key("settings.macos_vpp.uninstall_on_device_removal").HasValue("true"),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_017").Key("settings.macos_vpp.use_device_licensing").HasValue("true"),
				),
			},
		},
	})
}

// Test 018: Scenario 18 - microsoft_store_for_business settings round trip
func TestUnitResourceMobileAppAssignment_18_Settings_MicrosoftStoreForBusiness(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("018_settings_microsoft_store_for_business.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_018").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_018").Key("settings.microsoft_store_for_business.use_device_context").HasValue("true"),
				),
			},
		},
	})
}

// Test 019: Scenario 19 - windows_app_x settings round trip
func TestUnitResourceMobileAppAssignment_19_Settings_WindowsAppX(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("019_settings_windows_app_x.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_019").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_019").Key("settings.windows_app_x.use_device_context").HasValue("true"),
				),
			},
		},
	})
}

// Test 020: Scenario 20 - windows_universal_app_x settings round trip
func TestUnitResourceMobileAppAssignment_20_Settings_WindowsUniversalAppX(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mobileAppAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mobileAppAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("020_settings_windows_universal_app_x.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_020").Key("id").Exists(),
					check.That(graphBetaMobileAppAssignment.ResourceName+".test_020").Key("settings.windows_universal_app_x.use_device_context").HasValue("true"),
				),
			},
		},
	})
}

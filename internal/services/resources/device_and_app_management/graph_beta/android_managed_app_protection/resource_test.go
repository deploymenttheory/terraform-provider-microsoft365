package graphBetaDeviceAndAppManagementAndroidManagedAppProtection_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAndroidManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/android_managed_app_protection"
	wmapmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/android_managed_app_protection/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *wmapmocks.AndroidManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	amapMock := &wmapmocks.AndroidManagedAppProtectionMock{}
	amapMock.RegisterMocks()
	return mockClient, amapMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *wmapmocks.AndroidManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	amapMock := &wmapmocks.AndroidManagedAppProtectionMock{}
	amapMock.RegisterErrorMocks()
	return mockClient, amapMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceAndroidManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, amapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer amapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("display_name").HasValue("unit-test-android-managed-app-protection-minimal"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("print_blocked").HasValue("false"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("pin_required").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("encrypt_app_data").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceAndroidManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, amapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer amapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("display_name").HasValue("unit-test-android-managed-app-protection-maximal"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("description").HasValue("Maximal test configuration for Android managed app protection"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("blocked"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("screen_capture_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("pin_required").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_pin_length").HasValue("6"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("maximum_pin_retries").HasValue("10"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("pin_character_set").HasValue("alphanumericAndSymbol"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_required_os_version").HasValue("9.0"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_required_app_version").HasValue("2.0.0"),
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceAndroidManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, amapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer amapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-android-managed-app-protection-lifecycle"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("false"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-android-managed-app-protection-lifecycle"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("screen_capture_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("minimum_pin_length").HasValue("6"),
				),
			},
		},
	})
}

func TestUnitResourceAndroidManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, amapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer amapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-android-managed-app-protection-lifecycle"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("screen_capture_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("minimum_pin_length").HasValue("6"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-android-managed-app-protection-lifecycle"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("false"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
		},
	})
}

func TestUnitResourceAndroidManagedAppProtection_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, amapMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer amapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("005_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Android Managed App Protection data"),
			},
		},
	})
}

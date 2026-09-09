package graphBetaDeviceAndAppManagementIosManagedAppProtection_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIosManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_app_protection"
	imapmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_app_protection/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *imapmocks.IosManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	imapMock := &imapmocks.IosManagedAppProtectionMock{}
	imapMock.RegisterMocks()
	return mockClient, imapMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *imapmocks.IosManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	imapMock := &imapmocks.IosManagedAppProtectionMock{}
	imapMock.RegisterErrorMocks()
	return mockClient, imapMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceIosManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, imapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer imapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("display_name").HasValue("unit-test-ios-managed-app-protection-minimal"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("print_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("pin_required").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("face_id_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("app_data_encryption_type").HasValue("whenDeviceLocked"),
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceIosManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, imapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer imapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("display_name").HasValue("unit-test-ios-managed-app-protection-maximal"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("description").HasValue("Maximal test configuration for iOS managed app protection"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("blocked"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("face_id_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("third_party_keyboards_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("filter_open_in_to_only_managed_apps").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("app_data_encryption_type").HasValue("afterDeviceRestart"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("pin_required").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_pin_length").HasValue("6"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("maximum_pin_retries").HasValue("10"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("pin_character_set").HasValue("alphanumericAndSymbol"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_required_os_version").HasValue("15.0"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_required_app_version").HasValue("2.0.0"),
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceIosManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, imapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer imapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-ios-managed-app-protection-lifecycle"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-ios-managed-app-protection-lifecycle"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("face_id_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("minimum_pin_length").HasValue("6"),
				),
			},
		},
	})
}

func TestUnitResourceIosManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, imapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer imapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-ios-managed-app-protection-lifecycle"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("face_id_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("minimum_pin_length").HasValue("6"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-ios-managed-app-protection-lifecycle"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
		},
	})
}

func TestUnitResourceIosManagedAppProtection_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, imapMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer imapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("005_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid iOS Managed App Protection data"),
			},
		},
	})
}
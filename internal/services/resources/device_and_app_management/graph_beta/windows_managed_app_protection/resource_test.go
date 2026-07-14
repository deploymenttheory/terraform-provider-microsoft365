package graphBetaDeviceAndAppManagementWindowsManagedAppProtection_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/windows_managed_app_protection"
	wmapmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/windows_managed_app_protection/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *wmapmocks.WindowsManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	wmapMock := &wmapmocks.WindowsManagedAppProtectionMock{}
	wmapMock.RegisterMocks()
	return mockClient, wmapMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *wmapmocks.WindowsManagedAppProtectionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	wmapMock := &wmapmocks.WindowsManagedAppProtectionMock{}
	wmapMock.RegisterErrorMocks()
	return mockClient, wmapMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, wmapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer wmapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-managed-app-protection-minimal"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("print_blocked").HasValue("false"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("maximum_allowed_device_threat_level").HasValue("notConfigured"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("mobile_threat_defense_remediation_action").HasValue("block"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_wipe_is_enforced").HasValue("P90D"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_access_check").HasValue("P30D"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				ResourceName:      graphBetaWindowsManagedAppProtection.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, wmapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer wmapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-managed-app-protection-maximal"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("description").HasValue("Maximal test configuration for Windows managed app protection"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("print_blocked").HasValue("true"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("none"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("maximum_allowed_device_threat_level").HasValue("low"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("mobile_threat_defense_remediation_action").HasValue("wipe"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_required_os_version").HasValue("10.0.19041"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_warning_os_version").HasValue("10.0.18363"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_wipe_os_version").HasValue("10.0.17763"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_required_app_version").HasValue("1.0.0"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_warning_app_version").HasValue("1.1.0"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("minimum_wipe_app_version").HasValue("0.9.0"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("period_offline_before_wipe_is_enforced").HasValue("P30D"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("period_offline_before_access_check").HasValue("P7D"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsManagedAppProtection.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, wmapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer wmapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-managed-app-protection-lifecycle"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("false"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-managed-app-protection-lifecycle"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("true"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_003").Key("minimum_required_os_version").HasValue("10.0.19041"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, wmapMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer wmapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-managed-app-protection-lifecycle"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("true"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-managed-app-protection-lifecycle"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("false"),
					check.That(graphBetaWindowsManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsManagedAppProtection_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, wmapMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer wmapMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("005_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Windows Managed App Protection data"),
			},
		},
	})
}

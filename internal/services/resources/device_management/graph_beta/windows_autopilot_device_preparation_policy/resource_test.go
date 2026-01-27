package graphBetaWindowsAutopilotDevicePreparationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_device_preparation_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "name", "unit-test-windows-autopilot-device-preparation-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "description", "unit-test-windows-autopilot-device-preparation-policy-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "role_scope_tag_ids.*", "0"),

					// Device security group
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "device_security_group", "00000000-0000-0000-0000-000000000001"),

					// Deployment settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.deployment_mode", "enrollment_autopilot_dpp_deploymentmode_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.deployment_type", "enrollment_autopilot_dpp_deploymenttype_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.join_type", "enrollment_autopilot_dpp_jointype_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.account_type", "enrollment_autopilot_dpp_accountype_0"),

					// OOBE settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.custom_error_message", "Contact your organization's support person for help."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.allow_skip", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.allow_diagnostics", "false"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "assignments.include_group_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "assignments.include_group_ids.*", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "assignments.include_group_ids.*", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_02_MaximalSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "name", "unit-test-windows-autopilot-device-preparation-policy-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "description", "unit-test-windows-autopilot-device-preparation-policy-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "role_scope_tag_ids.*", "0"),

					// Device security group
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "device_security_group", "00000000-0000-0000-0000-000000000001"),

					// Deployment settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "deployment_settings.deployment_mode", "enrollment_autopilot_dpp_deploymentmode_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "deployment_settings.deployment_type", "enrollment_autopilot_dpp_deploymenttype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "deployment_settings.join_type", "enrollment_autopilot_dpp_jointype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "deployment_settings.account_type", "enrollment_autopilot_dpp_accountype_1"),

					// OOBE settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "oobe_settings.timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "oobe_settings.custom_error_message", "Please contact your IT administrator for assistance with device setup."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "oobe_settings.allow_skip", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "oobe_settings.allow_diagnostics", "true"),

					// Allowed apps
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.0.app_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.0.app_type", "win32LobApp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.1.app_id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.1.app_type", "winGetApp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.2.app_id", "00000000-0000-0000-0000-000000000005"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_apps.2.app_type", "officeSuiteApp"),

					// Allowed scripts
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_scripts.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_scripts.*", "00000000-0000-0000-0000-000000000006"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "allowed_scripts.*", "00000000-0000-0000-0000-000000000007"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "assignments.include_group_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "assignments.include_group_ids.*", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "assignments.include_group_ids.*", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.maximal", "assignments.include_group_ids.*", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_03_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Group ID in include_group_ids"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_autopilot_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_autopilot_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

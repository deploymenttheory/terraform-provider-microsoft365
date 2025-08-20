package graphBetaWindowsDeviceCompliancePolicies_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	complianceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *complianceMocks.WindowsDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceMock := &complianceMocks.WindowsDeviceCompliancePolicyMock{}
	complianceMock.RegisterMocks()
	return mockClient, complianceMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *complianceMocks.WindowsDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceMock := &complianceMocks.WindowsDeviceCompliancePolicyMock{}
	complianceMock.RegisterErrorMocks()
	return mockClient, complianceMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsDeviceCompliancePolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWindowsMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "display_name", "Windows Compliance Policy Minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_block_simple", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_minimum_length", "8"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_required_type", "alphanumeric"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "os_minimum_version", "10.0.19041.0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "os_maximum_version", "10.0.22631.3155"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_AllTerraformConfigurations(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWindowsMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "display_name", "Windows Compliance Policy Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "password_block_simple", "true"),
				),
			},
			{
				Config: testConfigWindowsMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "display_name", "Windows Compliance Policy Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_threat_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_threat_protection_required_security_level", "medium"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "valid_operating_system_build_ranges.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "scheduled_actions_for_rule.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "assignments.#", "1"),
				),
			},
			{
				Config: testConfigWindowsWSL(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "display_name", "Windows Compliance Policy with WSL"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "wsl_distributions.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "wsl_distributions.*", map[string]string{
						"distribution":       "Ubuntu",
						"minimum_os_version": "20.04",
						"maximum_os_version": "22.04",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_device_compliance_policy.wsl", "wsl_distributions.*", map[string]string{
						"distribution":       "Debian",
						"minimum_os_version": "11.0",
						"maximum_os_version": "12.0",
					}),
				),
			},
			{
				Config: testConfigWindowsCustomCompliance(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "display_name", "Windows Compliance Policy with Custom Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "custom_compliance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "device_compliance_policy_script.device_compliance_script_id", "44444444-4444-4444-4444-444444444444"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.custom_compliance", "device_compliance_policy_script.rules_content", "ZWNobyAiSGVsbG8gV29ybGQiCg=="),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_SecuritySettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWindowsMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "early_launch_anti_malware_driver_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "bit_locker_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "secure_boot_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "code_integrity_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "memory_integrity_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "kernel_dma_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "virtualization_based_security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "firmware_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "storage_require_encryption", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "active_firewall_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "defender_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "signature_out_of_date", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "rtp_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "antivirus_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "anti_spyware_required", "true"),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_ScheduledActions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWindowsMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "scheduled_actions_for_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "scheduled_actions_for_rule.0.scheduled_action_configurations.*", map[string]string{
						"action_type":        "block",
						"grace_period_hours": "24",
					}),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWindowsMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigWindowsMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Device Compliance Policy data"),
			},
		},
	})
}

// Configuration Functions
func testConfigWindowsMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_minimal.tf")
	if err != nil {
		panic("failed to load windows minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWindowsMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_maximal.tf")
	if err != nil {
		panic("failed to load windows maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWindowsWSL() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_wsl.tf")
	if err != nil {
		panic("failed to load windows WSL config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWindowsCustomCompliance() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_custom_compliance.tf")
	if err != nil {
		panic("failed to load windows custom compliance config: " + err.Error())
	}
	return unitTestConfig
}

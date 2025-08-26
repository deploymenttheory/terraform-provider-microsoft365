package graphBetaWindowsDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsDeviceCompliancePolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsDeviceCompliancePolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsDeviceCompliancePolicyResource_Schema(t *testing.T) {
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "display_name", "unit-test-windows-device-compliance-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "description", "unit-test-windows-device-compliance-policy-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "role_scope_tag_ids.*", "0"),

					// Microsoft Defender for Endpoint settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "microsoft_defender_for_endpoint.device_threat_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "microsoft_defender_for_endpoint.device_threat_protection_required_security_level", "medium"),

					// Scheduled actions for rule
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.#", "3"),

					// Block action
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.0.action_type", "block"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.0.grace_period_hours", "12"),

					// Notification action
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.1.action_type", "notification"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.1.grace_period_hours", "24"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.1.notification_template_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.1.notification_message_cc_list.#", "2"),

					// Retire action
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.2.action_type", "retire"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "scheduled_actions_for_rule.0.scheduled_action_configurations.2.grace_period_hours", "48"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.minimal", "assignments.#", "6"),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_MaximalSettings(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "display_name", "unit-test-windows-device-compliance-policy-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "description", "unit-test-windows-device-compliance-policy-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "role_scope_tag_ids.*", "0"),

					// Device health settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_health.bit_locker_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_health.secure_boot_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_health.code_integrity_enabled", "true"),

					// Microsoft Defender for Endpoint settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "microsoft_defender_for_endpoint.device_threat_protection_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "microsoft_defender_for_endpoint.device_threat_protection_required_security_level", "medium"),

					// Device properties
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_properties.os_minimum_version", "10.0.22631.5768"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_properties.os_maximum_version", "10.0.26100.9999"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_properties.mobile_os_minimum_version", "10.0.22631.5768"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_properties.mobile_os_maximum_version", "10.0.26100.9999"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_properties.valid_operating_system_build_ranges.#", "2"),

					// Custom compliance
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "custom_compliance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_compliance_policy_script.device_compliance_script_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "device_compliance_policy_script.rules_content"),

					// Scheduled actions for rule
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "scheduled_actions_for_rule.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "scheduled_actions_for_rule.0.scheduled_action_configurations.#", "3"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_policy.maximal", "assignments.#", "6"),
				),
			},
		},
	})
}

func TestWindowsDeviceCompliancePolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Group ID in notification_message_cc_list"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_windows_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

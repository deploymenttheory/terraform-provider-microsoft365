package graphBetaWindowsDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
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

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Scenario 1 - Minimal configuration without assignments
func TestWindowsDeviceCompliancePolicyResource_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_windows_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-wdcp-minimal"),
					check.That(resourceType+".minimal").Key("description").HasValue("unit-test-wdcp-minimal"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),

					// Microsoft Defender for Endpoint settings
					check.That(resourceType+".minimal").Key("microsoft_defender_for_endpoint.device_threat_protection_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("microsoft_defender_for_endpoint.device_threat_protection_required_security_level").HasValue("medium"),

					// Scheduled actions for rule
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.#").HasValue("3"),

					// Block action
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.action_type").HasValue("block"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.grace_period_hours").HasValue("12"),

					// Notification action
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.1.action_type").HasValue("notification"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.1.grace_period_hours").HasValue("24"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.1.notification_template_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.1.notification_message_cc_list.#").HasValue("2"),

					// Retire action
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.2.action_type").HasValue("retire"),
					check.That(resourceType+".minimal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.2.grace_period_hours").HasValue("48"),

					// Assignments
					check.That(resourceType+".minimal").Key("assignments.#").HasValue("6"),
				),
			},
		},
	})
}

// Test 02: Scenario 2 - Maximal configuration with assignments
func TestWindowsDeviceCompliancePolicyResource_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_windows_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-wdcp-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("unit-test-wdcp-maximal"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),

					// Device health settings
					check.That(resourceType+".maximal").Key("device_health.bit_locker_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("device_health.secure_boot_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("device_health.code_integrity_enabled").HasValue("true"),

					// Microsoft Defender for Endpoint settings
					check.That(resourceType+".maximal").Key("microsoft_defender_for_endpoint.device_threat_protection_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("microsoft_defender_for_endpoint.device_threat_protection_required_security_level").HasValue("medium"),

					// Device properties
					check.That(resourceType+".maximal").Key("device_properties.os_minimum_version").HasValue("10.0.22631.5768"),
					check.That(resourceType+".maximal").Key("device_properties.os_maximum_version").HasValue("10.0.26100.9999"),
					check.That(resourceType+".maximal").Key("device_properties.mobile_os_minimum_version").HasValue("10.0.22631.5768"),
					check.That(resourceType+".maximal").Key("device_properties.mobile_os_maximum_version").HasValue("10.0.26100.9999"),
					check.That(resourceType+".maximal").Key("device_properties.valid_operating_system_build_ranges.#").HasValue("2"),

					// Custom compliance
					check.That(resourceType+".maximal").Key("custom_compliance_required").HasValue("true"),
					check.That(resourceType+".maximal").Key("device_compliance_policy_script.device_compliance_script_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".maximal").Key("device_compliance_policy_script.rules_content").Exists(),

					// Scheduled actions for rule
					check.That(resourceType+".maximal").Key("scheduled_actions_for_rule.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("scheduled_actions_for_rule.0.scheduled_action_configurations.#").HasValue("3"),

					// Assignments
					check.That(resourceType+".maximal").Key("assignments.#").HasValue("6"),
				),
			},
		},
	})
}

// Test 03: Scenario 3 - Error handling
func TestWindowsDeviceCompliancePolicyResource_03_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_windows_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid Group ID in notification_message_cc_list"),
			},
		},
	})
}

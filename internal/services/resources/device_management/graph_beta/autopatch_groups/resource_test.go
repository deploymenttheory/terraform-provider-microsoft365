package graphBetaAutopatchGroups_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAutopatchGroups "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/autopatch_groups"
	autopatchGroupsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/autopatch_groups/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAutopatchGroups.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *autopatchGroupsMocks.AutopatchGroupsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	autopatchGroupsMock := &autopatchGroupsMocks.AutopatchGroupsMock{}
	autopatchGroupsMock.RegisterMocks()
	return mockClient, autopatchGroupsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *autopatchGroupsMocks.AutopatchGroupsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	autopatchGroupsMock := &autopatchGroupsMocks.AutopatchGroupsMock{}
	autopatchGroupsMock.RegisterErrorMocks()
	return mockClient, autopatchGroupsMock
}

// Test Autopatch Group with Mixed Distribution
func TestUnitResourceAutopatchGroups_01_MixedDistribution(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, autopatchGroupsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer autopatchGroupsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMixedDistribution(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					check.That(resourceType+".unit-test_autopatch_group").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".unit-test_autopatch_group").Key("name").HasValue("unit-test-autopatch-group"),
					check.That(resourceType+".unit-test_autopatch_group").Key("description").HasValue("unit-test"),
					check.That(resourceType+".unit-test_autopatch_group").Key("tenant_id").HasValue("2fd6bb84-ad40-4ec5-9369-a215b25c9952"),
					check.That(resourceType+".unit-test_autopatch_group").Key("type").HasValue("User"),
					check.That(resourceType+".unit-test_autopatch_group").Key("status").HasValue("Active"),
					check.That(resourceType+".unit-test_autopatch_group").Key("distribution_type").HasValue("Mixed"),
					check.That(resourceType+".unit-test_autopatch_group").Key("is_locked_by_policy").HasValue("false"),
					check.That(resourceType+".unit-test_autopatch_group").Key("read_only").HasValue("false"),
					check.That(resourceType+".unit-test_autopatch_group").Key("flow_status").HasValue("Succeeded"),
					check.That(resourceType+".unit-test_autopatch_group").Key("number_of_registered_devices").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("user_has_all_scope_tag").HasValue("true"),

					// Global User Managed AAD Groups
					check.That(resourceType+".unit-test_autopatch_group").Key("global_user_managed_aad_groups.#").HasValue("2"),
					check.That(resourceType+".unit-test_autopatch_group").Key("global_user_managed_aad_groups.0.id").HasValue("550dba96-abd7-4ef0-9cf1-25be705f676c"),
					check.That(resourceType+".unit-test_autopatch_group").Key("global_user_managed_aad_groups.0.type").HasValue("None"),
					check.That(resourceType+".unit-test_autopatch_group").Key("global_user_managed_aad_groups.1.id").HasValue("6a08c3a0-1693-4089-80cb-f9c2f8063a3b"),
					check.That(resourceType+".unit-test_autopatch_group").Key("global_user_managed_aad_groups.1.type").HasValue("None"),

					// Scope Tags
					check.That(resourceType+".unit-test_autopatch_group").Key("scope_tags.#").HasValue("3"),
					check.That(resourceType+".unit-test_autopatch_group").Key("scope_tags.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("scope_tags.*").ContainsTypeSetElement("1232"),
					check.That(resourceType+".unit-test_autopatch_group").Key("scope_tags.*").ContainsTypeSetElement("1234"),

					// Deployment Groups - Count only (sets are unordered, so we only verify count)
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.#").HasValue("4"),

					// Deployment Group 0 - Policy Settings
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.aad_group_name").HasValue("unit-test-autopatch-group - test"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.is_update_settings_modified").HasValue("false"),

					// Deployment Group 0 - Device Configuration Setting
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.policy_id").HasValue("b97bff48-1138-4ad2-b2b2-9f51f98df847"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.update_behavior").HasValue("AutoInstallAndRestart"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.notification_setting").HasValue("DefaultNotifications"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.deferral").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.deadline").HasValue("1"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.grace_period").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.feature_deployment_settings.deferral").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.device_configuration_setting.feature_deployment_settings.deadline").HasValue("5"),

					// Deployment Group 0 - Feature Update Anchor Cloud Setting
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.feature_update_anchor_cloud_setting.target_os_version").HasValue("Windows 11, version 25H2"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.feature_update_anchor_cloud_setting.install_latest_windows10_on_windows11_ineligible_device").HasValue("true"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.feature_update_anchor_cloud_setting.policy_id").HasValue("297070c5-0748-4ec0-9452-b51d03fe1d3a"),

					// Deployment Group 0 - DNF Update Cloud Setting
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.dnf_update_cloud_setting.approval_type").HasValue("Automatic"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.dnf_update_cloud_setting.deployment_deferral_in_days").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.dnf_update_cloud_setting.policy_id").HasValue("6e638095-13ee-407d-8eb8-2c46e47461b4"),

					// Deployment Group 0 - Edge DCv2 Setting
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.edge_dcv2_setting.target_channel").HasValue("Beta"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.edge_dcv2_setting.policy_id").HasValue("1d6ff71e-38c6-45cc-bb86-f61dd19c4022"),

					// Deployment Group 0 - Office DCv2 Setting
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.target_channel").HasValue("MonthlyEnterprise"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.deferral").HasValue("0"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.deadline").HasValue("1"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.hide_update_notifications").HasValue("false"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.enable_automatic_update").HasValue("true"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.hide_enable_disable_update").HasValue("true"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.enable_office_mgmt").HasValue("false"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.update_path").HasValue("http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.0.deployment_group_policy_settings.office_dcv2_setting.policy_id").HasValue("bae5a6d7-12bf-48d0-96ad-08f462f5e7e1"),

					// Deployment Group 1 - Ring1
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.1.name").HasValue("unit-test - Ring1"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.1.aad_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.1.distribution").HasValue("75"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.1.deployment_group_policy_settings.dnf_update_cloud_setting.deployment_deferral_in_days").HasValue("1"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.1.deployment_group_policy_settings.edge_dcv2_setting.target_channel").HasValue("Stable"),

					// Deployment Group 2 - Ring2
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.2.name").HasValue("unit-test-autopatch-group - Ring2"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.2.aad_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.2.distribution").HasValue("25"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.2.deployment_group_policy_settings.dnf_update_cloud_setting.approval_type").HasValue("Manual"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.2.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.deferral").HasValue("5"),

					// Deployment Group 3 - Last
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.3.name").HasValue("unit-test-autopatch-group - Last"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.3.aad_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.3.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.deferral").HasValue("9"),
					check.That(resourceType+".unit-test_autopatch_group").Key("deployment_groups.3.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings.deadline").HasValue("5"),
				),
			},
			{
				ResourceName:      resourceType + ".unit-test_autopatch_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test Autopatch Group Create Error
func TestUnitResourceAutopatchGroups_02_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, autopatchGroupsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer autopatchGroupsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMixedDistribution(),
				ExpectError: regexp.MustCompile(`Bad Request`),
			},
		},
	})
}

func testConfigMixedDistribution() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_autopatch_group_test.tf")
	if err != nil {
		panic("failed to load autopatch group test config: " + err.Error())
	}
	return unitTestConfig
}

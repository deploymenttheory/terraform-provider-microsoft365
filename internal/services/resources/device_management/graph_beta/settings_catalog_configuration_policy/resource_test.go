package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUnitSettingsCatalogConfigurationPolicyResource(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Create and Read with minimal configuration
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.*", "0"),
				),
			},
			// Test Create and Read with maximal configuration
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "description", "Configure the FileVault payload to manage FileVault disk encryption settings on devices."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.#", "1"),
				),
			},
			// Test Create and Read with all assignment types including filters
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_allAssignmentTypes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "name", "Test All Assignment Types Settings Catalog Policy - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					// Verify assignment filters are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"filter_type": "include"}),
				),
			},
			// Test Create and Read with group assignments only
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_groupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "name", "Test Group Assignments Settings Catalog Policy - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.#", "2"),
					// Verify both group assignments with filters
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "include"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333", "filter_type": "include"}),
				),
			},
			// Test Create and Read with all devices assignment
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_allDevicesAssignment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "name", "Test All Devices Assignment Settings Catalog Policy - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget", "filter_type": "include"}),
				),
			},
			// Test Create and Read with all users assignment
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_allUsersAssignment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "name", "Test All Users Assignment Settings Catalog Policy - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget", "filter_type": "include"}),
				),
			},
			// Test Create and Read with exclusion assignment
			{
				Config: testUnitSettingsCatalogConfigurationPolicyResourceConfig_exclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "name", "Test Exclusion Assignment Settings Catalog Policy - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "7777777-7777-7777-7777-777777777777", "filter_type": "include"}),
				),
			},
		},
	})
}

// Test configuration functions
func testUnitSettingsCatalogConfigurationPolicyResourceConfig_minimal() string {
	config := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	if config == "" {
		panic("minimal config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_maximal() string {
	config := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	if config == "" {
		panic("maximal config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_allAssignmentTypes() string {
	config := mocks.LoadLocalTerraformConfig("resource_with_all_assignment_types.tf")
	if config == "" {
		panic("all assignment types config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_groupAssignments() string {
	config := mocks.LoadLocalTerraformConfig("resource_with_group_assignments.tf")
	if config == "" {
		panic("group assignments config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_allDevicesAssignment() string {
	config := mocks.LoadLocalTerraformConfig("resource_with_all_devices_assignment.tf")
	if config == "" {
		panic("all devices assignment config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_allUsersAssignment() string {
	config := mocks.LoadLocalTerraformConfig("resource_with_all_users_assignment.tf")
	if config == "" {
		panic("all users assignment config is empty")
	}
	return config
}

func testUnitSettingsCatalogConfigurationPolicyResourceConfig_exclusionAssignment() string {
	config := mocks.LoadLocalTerraformConfig("resource_with_exclusion_assignment.tf")
	if config == "" {
		panic("exclusion assignment config is empty")
	}
	return config
}


package graphBetaSettingsCatalogInventoryPolicy_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	inventoryPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_inventory_policy"
	inventoryPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_inventory_policy/mocks"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	terraformResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *inventoryPolicyMocks.SettingsCatalogInventoryPolicyMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	inventoryMock := &inventoryPolicyMocks.SettingsCatalogInventoryPolicyMock{}
	inventoryMock.RegisterMocks()

	return mockClient, inventoryMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *inventoryPolicyMocks.SettingsCatalogInventoryPolicyMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	inventoryMock := &inventoryPolicyMocks.SettingsCatalogInventoryPolicyMock{}
	inventoryMock.RegisterErrorMocks()

	return mockClient, inventoryMock
}

func TestUnitResourceSettingsCatalogInventoryPolicy_01_SchemaValidation(t *testing.T) {
	t.Run("resource schema validation", func(t *testing.T) {
		startTime := time.Now()

		resourceInstance := inventoryPolicy.NewInventoryPolicyResource()

		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		resourceInstance.Schema(context.Background(), req, resp)

		elapsed := time.Since(startTime)
		if elapsed > time.Second*30 {
			t.Fatalf("Schema construction took too long: %v", elapsed)
		}

		if resp.Schema.Attributes == nil {
			t.Fatal("Schema attributes should not be nil")
		}

		expectedAttrs := []string{"id", "name", "description", "platforms", "technologies", "configuration_policy", "assignments"}
		for _, attr := range expectedAttrs {
			if _, exists := resp.Schema.Attributes[attr]; !exists {
				t.Fatalf("Resource attribute %s should exist", attr)
			}
		}

		// Verify template_reference does NOT exist (key difference from configuration_policy)
		if _, exists := resp.Schema.Attributes["template_reference"]; exists {
			t.Fatal("template_reference attribute should NOT exist on inventory policy")
		}

		// Verify is_assigned does NOT exist
		if _, exists := resp.Schema.Attributes["is_assigned"]; exists {
			t.Fatal("is_assigned attribute should NOT exist on inventory policy")
		}

		// Verify technologies is a StringAttribute (not ListAttribute)
		techAttr, exists := resp.Schema.Attributes["technologies"]
		if !exists {
			t.Fatal("technologies attribute should exist")
		}
		if _, ok := techAttr.(schema.StringAttribute); !ok {
			t.Fatal("technologies should be a StringAttribute (not ListAttribute)")
		}

		// Verify configuration_policy is correctly structured
		configPolicyAttr, exists := resp.Schema.Attributes["configuration_policy"]
		if !exists {
			t.Fatal("configuration_policy attribute should exist")
		}

		singleNestedAttr, ok := configPolicyAttr.(schema.SingleNestedAttribute)
		if !ok {
			t.Fatal("configuration_policy should be a SingleNestedAttribute")
		}

		if _, exists := singleNestedAttr.Attributes["settings"]; !exists {
			t.Fatal("settings attribute should exist within configuration_policy")
		}

		t.Logf("Resource schema validation passed in %v", elapsed)
	})

	t.Run("platform and technology validation", func(t *testing.T) {
		resourceInstance := inventoryPolicy.NewInventoryPolicyResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		resourceInstance.Schema(context.Background(), req, resp)

		if platformAttr, exists := resp.Schema.Attributes["platforms"]; exists {
			stringAttr, ok := platformAttr.(schema.StringAttribute)
			if !ok {
				t.Fatal("platforms should be a StringAttribute")
			}
			if !stringAttr.Optional {
				t.Fatal("platforms should be optional")
			}
			if !stringAttr.Computed {
				t.Fatal("platforms should be computed")
			}
		} else {
			t.Fatal("platforms attribute should exist")
		}

		if techAttr, exists := resp.Schema.Attributes["technologies"]; exists {
			stringAttr, ok := techAttr.(schema.StringAttribute)
			if !ok {
				t.Fatal("technologies should be a StringAttribute")
			}
			if !stringAttr.Optional {
				t.Fatal("technologies should be optional")
			}
			if !stringAttr.Computed {
				t.Fatal("technologies should be computed")
			}
		} else {
			t.Fatal("technologies attribute should exist")
		}

		t.Log("Platform and technology validation passed")
	})
}

func TestUnitResourceSettingsCatalogInventoryPolicy_02_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "test" {
  platforms    = "windows10"
  technologies = "extensibility"
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`Missing required argument|name`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "test" {
  name         = "Test Policy"
  platforms    = "invalid_platform"
  technologies = "extensibility"
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute platforms value must be one of|invalid_platform`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "test" {
  name         = "Test Policy"
  platforms    = "windows10"
  technologies = "invalid_technology"
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute technologies value must be one of|invalid_technology`),
			},
		},
	})
}

func TestUnitResourceSettingsCatalogInventoryPolicy_03_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: testUnitInventoryPolicyMinimal(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.minimal", "name", "Test Inventory Policy Minimal - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.minimal", "platforms", "windows10"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.minimal", "technologies", "extensibility"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.minimal", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
				),
			},
			{
				Config: testUnitInventoryPolicyGroupCollection(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_collection", "name", "Test Inventory Policy Group Collection - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_collection", "description", "Test inventory policy with full application inventory settings"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.#", "1"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.#", "4"),
				),
			},
		},
	})
}

func TestUnitResourceSettingsCatalogInventoryPolicy_04_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: testUnitInventoryPolicyWithAllDevicesAssignment(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_devices_assignment", "name", "Test All Devices Assignment Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_devices_assignment", "platforms", "windows10"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_devices_assignment", "technologies", "extensibility"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_devices_assignment", "assignments.#", "1"),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_devices_assignment", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
				),
			},
			{
				Config: testUnitInventoryPolicyWithAllUsersAssignment(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_users_assignment", "name", "Test All Users Assignment Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_users_assignment", "assignments.#", "1"),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_users_assignment", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
				),
			},
			{
				Config: testUnitInventoryPolicyWithGroupAssignments(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_assignments", "name", "Test Group Assignments Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_assignments", "assignments.#", "2"),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_assignments", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.group_assignments", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
			{
				Config: testUnitInventoryPolicyWithExclusionAssignment(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.exclusion_assignment", "name", "Test Exclusion Assignment Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.exclusion_assignment", "assignments.#", "1"),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.exclusion_assignment", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				Config: testUnitInventoryPolicyWithAllAssignmentTypes(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "name", "Test All Assignment Types Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "description", "Inventory policy with all assignment types for unit testing"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "assignments.#", "5"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "role_scope_tag_ids.#", "2"),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					terraformResource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.all_assignment_types", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
				),
			},
		},
	})
}

func TestUnitResourceSettingsCatalogInventoryPolicy_05_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, inventoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer inventoryMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			{
				Config: testUnitInventoryPolicyMaximal(),
				Check: terraformResource.ComposeTestCheckFunc(
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "name", "Test Maximal Inventory Policy - Unit"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "description", "Comprehensive inventory policy with full settings and all assignment types"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "platforms", "windows10"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "technologies", "extensibility"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "role_scope_tag_ids.#", "2"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.#", "1"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.#", "6"),
					terraformResource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_inventory_policy.maximal", "assignments.#", "5"),
				),
			},
		},
	})
}

func testUnitInventoryPolicyMinimal() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyGroupCollection() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_collection.tf")
	if err != nil {
		panic("failed to load group collection config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyWithAllDevicesAssignment() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_devices_assignment.tf")
	if err != nil {
		panic("failed to load all devices assignment config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyWithAllUsersAssignment() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_users_assignment.tf")
	if err != nil {
		panic("failed to load all users assignment config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyWithGroupAssignments() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_group_assignments.tf")
	if err != nil {
		panic("failed to load group assignments config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyWithExclusionAssignment() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_exclusion_assignment.tf")
	if err != nil {
		panic("failed to load exclusion assignment config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyWithAllAssignmentTypes() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_assignment_types.tf")
	if err != nil {
		panic("failed to load all assignment types config: " + err.Error())
	}
	return config
}

func testUnitInventoryPolicyMaximal() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return config
}

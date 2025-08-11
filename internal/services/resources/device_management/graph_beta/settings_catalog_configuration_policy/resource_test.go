package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	settingsCatalogMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *settingsCatalogMocks.SettingsCatalogConfigurationPolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	settingsCatalogMock := &settingsCatalogMocks.SettingsCatalogConfigurationPolicyMock{}
	settingsCatalogMock.RegisterMocks()

	return mockClient, settingsCatalogMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *settingsCatalogMocks.SettingsCatalogConfigurationPolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	settingsCatalogMock := &settingsCatalogMocks.SettingsCatalogConfigurationPolicyMock{}
	settingsCatalogMock.RegisterErrorMocks()

	return mockClient, settingsCatalogMock
}

func TestUnitSettingsCatalogConfigurationPolicyResource(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
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

// Unit tests for all setting types covering construct_configuration_policy_settings.go functionality
func TestUnitConstructSettingsCatalogSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Simple Setting - String value from FileVault location setting
			{
				Config: testUnitSettingsCatalogSimpleString(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "name", "Test Simple String Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "configuration_policy.settings.0.setting_instance.simple_setting_value.odata_type", "#microsoft.graph.deviceManagementConfigurationStringSettingValue"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "configuration_policy.settings.0.setting_instance.simple_setting_value.value", "Personal recovery key location message"),
				),
			},
			// Test Simple Setting - Secret value (password)
			{
				Config: testUnitSettingsCatalogSimpleSecret(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "name", "Test Simple Secret Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "configuration_policy.settings.0.setting_instance.simple_setting_value.odata_type", "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "configuration_policy.settings.0.setting_instance.simple_setting_value.value"),
				),
			},
			// Test Choice Setting - From Edge security settings
			{
				Config: testUnitSettingsCatalogChoice(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "name", "Test Choice Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "configuration_policy.settings.0.setting_instance.choice_setting_value.value", "com.apple.managedclient.preferences_smartscreenenabled_true"),
				),
			},
			// Test Simple Collection Setting - From Edge extensions
			{
				Config: testUnitSettingsCatalogSimpleCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "name", "Test Simple Collection Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.simple_setting_collection_value.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.simple_setting_collection_value.0.value", "extension_id_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.simple_setting_collection_value.1.value", "extension_id_2"),
				),
			},
			// Test Choice Collection Setting - Multiple choice values
			{
				Config: testUnitSettingsCatalogChoiceCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "name", "Test Choice Collection Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "configuration_policy.settings.0.setting_instance.choice_setting_collection_value.#", "2"),
				),
			},
			// Test Group Collection Setting - From FileVault configuration
			{
				Config: testUnitSettingsCatalogGroupCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "name", "Test Group Collection Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.#", "3"),
				),
			},
			// Test Complex Group Collection with Nested Simple Collection - From System Preferences
			{
				Config: testUnitSettingsCatalogComplexGroupCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "name", "Test Complex Group Collection Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.0.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"),
				),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_ErrorHandling tests error scenarios
func TestSettingsCatalogConfigurationPolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid configuration - missing required name field
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  platforms = "macOS"
  technologies = ["mdm"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`Missing required argument|name`),
			},
			// Test invalid platforms value
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Policy"
  platforms = "invalid_platform"
  technologies = ["mdm"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute platforms value must be one of|invalid_platform`),
			},
			// Test invalid technologies value
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Policy"
  platforms = "macOS"
  technologies = ["invalid_technology"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = []
  }
}
`,
				ExpectError: regexp.MustCompile(`invalid value for technologies|invalid_technology`),
			},
			// Test server error during creation (BadRequest from error mock)
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Error Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "test.error.setting"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value = "error_value"
          }
        }
        id = "0"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`Bad Request - 400|Invalid request body|BadRequest`),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_SettingTypeErrors tests specific setting type error scenarios
func TestSettingsCatalogConfigurationPolicyResource_SettingTypeErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid choice setting value
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Invalid Choice Value Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id = "test.choice.setting"
          choice_setting_value = {
            children = []
            value = "" # Empty/invalid choice value
          }
        }
        id = "0"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`Bad Request - 400|Invalid request body|BadRequest|empty.*value`),
			},
			// Test secret setting with invalid value_state
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Invalid Secret State Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "test.secret.setting"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            value = "secret_value"
            value_state = "invalidState"
          }
        }
        id = "0"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|value must be one of|invalidState`),
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

// TestSettingsCatalogConfigurationPolicyResource_Schema validates the resource schema
func TestSettingsCatalogConfigurationPolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Simple String Setting Schema
			{
				Config: testUnitSettingsCatalogSimpleString(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "name", "Test Simple String Setting - Unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "technologies.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "technologies.*", "mdm"),
					// Check simple string setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_string", "configuration_policy.settings.0.setting_instance.simple_setting_value.odata_type", "#microsoft.graph.deviceManagementConfigurationStringSettingValue"),
				),
			},
			// Test Simple Secret Setting Schema  
			{
				Config: testUnitSettingsCatalogSimpleSecret(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "name", "Test Simple Secret Setting - Unit"),
					// Check simple secret setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_secret", "configuration_policy.settings.0.setting_instance.simple_setting_value.odata_type", "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"),
				),
			},
			// Test Choice Setting Schema
			{
				Config: testUnitSettingsCatalogChoice(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "name", "Test Choice Setting - Unit"),
					// Check choice setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice", "configuration_policy.settings.0.setting_instance.choice_setting_value.value", "com.apple.managedclient.preferences_smartscreenenabled_true"),
				),
			},
			// Test Simple Collection Setting Schema
			{
				Config: testUnitSettingsCatalogSimpleCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "name", "Test Simple Collection Setting - Unit"),
					// Check simple collection setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.simple_setting_collection_value.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.simple_collection", "configuration_policy.settings.0.setting_instance.simple_setting_collection_value.0.odata_type", "#microsoft.graph.deviceManagementConfigurationStringSettingValue"),
				),
			},
			// Test Choice Collection Setting Schema
			{
				Config: testUnitSettingsCatalogChoiceCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "name", "Test Choice Collection Setting - Unit"),
					// Check choice collection setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.choice_collection", "configuration_policy.settings.0.setting_instance.choice_setting_collection_value.#", "2"),
				),
			},
			// Test Group Collection Setting Schema
			{
				Config: testUnitSettingsCatalogGroupCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "name", "Test Group Collection Setting - Unit"),
					// Check group collection setting structure
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.#", "3"),
				),
			},
			// Test Complex Group Collection Setting Schema
			{
				Config: testUnitSettingsCatalogComplexGroupCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "name", "Test Complex Group Collection Setting - Unit"),
					// Check complex group collection with nested simple collection
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "configuration_policy.settings.0.setting_instance.odata_type", "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.0.odata_type", "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.complex_group_collection", "configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.0.simple_setting_collection_value.#", "3"),
				),
			},
		},
	})
}

// Test configuration functions for different setting types
func testUnitSettingsCatalogSimpleString() string {
	config := mocks.LoadLocalTerraformConfig("resource_simple_string.tf")
	if config == "" {
		panic("simple string config is empty")
	}
	return config
}

func testUnitSettingsCatalogSimpleSecret() string {
	config := mocks.LoadLocalTerraformConfig("resource_simple_secret.tf")
	if config == "" {
		panic("simple secret config is empty")
	}
	return config
}

func testUnitSettingsCatalogChoice() string {
	config := mocks.LoadLocalTerraformConfig("resource_choice.tf")
	if config == "" {
		panic("choice config is empty")
	}
	return config
}

func testUnitSettingsCatalogSimpleCollection() string {
	config := mocks.LoadLocalTerraformConfig("resource_simple_collection.tf")
	if config == "" {
		panic("simple collection config is empty")
	}
	return config
}

func testUnitSettingsCatalogChoiceCollection() string {
	config := mocks.LoadLocalTerraformConfig("resource_choice_collection.tf")
	if config == "" {
		panic("choice collection config is empty")
	}
	return config
}

func testUnitSettingsCatalogGroupCollection() string {
	config := mocks.LoadLocalTerraformConfig("resource_group_collection.tf")
	if config == "" {
		panic("group collection config is empty")
	}
	return config
}

func testUnitSettingsCatalogComplexGroupCollection() string {
	config := mocks.LoadLocalTerraformConfig("resource_complex_group_collection.tf")
	if config == "" {
		panic("complex group collection config is empty")
	}
	return config
}

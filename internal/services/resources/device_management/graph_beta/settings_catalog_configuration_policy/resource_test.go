package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	settingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *settingsMocks.SettingsCatalogConfigurationPolicyMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	settingsMock := &settingsMocks.SettingsCatalogConfigurationPolicyMock{}
	settingsMock.RegisterMocks()

	return mockClient, settingsMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *settingsMocks.SettingsCatalogConfigurationPolicyMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	settingsMock := &settingsMocks.SettingsCatalogConfigurationPolicyMock{}
	settingsMock.RegisterErrorMocks()

	return mockClient, settingsMock
}

// testCheckExists ensures the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigGroupAssignments returns the group assignments configuration for testing
func testConfigGroupAssignments() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_group_assignments.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllUsersAssignment returns the all users assignment configuration for testing
func testConfigAllUsersAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_users_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllDevicesAssignment returns the all devices assignment configuration for testing
func testConfigAllDevicesAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_devices_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigExclusionAssignment returns the exclusion assignment configuration for testing
func testConfigExclusionAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_exclusion_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllAssignmentTypes returns the all assignment types configuration for testing
func testConfigAllAssignmentTypes() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_assignment_types.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestSettingsCatalogConfigurationPolicyResource_Schema validates the resource schema
func TestSettingsCatalogConfigurationPolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "name", "Test Minimal Settings Catalog Policy - Unique"),
					// Computed attributes
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_Minimal tests basic CRUD operations
func TestSettingsCatalogConfigurationPolicyResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "name", "Test Minimal Settings Catalog Policy - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "name", "Test Maximal Settings Catalog Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "description", "Maximal settings catalog policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_UpdateInPlace tests in-place updates
func TestSettingsCatalogConfigurationPolicyResource_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.minimal", "name", "Test Minimal Settings Catalog Policy - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "name", "Test Maximal Settings Catalog Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "description", "Maximal settings catalog policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_RequiredFields tests required field validation
func TestSettingsCatalogConfigurationPolicyResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  # Missing name
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "test.setting"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "value"
          }
        }
        id = "0"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyResource_ErrorHandling tests error scenarios
func TestSettingsCatalogConfigurationPolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Settings Catalog Policy"
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "test.setting"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "value"
          }
        }
        id = "0"
      }
    ]
  }
}

// TestSettingsCatalogConfigurationPolicyResource_GroupAssignments tests group assignment functionality
func TestSettingsCatalogConfigurationPolicyResource_GroupAssignments(t *testing.T) {
    mocks.SetupUnitTestEnvironment(t)
    _, settingsMock := setupMockEnvironment()
    defer httpmock.DeactivateAndReset()
    defer settingsMock.CleanupMockState()

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testConfigGroupAssignments(),
                Check: resource.ComposeTestCheckFunc(
                    testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "name", "Test Group Assignments Settings Catalog Policy - Unique"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.#", "2"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
                ),
            },
        },
    })
}

// TestSettingsCatalogConfigurationPolicyResource_AllUsersAssignment tests all licensed users assignment functionality
func TestSettingsCatalogConfigurationPolicyResource_AllUsersAssignment(t *testing.T) {
    mocks.SetupUnitTestEnvironment(t)
    _, settingsMock := setupMockEnvironment()
    defer httpmock.DeactivateAndReset()
    defer settingsMock.CleanupMockState()

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testConfigAllUsersAssignment(),
                Check: resource.ComposeTestCheckFunc(
                    testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "name", "Test All Users Assignment Settings Catalog Policy - Unique"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "assignments.#", "1"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_users_assignment", "assignments.0.type", "allLicensedUsersAssignmentTarget"),
                ),
            },
        },
    })
}

// TestSettingsCatalogConfigurationPolicyResource_AllDevicesAssignment tests all devices assignment functionality
func TestSettingsCatalogConfigurationPolicyResource_AllDevicesAssignment(t *testing.T) {
    mocks.SetupUnitTestEnvironment(t)
    _, settingsMock := setupMockEnvironment()
    defer httpmock.DeactivateAndReset()
    defer settingsMock.CleanupMockState()

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testConfigAllDevicesAssignment(),
                Check: resource.ComposeTestCheckFunc(
                    testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "name", "Test All Devices Assignment Settings Catalog Policy - Unique"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "assignments.#", "1"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_devices_assignment", "assignments.0.type", "allDevicesAssignmentTarget"),
                ),
            },
        },
    })
}

// TestSettingsCatalogConfigurationPolicyResource_ExclusionAssignment tests exclusion group assignment functionality
func TestSettingsCatalogConfigurationPolicyResource_ExclusionAssignment(t *testing.T) {
    mocks.SetupUnitTestEnvironment(t)
    _, settingsMock := setupMockEnvironment()
    defer httpmock.DeactivateAndReset()
    defer settingsMock.CleanupMockState()

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testConfigExclusionAssignment(),
                Check: resource.ComposeTestCheckFunc(
                    testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "name", "Test Exclusion Assignment Settings Catalog Policy - Unique"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "assignments.#", "1"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
                ),
            },
        },
    })
}

// TestSettingsCatalogConfigurationPolicyResource_AllAssignmentTypes tests all assignment types together
func TestSettingsCatalogConfigurationPolicyResource_AllAssignmentTypes(t *testing.T) {
    mocks.SetupUnitTestEnvironment(t)
    _, settingsMock := setupMockEnvironment()
    defer httpmock.DeactivateAndReset()
    defer settingsMock.CleanupMockState()

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testConfigAllAssignmentTypes(),
                Check: resource.ComposeTestCheckFunc(
                    testCheckExists("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "name", "Test All Assignment Types Settings Catalog Policy - Unique"),
                    resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.#", "5"),
                    resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
                    resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
                    resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
                    resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.all_assignment_types", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
                ),
            },
        },
    })
}
`,
				ExpectError: regexp.MustCompile(`Invalid settings catalog policy data|BadRequest`),
			},
		},
	})
}

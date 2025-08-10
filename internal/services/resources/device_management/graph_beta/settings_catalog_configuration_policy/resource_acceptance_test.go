package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSettingsCatalogConfigurationPolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore secret password field as Microsoft Graph returns different UUIDs for security
				ImportStateVerifyIgnore: []string{
					"configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.6.simple_setting_value.value",
				},
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "name", "Test Acceptance Settings Catalog Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			// Create with all assignment types including filters
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "name", "Test All Assignment Types Settings Catalog Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					// Verify assignment filters are included
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"filter_type": "include"}),
					// Verify role scope tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy:             testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy:             testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform(),
				ExpectError: regexp.MustCompile("Attribute platforms value must be one of"),
			},
		},
	})
}


// testAccCheckSettingsCatalogConfigurationPolicyDestroy verifies that settings catalog configuration policies have been destroyed
func testAccCheckSettingsCatalogConfigurationPolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" {
			continue
		}

		// Attempt to get the settings catalog configuration policy by ID
		_, err := graphClient.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			
			if errorInfo.StatusCode == 404 ||
				errorInfo.StatusCode == 400 ||  // Microsoft Graph sometimes returns 400 for deleted resources
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if settings catalog configuration policy %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("settings catalog configuration policy %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccSettingsCatalogConfigurationPolicyConfig_minimal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_maximal() string {
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_assignments() string {
	groups := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	assignmentFilters := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingName() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  platforms = "windows10"
  configuration_policy = {
    name = "Test Policy Settings"
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Policy"
  configuration_policy = {
    name = "Test Policy Settings"
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name      = "Test Policy"
  platforms = "invalid"
  configuration_policy = {
    name = "Test Policy Settings"
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}
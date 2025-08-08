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
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "name", "Test Acceptance Settings Catalog Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.*", "0"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "name", "Test Acceptance Settings Catalog Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.#", "2"),
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
		CheckDestroy: testAccCheckSettingsCatalogConfigurationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccSettingsCatalogConfigurationPolicyResource_WithAssignments(t *testing.T) {
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
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "name", "Test Settings Catalog Policy with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "assignments.#", "5"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test_assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

func testAccSettingsCatalogConfigurationPolicyConfig_minimal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_maximal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_withAssignments() string {
	groups := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	assignmentFilters := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_with_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingName() string {
	return `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  description = "desc"
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
`
}

// testAccCheckSettingsCatalogConfigurationPolicyDestroy verifies that settings catalog policies have been destroyed
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

		_, err := graphClient.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if settings catalog policy %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("settings catalog policy %s still exists", rs.Primary.ID)
	}

	return nil
}

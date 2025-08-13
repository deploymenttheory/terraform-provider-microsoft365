package graphBetaWindowsQualityUpdatePolicy_test

import (
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsQualityUpdatePolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsQualityUpdatePolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "display_name", "Acceptance - Windows Quality Update Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "role_scope_tag_ids.*", "0"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_quality_update_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWindowsQualityUpdatePolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "display_name", "Acceptance - Windows Quality Update Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test", "hotpatch_enabled", "true"),
				),
			},
		},
	})
}

func TestAccWindowsQualityUpdatePolicyResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsQualityUpdatePolicyConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_policy.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test_assignments", "display_name", "Acceptance - Windows Quality Update Policy with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.test_assignments", "assignments.1.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

func testAccWindowsQualityUpdatePolicyConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsQualityUpdatePolicyConfig_maximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsQualityUpdatePolicyConfig_withAssignments() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_with_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load assignments test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

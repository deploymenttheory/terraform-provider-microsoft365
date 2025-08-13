package graphBetaWindowsQualityUpdateExpeditePolicy_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsQualityUpdateExpeditePolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsQualityUpdateExpeditePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: ">= 3.7.2"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsQualityUpdateExpeditePolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "display_name", "Acceptance - Windows Quality Update Expedite Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "role_scope_tag_ids.*", "0"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWindowsQualityUpdateExpeditePolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "display_name", "Acceptance - Windows Quality Update Expedite Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test", "description", "Updated description for acceptance testing"),
				),
			},
		},
	})
}

func TestAccWindowsQualityUpdateExpeditePolicyResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsQualityUpdateExpeditePolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: ">= 3.7.2"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsQualityUpdateExpeditePolicyConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test_assignments", "display_name", "Acceptance - Windows Quality Update Expedite Policy with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test_assignments", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.test_assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

func testAccWindowsQualityUpdateExpeditePolicyConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsQualityUpdateExpeditePolicyConfig_maximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsQualityUpdateExpeditePolicyConfig_withAssignments() string {
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

func testAccCheckWindowsQualityUpdateExpeditePolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			WindowsQualityUpdateProfiles().
			ByWindowsQualityUpdateProfileId(rs.Primary.ID).
			Get(ctx, nil)
		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows quality update expedite policy %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows quality update expedite policy %s still exists", rs.Primary.ID)
	}
	return nil
}

package graphBetaWindowsDriverUpdateProfile_test

import (
	"log"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsDriverUpdateProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_driver_update_profile"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaWindowsDriverUpdateProfile.ResourceName

var testResource = graphBetaWindowsDriverUpdateProfile.WindowsDriverUpdateProfileTestResource{}

func TestAccResourceWindowsDriverUpdateProfile_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDriverUpdateProfileConfig_manual(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_manual", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_manual", "display_name", "Acceptance - Windows Driver Update Profile Manual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_manual", "approval_type", "manual"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_driver_update_profile.test_manual", ImportState: true, ImportStateVerify: true},
			{
				Config: testAccWindowsDriverUpdateProfileConfig_automatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_automatic", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_automatic", "display_name", "Acceptance - Windows Driver Update Profile Automatic"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_automatic", "approval_type", "automatic"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_automatic", "deployment_deferral_in_days", "7"),
				),
			},
		},
	})
}

func TestAccResourceWindowsDriverUpdateProfile_02_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDriverUpdateProfileConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_assignments", "display_name", "Acceptance - Windows Driver Update Profile with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_assignments", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_driver_update_profile.test_assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

func testAccWindowsDriverUpdateProfileConfig_withAssignments() string {
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

func testAccWindowsDriverUpdateProfileConfig_manual() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_manual.tf")
	if err != nil {
		log.Fatalf("Failed to load manual approval test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsDriverUpdateProfileConfig_automatic() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_automatic.tf")
	if err != nil {
		log.Fatalf("Failed to load automatic approval test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

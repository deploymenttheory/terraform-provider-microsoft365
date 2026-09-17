package graphBetaDeviceCategory_test

import (
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDeviceCategory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_category"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaDeviceCategory.ResourceName

var testResource = graphBetaDeviceCategory.DeviceCategoryTestResource{}

func TestAccResourceDeviceCategory_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccDeviceCategoryConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_category.test", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "display_name", regexp.MustCompile(`^Test Acceptance Device Category - [a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_device_category.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccDeviceCategoryConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_category.test", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "display_name", regexp.MustCompile(`^Test Acceptance Device Category - Updated - [a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.#", "3"),
				),
			},
		},
	})
}

func TestAccResourceDeviceCategory_02_RoleScopeTags(t *testing.T) {
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
				Config: testAccDeviceCategoryConfig_roleScopeTags(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_category.role_tags", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.role_tags", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.role_tags", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.role_tags", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.role_tags", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

func TestAccResourceDeviceCategory_03_Description(t *testing.T) {
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
				Config: testAccDeviceCategoryConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_category.description", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_device_category.description", "display_name", regexp.MustCompile(`^Test Description Device Category - [a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.description", "description", "This is a test device category with description"),
				),
			},
		},
	})
}

func testAccDeviceCategoryConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccDeviceCategoryConfig_maximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccDeviceCategoryConfig_roleScopeTags() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccDeviceCategoryConfig_description() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_description.tf")
	if err != nil {
		log.Fatalf("Failed to load description test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

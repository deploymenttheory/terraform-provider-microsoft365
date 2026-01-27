package graphBetaDeviceCategory_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

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
		CheckDestroy: testAccCheckDeviceCategoryDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccDeviceCategoryConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_category.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "display_name", "Test Acceptance Device Category"),
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "display_name", "Test Acceptance Device Category - Updated"),
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
		CheckDestroy:             testAccCheckDeviceCategoryDestroy,
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
		CheckDestroy:             testAccCheckDeviceCategoryDestroy,
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.description", "display_name", "Test Description Device Category"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.description", "description", "This is a test device category with description"),
				),
			},
		},
	})
}

// testAccCheckDeviceCategoryDestroy verifies that device categories have been destroyed
func testAccCheckDeviceCategoryDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_device_category" {
			continue
		}

		// Attempt to get the device category by ID
		_, err := graphClient.
			DeviceManagement().
			DeviceCategories().
			ByDeviceCategoryId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if device category %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("device category %s still exists", rs.Primary.ID)
	}

	return nil
}

// Test configuration functions

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

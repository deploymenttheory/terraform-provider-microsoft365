package graphBetaDeviceCategory_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccPreCheck verifies necessary test prerequisites
func testAccPreCheck(t *testing.T) {
	// Check for required environment variables
	requiredEnvVars := []string{
		"M365_CLIENT_ID",
		"M365_CLIENT_SECRET",
		"M365_TENANT_ID",
		"M365_AUTH_METHOD",
		"M365_CLOUD",
	}

	for _, envVar := range requiredEnvVars {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
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
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if device category %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("device category %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccDeviceCategoryResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDeviceCategoryDestroy,
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

func TestAccDeviceCategoryResource_RoleScopeTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDeviceCategoryDestroy,
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

func TestAccDeviceCategoryResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDeviceCategoryDestroy,
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

// Test configuration functions

func testAccDeviceCategoryConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccDeviceCategoryConfig_maximal() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccDeviceCategoryConfig_roleScopeTags() string {
	config := mocks.LoadTerraformConfigFile("resource_role_scope_tags.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccDeviceCategoryConfig_description() string {
	config := mocks.LoadTerraformConfigFile("resource_description.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

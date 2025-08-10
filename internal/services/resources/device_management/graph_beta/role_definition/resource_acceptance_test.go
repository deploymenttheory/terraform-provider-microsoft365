package graphBetaRoleDefinition_test

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

// testAccCheckRoleDefinitionDestroy verifies that role definitions have been destroyed
func testAccCheckRoleDefinitionDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_role_definition" {
			continue
		}

		// Attempt to get the role definition by ID
		_, err := graphClient.
			DeviceManagement().
			RoleDefinitions().
			ByRoleDefinitionId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			fmt.Printf("DEBUG: Error details - StatusCode: %d, ErrorCode: %s, ErrorMessage: %s\n",
				errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if role definition %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("role definition %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRoleDefinitionResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccRoleDefinitionConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "display_name", "Test Acceptance Role Definition"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "description", ""),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.test", "is_built_in_role_definition"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.test", "is_built_in"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_permissions.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_definition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccRoleDefinitionConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "display_name", "Test Acceptance Role Definition - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.test", "role_permissions.#", "1"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDefinitionConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_definition.description", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.description", "display_name", "Test Description Role Definition"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_definition.description", "description", "This is a test role definition with description"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccRoleDefinitionConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleDefinitionConfig_maximal() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleDefinitionConfig_description() string {
	config := mocks.LoadTerraformConfigFile("resource_description.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

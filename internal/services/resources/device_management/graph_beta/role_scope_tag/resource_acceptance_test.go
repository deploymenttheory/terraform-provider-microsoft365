package graphBetaRoleScopeTag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRoleScopeTagResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleScopeTagDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccRoleScopeTagConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Acceptance Role Scope Tag"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "is_built_in", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_scope_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccRoleScopeTagConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Acceptance Role Scope Tag - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", "Updated description for acceptance testing"),
				),
			},
		},
	})
}

func TestAccRoleScopeTagResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleScopeTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleScopeTagConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.description", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.description", "display_name", "Test Description Role Scope Tag"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.description", "description", "This is a test role scope tag with description"),
				),
			},
		},
	})
}

func TestAccRoleScopeTagResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleScopeTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleScopeTagConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "assignments.#", "2"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccRoleScopeTagConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleScopeTagConfig_maximal() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleScopeTagConfig_description() string {
	config := mocks.LoadTerraformConfigFile("resource_description.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleScopeTagConfig_assignments() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

// testAccCheckRoleScopeTagDestroy verifies that role scope tags have been destroyed
func testAccCheckRoleScopeTagDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_role_scope_tag" {
			continue
		}

		// Attempt to get the role scope tag by ID
		_, err := graphClient.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(rs.Primary.ID).
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
			return fmt.Errorf("error checking if role scope tag %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("role scope tag %s still exists", rs.Primary.ID)
	}

	return nil
}

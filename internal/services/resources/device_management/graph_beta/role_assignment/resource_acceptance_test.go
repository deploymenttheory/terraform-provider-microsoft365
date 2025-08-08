package graphBetaRoleDefinitionAssignment_test

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

// testAccCheckRoleAssignmentDestroy verifies that role assignments have been destroyed
func testAccCheckRoleAssignmentDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_role_assignment" {
			continue
		}

		// Attempt to get the role assignment by ID
		_, err := graphClient.
			DeviceManagement().
			RoleAssignments().
			ByDeviceAndAppManagementRoleAssignmentId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if role assignment %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("role assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRoleAssignmentResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccRoleAssignmentConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.test", "display_name"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "role_definition_id", "0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "scope_configuration.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "scope_configuration.0.type", "AllLicensedUsers"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["microsoft365_graph_beta_device_management_role_assignment.test"]
					if !ok {
						return "", fmt.Errorf("not found: microsoft365_graph_beta_device_management_role_assignment.test")
					}
					id := rs.Primary.ID
					roleDefId := rs.Primary.Attributes["role_definition_id"]
					compositeId := fmt.Sprintf("%s/%s", id, roleDefId)
					fmt.Printf("DEBUG: ImportStateIdFunc - id: %s, roleDefId: %s, compositeId: %s\n", id, roleDefId, compositeId)
					return compositeId, nil
				},
			},
			// Update to maximal configuration
			{
				Config: testAccRoleAssignmentConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.test", "display_name"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "role_definition_id", "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "members.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.test", "scope_configuration.0.type", "AllDevices"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentResource_ResourceScopes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentConfig_resourceScopes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.resource_scopes", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.resource_scopes", "scope_configuration.0.type", "ResourceScopes"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.resource_scopes", "scope_configuration.0.resource_scopes.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.resource_scopes", "members.#", "2"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentResource_AllDevicesScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentConfig_allDevices(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.all_devices", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.all_devices", "scope_configuration.0.type", "AllDevices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.all_devices", "members.#", "2"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentResource_AllUsersScope(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentConfig_allUsers(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_assignment.all_users", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.all_users", "scope_configuration.0.type", "AllLicensedUsers"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_assignment.all_users", "members.#", "2"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccRoleAssignmentConfig_minimal() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleAssignmentConfig_maximal() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleAssignmentConfig_resourceScopes() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_resource_scopes.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleAssignmentConfig_allDevices() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_all_devices.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleAssignmentConfig_allUsers() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_all_users.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

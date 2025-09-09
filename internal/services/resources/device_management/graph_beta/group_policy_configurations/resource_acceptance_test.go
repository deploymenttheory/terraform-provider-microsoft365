package graphBetaGroupPolicyConfigurations_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupPolicyConfigurationResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckGroupPolicyConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccGroupPolicyConfigurationConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_group_policy_configuration.test", "display_name", func(value string) error {
						if matched, _ := regexp.MatchString(`^Test Acceptance Group Policy Configuration - [0-9a-f-]+$`, value); !matched {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "last_modified_date_time"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_group_policy_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccGroupPolicyConfigurationConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_group_policy_configuration.test", "display_name", func(value string) error {
						if matched, _ := regexp.MatchString(`^Test Acceptance Group Policy Configuration - Updated - [0-9a-f-]+$`, value); !matched {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccGroupPolicyConfigurationResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckGroupPolicyConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupPolicyConfigurationConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "display_name"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "last_modified_date_time"),
				),
			},
		},
	})
}

func TestAccGroupPolicyConfigurationResource_OptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckGroupPolicyConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupPolicyConfigurationConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_group_policy_configuration.test", "display_name"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "role_scope_tag_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_group_policy_configuration.test", "role_scope_tag_ids.*", "1"),
				),
			},
		},
	})
}

func testAccCheckGroupPolicyConfigurationDestroy(s *terraform.State) error {
	// Get a Graph client using the same configuration as acceptance tests
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_group_policy_configuration" {
			continue
		}

		// Attempt to get the group policy configuration by ID
		_, err := graphClient.
			DeviceManagement().
			GroupPolicyConfigurations().
			ByGroupPolicyConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue // Resource successfully destroyed
			}

			return fmt.Errorf("unexpected error checking for Group Policy Configuration %s: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("Group Policy Configuration %s still exists", rs.Primary.ID)
	}

	return nil
}

// testAccGroupPolicyConfigurationConfig_minimal returns the minimal configuration for acceptance testing
func testAccGroupPolicyConfigurationConfig_minimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "resource_minimal.tf"))
	if err != nil {
		return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test" {
  display_name = "Test Acceptance Group Policy Configuration - ${random_uuid.test.result}"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}`
	}
	return string(content)
}

// testAccGroupPolicyConfigurationConfig_maximal returns the maximal configuration for acceptance testing
func testAccGroupPolicyConfigurationConfig_maximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "resource_maximal.tf"))
	if err != nil {
		return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test" {
  display_name       = "Test Acceptance Group Policy Configuration - Updated - ${random_uuid.test.result}"
  description        = "Updated description for acceptance testing"
  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}`
	}
	return string(content)
}

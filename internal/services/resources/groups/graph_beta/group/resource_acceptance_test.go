package graphBetaGroup_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccGroupConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.test", "created_date_time"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.test", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-test-group-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.test", "mail_nickname", func(value string) error {
						if !regexp.MustCompile(`^acctestgroup[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("mail_nickname does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "security_enabled", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_groups_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroupResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccGroupConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.test", "created_date_time"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.test", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-test-group-updated-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.test", "mail_nickname", func(value string) error {
						if !regexp.MustCompile(`^acctestgroup[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("mail_nickname does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "visibility", "Private"),
				),
			},
		},
	})
}

func TestAccGroupResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccGroupConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccGroupConfig_missingMailNickname(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccGroupConfig_missingMailEnabled(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccGroupConfig_missingSecurityEnabled(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccGroupResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccGroupConfig_invalidVisibility(),
				ExpectError: regexp.MustCompile("Attribute visibility value must be one of"),
			},
			{
				Config:      testAccGroupConfig_invalidMailNickname(),
				ExpectError: regexp.MustCompile("String contains forbidden character"),
			},
		},
	})
}

// TestAccGroupResource_Scenario1_SecurityGroupAssigned tests security group with assigned membership
func TestAccGroupResource_Scenario1_SecurityGroupAssigned(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_1", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_1", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-security-group-assigned-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_1", "mail_nickname", func(value string) error {
						if !regexp.MustCompile(`^accsg1[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("mail_nickname does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "description", "Acceptance test - Security group with assigned membership"),
				),
			},
		},
	})
}

// TestAccGroupResource_Scenario2_SecurityGroupDynamicUser tests security group with dynamic user membership
func TestAccGroupResource_Scenario2_SecurityGroupDynamicUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_2", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_2", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-security-group-dynamic-user-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "membership_rule", "(user.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "membership_rule_processing_state", "On"),
				),
			},
		},
	})
}

// TestAccGroupResource_Scenario3_SecurityGroupDynamicDevice tests security group with dynamic device membership
func TestAccGroupResource_Scenario3_SecurityGroupDynamicDevice(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario3(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_3", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_3", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-security-group-dynamic-device-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "membership_rule", "(device.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "membership_rule_processing_state", "On"),
				),
			},
		},
	})
}

// TestAccGroupResource_Scenario4_SecurityGroupRoleAssignable tests security group with Entra role assignment capability
func TestAccGroupResource_Scenario4_SecurityGroupRoleAssignable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_4", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_4", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-security-group-role-assignable-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "is_assignable_to_role", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "visibility", "Private"),
				),
			},
		},
	})
}

// TestAccGroupResource_Scenario5_M365GroupDynamicUser tests M365 group with dynamic user membership
func TestAccGroupResource_Scenario5_M365GroupDynamicUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario5(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_5", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_5", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-m365-group-dynamic-user-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "group_types.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "membership_rule", "(user.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "membership_rule_processing_state", "On"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "visibility", "Private"),
				),
			},
		},
	})
}

// TestAccGroupResource_Scenario6_M365GroupAssigned tests M365 group with assigned membership
func TestAccGroupResource_Scenario6_M365GroupAssigned(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_scenario6(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_groups_group.scenario_6", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_groups_group.scenario_6", "display_name", func(value string) error {
						if !regexp.MustCompile(`^acc-m365-group-assigned-[a-zA-Z0-9]{8}$`).MatchString(value) {
							return fmt.Errorf("display_name does not match expected pattern: %s", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "description", "Acceptance test - M365 group with assigned membership"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "is_assignable_to_role", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "visibility", "Private"),
				),
			},
		},
	})
}

func testAccGroupConfig_minimal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_maximal() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario1() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_1_security_group_assigned.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario2() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_2_security_group_dynamic_user.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario3() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_3_security_group_dynamic_device.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario4() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_4_security_group_role_assignable.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario5() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_5_m365_group_dynamic_user.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_scenario6() string {
	accTestConfig := mocks.LoadLocalTerraformConfig("resource_scenario_6_m365_group_assigned.tf")
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccGroupConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  mail_nickname    = "testgroup"
  mail_enabled     = false
  security_enabled = true
}
`
}

func testAccGroupConfig_missingMailNickname() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_enabled     = false
  security_enabled = true
}
`
}

func testAccGroupConfig_missingMailEnabled() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name  = "Test Group"
  mail_nickname = "testgroup"
  security_enabled = true
}
`
}

func testAccGroupConfig_missingSecurityEnabled() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name  = "Test Group"
  mail_nickname = "testgroup"
  mail_enabled  = false
}
`
}

func testAccGroupConfig_invalidVisibility() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_nickname    = "testgroup"
  mail_enabled     = false
  security_enabled = true
  visibility       = "invalid"
}
`
}

func testAccGroupConfig_invalidMailNickname() string {
	return `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_nickname    = "test@group"
  mail_enabled     = false
  security_enabled = true
}
`
}

// testAccCheckGroupDestroy verifies that groups have been destroyed
func testAccCheckGroupDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_groups_group" {
			continue
		}

		// Attempt to get the group by ID
		_, err := graphClient.
			Groups().
			ByGroupId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			fmt.Printf("DEBUG: Error details - StatusCode: %d, ErrorCode: %s, ErrorMessage: %s\n",
				errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage)

			if errorInfo.StatusCode == 404 ||
				errorInfo.StatusCode == 400 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "Request_ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/400/NotFound)\n", rs.Primary.ID)
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if group %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("group %s still exists", rs.Primary.ID)
	}

	return nil
}

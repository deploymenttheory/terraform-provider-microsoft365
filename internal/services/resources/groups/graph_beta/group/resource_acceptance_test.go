package graphBetaGroup_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// testResource is the test resource implementation for groups
	testResource = graphBetaGroup.GroupTestResource{}
)

func TestAccGroupResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal group")
				},
				Config: testAccGroupConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-group-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".test").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^acctestgroup[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".test").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("security_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal group")
				},
				ResourceName: resourceType + ".test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroupResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal group")
				},
				Config: testAccGroupConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-group-updated-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".test").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^acctestgroup[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That(resourceType+".test").Key("visibility").HasValue("Private"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal group")
				},
				ResourceName: resourceType + ".test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroupResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
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
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
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
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating security group with assigned membership")
				},
				Config: testAccGroupConfig_scenario1(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_1").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_1").Key("id").Exists(),
					check.That(resourceType+".scenario_1").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-security-group-assigned-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_1").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^accsg1[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_1").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_1").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_1").Key("description").HasValue("Acceptance test - Security group with assigned membership"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing security group with assigned membership")
				},
				ResourceName: resourceType + ".scenario_1",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_1"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_1")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupResource_Scenario2_SecurityGroupDynamicUser tests security group with dynamic user membership
func TestAccGroupResource_Scenario2_SecurityGroupDynamicUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating security group with dynamic user membership")
				},
				Config: testAccGroupConfig_scenario2(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_2").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_2").Key("id").Exists(),
					check.That(resourceType+".scenario_2").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-security-group-dynamic-user-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_2").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_2").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_2").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_2").Key("membership_rule").HasValue("(user.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_2").Key("membership_rule_processing_state").HasValue("On"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing security group with dynamic user membership")
				},
				ResourceName: resourceType + ".scenario_2",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_2"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_2")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupResource_Scenario3_SecurityGroupDynamicDevice tests security group with dynamic device membership
func TestAccGroupResource_Scenario3_SecurityGroupDynamicDevice(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating security group with dynamic device membership")
				},
				Config: testAccGroupConfig_scenario3(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_3").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_3").Key("id").Exists(),
					check.That(resourceType+".scenario_3").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-security-group-dynamic-device-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_3").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_3").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_3").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_3").Key("membership_rule").HasValue("(device.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_3").Key("membership_rule_processing_state").HasValue("On"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing security group with dynamic device membership")
				},
				ResourceName: resourceType + ".scenario_3",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_3"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_3")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupResource_Scenario4_SecurityGroupRoleAssignable tests security group with Entra role assignment capability
func TestAccGroupResource_Scenario4_SecurityGroupRoleAssignable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating security group with Entra role assignment capability")
				},
				Config: testAccGroupConfig_scenario4(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_4").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_4").Key("id").Exists(),
					check.That(resourceType+".scenario_4").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-security-group-role-assignable-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_4").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_4").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_4").Key("is_assignable_to_role").HasValue("true"),
					check.That(resourceType+".scenario_4").Key("visibility").HasValue("Private"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing security group with Entra role assignment capability")
				},
				ResourceName: resourceType + ".scenario_4",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_4"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_4")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupResource_Scenario5_M365GroupDynamicUser tests M365 group with dynamic user membership
func TestAccGroupResource_Scenario5_M365GroupDynamicUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating M365 group with dynamic user membership")
				},
				Config: testAccGroupConfig_scenario5(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_5").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_5").Key("id").Exists(),
					check.That(resourceType+".scenario_5").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-m365-group-dynamic-user-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_5").Key("description").HasValue("Acceptance test - M365 group with dynamic user membership"),
					check.That(resourceType+".scenario_5").Key("mail_enabled").HasValue("true"),
					check.That(resourceType+".scenario_5").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_5").Key("group_types.#").HasValue("2"),
					check.That(resourceType+".scenario_5").Key("membership_rule").HasValue("(user.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_5").Key("membership_rule_processing_state").HasValue("On"),
					check.That(resourceType+".scenario_5").Key("visibility").HasValue("Private"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing M365 group with dynamic user membership")
				},
				ResourceName: resourceType + ".scenario_5",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_5"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_5")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupResource_Scenario6_M365GroupAssigned tests M365 group with assigned membership
func TestAccGroupResource_Scenario6_M365GroupAssigned(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating M365 group with assigned membership")
				},
				Config: testAccGroupConfig_scenario6(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".scenario_6").ExistsInGraph(testResource),
					check.That(resourceType+".scenario_6").Key("id").Exists(),
					check.That(resourceType+".scenario_6").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-m365-group-assigned-[a-zA-Z0-9]{8}$`)),
					check.That(resourceType+".scenario_6").Key("mail_enabled").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_6").Key("description").HasValue("Acceptance test - M365 group with assigned membership"),
					check.That(resourceType+".scenario_6").Key("is_assignable_to_role").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("visibility").HasValue("Private"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing M365 group with assigned membership")
				},
				ResourceName: resourceType + ".scenario_6",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".scenario_6"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".scenario_6")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// Config loader functions
func testAccGroupConfig_minimal() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_maximal() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario1() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_1_security_group_assigned.tf")
	if err != nil {
		panic("failed to load scenario 1 config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario2() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_2_security_group_dynamic_user.tf")
	if err != nil {
		panic("failed to load scenario 2 config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario3() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_3_security_group_dynamic_device.tf")
	if err != nil {
		panic("failed to load scenario 3 config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario4() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_4_security_group_role_assignable.tf")
	if err != nil {
		panic("failed to load scenario 4 config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario5() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_5_m365_group_dynamic_user.tf")
	if err != nil {
		panic("failed to load scenario 5 config: " + err.Error())
	}
	return config
}

func testAccGroupConfig_scenario6() string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scenario_6_m365_group_assigned.tf")
	if err != nil {
		panic("failed to load scenario 6 config: " + err.Error())
	}
	return config
}

// Inline validation configs
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

package graphBetaGroup_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()
	return mockClient, groupMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterErrorMocks()
	return mockClient, groupMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestGroupResource_RequiredFields tests required field validation
func TestUnitResourceGroup_01_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  # Missing display_name
  mail_nickname    = "testgroup"
  mail_enabled     = false
  security_enabled = true
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name = "Test Group"
  # Missing mail_nickname
  mail_enabled     = false
  security_enabled = true
}
`,
				ExpectError: regexp.MustCompile(`The argument "mail_nickname" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name  = "Test Group"
  mail_nickname = "testgroup"
  # Missing mail_enabled
  security_enabled = true
}
`,
				ExpectError: regexp.MustCompile(`The argument "mail_enabled" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name  = "Test Group"
  mail_nickname = "testgroup"
  mail_enabled  = false
  # Missing security_enabled
}
`,
				ExpectError: regexp.MustCompile(`The argument "security_enabled" is required`),
			},
		},
	})
}

// TestGroupResource_ErrorHandling tests error scenarios
func TestUnitResourceGroup_02_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_nickname    = "testgroup"
  mail_enabled     = false
  security_enabled = true
}
`,
				ExpectError: regexp.MustCompile(`Bad Request|BadRequest`),
			},
		},
	})
}

// TestGroupResource_InvalidValues tests invalid value validation
func TestUnitResourceGroup_03_InvalidValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_nickname    = "testgroup"
  mail_enabled     = false
  security_enabled = true
  visibility       = "invalid"
}
`,
				ExpectError: regexp.MustCompile("Attribute visibility value must be one of"),
			},
		},
	})
}

// TestGroupResource_MailNicknameValidation tests mail_nickname validation
func TestUnitResourceGroup_04_MailNicknameValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "Test Group"
  mail_nickname    = "test@group"
  mail_enabled     = false
  security_enabled = true
}
`,
				ExpectError: regexp.MustCompile("String contains forbidden character"),
			},
		},
	})
}

// TestGroupResource_Scenario1_SecurityGroupAssigned tests security group with assigned membership
func TestUnitResourceGroup_05_Scenario1_SecurityGroupAssigned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_1_security_group_assigned.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_1").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_1").Key("display_name").HasValue("acc-security-group-with-assigned-membership-type"),
					check.That(resourceType+".scenario_1").Key("mail_nickname").HasValue("c660a1b4-5"),
					check.That(resourceType+".scenario_1").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_1").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_1").Key("description").HasValue("test"),
				),
			},
			{
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

// TestGroupResource_Scenario2_SecurityGroupDynamicUser tests security group with dynamic user membership
func TestUnitResourceGroup_06_Scenario2_SecurityGroupDynamicUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_2_security_group_dynamic_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_2").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_2").Key("display_name").HasValue("acc-security-group-with-dynamic-user-membership-type"),
					check.That(resourceType+".scenario_2").Key("mail_nickname").HasValue("f9a72987-7"),
					check.That(resourceType+".scenario_2").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_2").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_2").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_2").Key("membership_rule").HasValue("(user.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_2").Key("membership_rule_processing_state").HasValue("On"),
				),
			},
			{
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

// TestGroupResource_Scenario3_SecurityGroupDynamicDevice tests security group with dynamic device membership
func TestUnitResourceGroup_07_Scenario3_SecurityGroupDynamicDevice(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_3_security_group_dynamic_device.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_3").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_3").Key("display_name").HasValue("acc-security-group-with-dynamic-device-membership-type"),
					check.That(resourceType+".scenario_3").Key("mail_nickname").HasValue("17bf0e02-0"),
					check.That(resourceType+".scenario_3").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_3").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_3").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_3").Key("membership_rule").HasValue("(device.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_3").Key("membership_rule_processing_state").HasValue("On"),
				),
			},
			{
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

// TestGroupResource_Scenario4_SecurityGroupRoleAssignable tests security group with Entra role assignment capability
func TestUnitResourceGroup_08_Scenario4_SecurityGroupRoleAssignable(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_4_security_group_role_assignable.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_4").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_4").Key("display_name").HasValue("acc-security-group-with-entra-role-assignment"),
					check.That(resourceType+".scenario_4").Key("mail_nickname").HasValue("dec34327-9"),
					check.That(resourceType+".scenario_4").Key("mail_enabled").HasValue("false"),
					check.That(resourceType+".scenario_4").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_4").Key("is_assignable_to_role").HasValue("true"),
					check.That(resourceType+".scenario_4").Key("visibility").HasValue("Private"),
				),
			},
			{
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

// TestGroupResource_Scenario5_M365GroupDynamicUser tests M365 group with dynamic user membership
func TestUnitResourceGroup_09_Scenario5_M365GroupDynamicUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_5_m365_group_dynamic_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_5").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_5").Key("display_name").HasValue("acc-m365-group-with-dynamic-user-membership-type"),
					check.That(resourceType+".scenario_5").Key("mail_nickname").HasValue("some-string"),
					check.That(resourceType+".scenario_5").Key("mail_enabled").HasValue("true"),
					check.That(resourceType+".scenario_5").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_5").Key("group_types.#").HasValue("2"),
					check.That(resourceType+".scenario_5").Key("membership_rule").HasValue("(user.accountEnabled -eq true)"),
					check.That(resourceType+".scenario_5").Key("membership_rule_processing_state").HasValue("On"),
					check.That(resourceType+".scenario_5").Key("visibility").HasValue("Private"),
				),
			},
			{
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

// TestGroupResource_Scenario6_M365GroupAssigned tests M365 group with assigned membership
func TestUnitResourceGroup_10_Scenario6_M365GroupAssigned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_6_m365_group_assigned.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scenario_6").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scenario_6").Key("display_name").HasValue("acc-m365-group-with-assigned-membership-type"),
					check.That(resourceType+".scenario_6").Key("mail_nickname").HasValue("acc-m365-group-with-assigned-membership-type"),
					check.That(resourceType+".scenario_6").Key("mail_enabled").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("security_enabled").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("group_types.#").HasValue("1"),
					check.That(resourceType+".scenario_6").Key("description").HasValue("something"),
					check.That(resourceType+".scenario_6").Key("is_assignable_to_role").HasValue("true"),
					check.That(resourceType+".scenario_6").Key("visibility").HasValue("Private"),
				),
			},
			{
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

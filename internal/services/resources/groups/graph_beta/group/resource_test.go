package graphBetaGroup_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	return mockClient, groupMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterErrorMocks()

	return mockClient, groupMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// TestGroupResource_RequiredFields tests required field validation
func TestGroupResource_RequiredFields(t *testing.T) {
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
func TestGroupResource_ErrorHandling(t *testing.T) {
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
func TestGroupResource_InvalidValues(t *testing.T) {
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
func TestGroupResource_MailNicknameValidation(t *testing.T) {
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
func TestGroupResource_Scenario1_SecurityGroupAssigned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_1_security_group_assigned.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "display_name", "acc-security-group-with-assigned-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "mail_nickname", "c660a1b4-5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "description", "test"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_1", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestGroupResource_Scenario2_SecurityGroupDynamicUser tests security group with dynamic user membership
func TestGroupResource_Scenario2_SecurityGroupDynamicUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_2_security_group_dynamic_user.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "display_name", "acc-security-group-with-dynamic-user-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "mail_nickname", "f9a72987-7"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "membership_rule", "(user.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "membership_rule_processing_state", "On"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_2", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestGroupResource_Scenario3_SecurityGroupDynamicDevice tests security group with dynamic device membership
func TestGroupResource_Scenario3_SecurityGroupDynamicDevice(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_3_security_group_dynamic_device.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "display_name", "acc-security-group-with-dynamic-device-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "mail_nickname", "17bf0e02-0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "membership_rule", "(device.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "membership_rule_processing_state", "On"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_3", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestGroupResource_Scenario4_SecurityGroupRoleAssignable tests security group with Entra role assignment capability
func TestGroupResource_Scenario4_SecurityGroupRoleAssignable(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_4_security_group_role_assignable.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "display_name", "acc-security-group-with-entra-role-assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "mail_nickname", "dec34327-9"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "is_assignable_to_role", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "visibility", "Private"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_4", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestGroupResource_Scenario5_M365GroupDynamicUser tests M365 group with dynamic user membership
func TestGroupResource_Scenario5_M365GroupDynamicUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_5_m365_group_dynamic_user.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "display_name", "acc-m365-group-with-dynamic-user-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "mail_nickname", "some-string"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "group_types.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "membership_rule", "(user.accountEnabled -eq true)"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "membership_rule_processing_state", "On"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "visibility", "Private"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_5", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestGroupResource_Scenario6_M365GroupAssigned tests M365 group with assigned membership
func TestGroupResource_Scenario6_M365GroupAssigned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_scenario_6_m365_group_assigned.tf"))
	if err != nil {
		t.Fatalf("Failed to read test config: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(content),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.scenario_6"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "display_name", "acc-m365-group-with-assigned-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "mail_nickname", "acc-m365-group-with-assigned-membership-type"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "group_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "description", "something"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "is_assignable_to_role", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "visibility", "Private"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_groups_group.scenario_6", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

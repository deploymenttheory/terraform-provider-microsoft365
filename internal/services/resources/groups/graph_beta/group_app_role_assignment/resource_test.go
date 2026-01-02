package graphBetaGroupAppRoleAssignment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	appRoleAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_app_role_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const (
	resourceType = "microsoft365_graph_beta_groups_group_app_role_assignment"
)

func setupMockEnvironment() (*mocks.Mocks, *appRoleAssignmentMocks.GroupAppRoleAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &appRoleAssignmentMocks.GroupAppRoleAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *appRoleAssignmentMocks.GroupAppRoleAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &appRoleAssignmentMocks.GroupAppRoleAssignmentMock{}
	assignmentMock.RegisterErrorMocks()
	return mockClient, assignmentMock
}

func testConfigMinimal() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	return config
}

func testConfigMaximal() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	return config
}

// TestUnitGroupAppRoleAssignmentResource_Minimal tests minimal configuration
func TestUnitGroupAppRoleAssignmentResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("target_group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("resource_object_id").HasValue("00000000-0000-0000-0000-000000000010"),
					check.That(resourceType+".minimal").Key("app_role_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".minimal").Key("principal_display_name").HasValue("Minimal Group"),
					check.That(resourceType+".minimal").Key("resource_display_name").HasValue("Microsoft Graph"),
					check.That(resourceType+".minimal").Key("principal_type").HasValue("Group"),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".minimal"]
					if !ok {
						return "", fmt.Errorf("Resource not found: %s", resourceType+".minimal")
					}
					groupID := rs.Primary.Attributes["target_group_id"]
					assignmentID := rs.Primary.ID
					return fmt.Sprintf("%s/%s", groupID, assignmentID), nil
				},
			},
		},
	})
}

// TestUnitGroupAppRoleAssignmentResource_Maximal tests maximal configuration
func TestUnitGroupAppRoleAssignmentResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("target_group_id").HasValue("00000000-0000-0000-0000-000000000003"),
					check.That(resourceType+".maximal").Key("resource_object_id").HasValue("00000000-0000-0000-0000-000000000011"),
					check.That(resourceType+".maximal").Key("app_role_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".maximal").Key("principal_display_name").HasValue("Maximal Group"),
					check.That(resourceType+".maximal").Key("resource_display_name").HasValue("SharePoint Online"),
					check.That(resourceType+".maximal").Key("principal_type").HasValue("Group"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".maximal"]
					if !ok {
						return "", fmt.Errorf("Resource not found: %s", resourceType+".maximal")
					}
					groupID := rs.Primary.Attributes["target_group_id"]
					assignmentID := rs.Primary.ID
					return fmt.Sprintf("%s/%s", groupID, assignmentID), nil
				},
			},
		},
	})
}

// TestUnitGroupAppRoleAssignmentResource_Delete tests resource deletion
func TestUnitGroupAppRoleAssignmentResource_Delete(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".minimal").Key("id").Exists(),
				),
			},
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					_, exists := s.RootModule().Resources[resourceType+".minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupAppRoleAssignmentResource_RequiredFields tests required field validation
func TestUnitGroupAppRoleAssignmentResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  resource_object_id = "00000000-0000-0000-0000-000000000010"
  app_role_id        = "00000000-0000-0000-0000-000000000000"
}
`,
				ExpectError: regexp.MustCompile(`The argument "target_group_id" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  target_group_id = "00000000-0000-0000-0000-000000000002"
  app_role_id     = "00000000-0000-0000-0000-000000000000"
}
`,
				ExpectError: regexp.MustCompile(`The argument "resource_object_id" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  target_group_id    = "00000000-0000-0000-0000-000000000002"
  resource_object_id = "00000000-0000-0000-0000-000000000010"
}
`,
				ExpectError: regexp.MustCompile(`The argument "app_role_id" is required`),
			},
		},
	})
}

// TestUnitGroupAppRoleAssignmentResource_InvalidValues tests invalid value validation
func TestUnitGroupAppRoleAssignmentResource_InvalidValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  target_group_id    = "invalid-uuid"
  resource_object_id = "00000000-0000-0000-0000-000000000010"
  app_role_id        = "00000000-0000-0000-0000-000000000000"
}
`,
				ExpectError: regexp.MustCompile(`Must be a valid UUID format`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  target_group_id    = "00000000-0000-0000-0000-000000000002"
  resource_object_id = "invalid-uuid"
  app_role_id        = "00000000-0000-0000-0000-000000000000"
}
`,
				ExpectError: regexp.MustCompile(`Must be a valid UUID format`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "test" {
  target_group_id    = "00000000-0000-0000-0000-000000000002"
  resource_object_id = "00000000-0000-0000-0000-000000000010"
  app_role_id        = "invalid-uuid"
}
`,
				ExpectError: regexp.MustCompile(`Must be a valid UUID format`),
			},
		},
	})
}

// TestUnitGroupAppRoleAssignmentResource_ErrorHandling tests API error handling
func TestUnitGroupAppRoleAssignmentResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile(`Bad Request|400|ApiError`),
			},
		},
	})
}

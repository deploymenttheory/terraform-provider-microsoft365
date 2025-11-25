package graphBetaGroupLifecycleExpirationPolicyAssignment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupLifecycleExpirationPolicyAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_lifecycle_expiration_policy_assignment"
	assignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_lifecycle_expiration_policy_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaGroupLifecycleExpirationPolicyAssignment.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *assignmentMocks.GroupLifecycleExpirationPolicyAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &assignmentMocks.GroupLifecycleExpirationPolicyAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *assignmentMocks.GroupLifecycleExpirationPolicyAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &assignmentMocks.GroupLifecycleExpirationPolicyAssignmentMock{}
	assignmentMock.RegisterErrorMocks()
	return mockClient, assignmentMock
}

func TestUnitGroupLifecyclePolicyAssignmentResource_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("group_id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_Delete(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("id").Exists(),
				),
			},
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					_, exists := s.RootModule().Resources[resourceType+".test"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment" "test" { }`,
				ExpectError: regexp.MustCompile(`The argument "group_id" is required`),
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_InvalidGroupID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment" "test" {
  group_id = "not-a-valid-uuid"
}
`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_ErrorHandling_NoPolicyExists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &assignmentMocks.GroupLifecycleExpirationPolicyAssignmentMock{}
	assignmentMock.RegisterNoPolicyErrorMocks() // Use specific no-policy error mocks
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`No lifecycle policy found|no lifecycle policy exists`),
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_ErrorHandling_BadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Bad Request|400|ApiError`),
			},
		},
	})
}

func TestUnitGroupLifecyclePolicyAssignmentResource_UpdateNotSupported(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("id").Exists(),
				),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "selected" {
  group_lifetime_in_days        = 180
  managed_group_types           = "Selected"
  alternate_notification_emails = "admin@deploymenttheory.com"
}

resource "microsoft365_graph_beta_groups_group" "test_group" {
  display_name     = "Test M365 Group for Lifecycle Policy"
  mail_enabled     = true
  security_enabled = true
  group_types      = ["Unified"]
  description      = "Unit test - M365 group for lifecycle policy assignment"
  mail_nickname    = "testm365group"
  visibility       = "Private"
}

resource "microsoft365_graph_beta_groups_group" "test_group_2" {
  display_name     = "Second Test M365 Group"
  mail_enabled     = true
  security_enabled = true
  group_types      = ["Unified"]
  description      = "Second test group"
  mail_nickname    = "testm365group2"
  visibility       = "Private"
}

resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment" "test" {
  group_id = microsoft365_graph_beta_groups_group.test_group_2.id
}
`,
				// Changing group_id should force replacement (RequiresReplace)
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Config loader functions
func testConfigBasic() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_basic.tf")
	if err != nil {
		panic("failed to load basic config: " + err.Error())
	}
	return config
}

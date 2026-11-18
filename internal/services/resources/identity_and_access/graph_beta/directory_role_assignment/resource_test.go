package graphBetaDirectoryRoleAssignment_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/directory_role_assignment"
	assignmentmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/directory_role_assignment/mocks"
)

var resourceType = assignment.ResourceName

func setupMockEnvironment(t *testing.T) *assignmentmocks.DirectoryRoleAssignmentMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &assignmentmocks.DirectoryRoleAssignmentMock{}
	mock.CleanupMockState()
	mock.RegisterMocks()
	t.Cleanup(func() { mock.CleanupMockState(); httpmock.DeactivateAndReset() })
	return mock
}

func loadUnitTestTerraform(t *testing.T, file string) string {
	t.Helper()
	value, err := helpers.ParseHCLFile("tests/terraform/unit/" + file)
	require.NoError(t, err)
	return value
}

func TestUnitResourceDirectoryRoleAssignment_Lifecycle(t *testing.T) {
	setupMockEnvironment(t)
	config := loadUnitTestTerraform(t, "resource_minimal.tf")
	address := resourceType + ".test"
	replacement := strings.Replace(config, "11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", 1)
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "id", assignmentmocks.AssignmentID),
				resource.TestCheckResourceAttr(address, "directory_scope_id", "/"),
			),
		},
		{ResourceName: address, ImportState: true, ImportStateVerify: true},
		{
			Config:           replacement,
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionDestroyBeforeCreate)}},
		},
		{Config: replacement, PlanOnly: true},
	}})
}

func TestUnitResourceDirectoryRoleAssignment_GeneratedImport(t *testing.T) {
	mock := setupMockEnvironment(t)
	mock.Seed("/")
	address := resourceType + ".test"
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: loadUnitTestTerraform(
				t,
				"import.tf",
			),
			ResourceName:           address,
			ImportState:            true,
			ImportStateId:          assignmentmocks.AssignmentID,
			GenerateConfig:         true,
			ImportStateKind:        resource.ImportBlockWithID,
			ImportStateConfigExact: true,
			ImportPlanChecks:       resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
		},
	}})
}

func TestUnitResourceDirectoryRoleAssignment_RejectScopedImport(t *testing.T) {
	mock := setupMockEnvironment(t)
	mock.Seed("/administrativeUnits/22222222-2222-4222-8222-222222222222")
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: loadUnitTestTerraform(
				t,
				"resource_minimal.tf",
			),
			ResourceName:  resourceType + ".test",
			ImportState:   true,
			ImportStateId: assignmentmocks.AssignmentID,
			ExpectError:   regexp.MustCompile("Unsupported directory role assignment scope"),
		},
	}})
}

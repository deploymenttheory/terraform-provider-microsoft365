package graphBetaRoleScopeTag_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitResourceRoleScopeTag_10_AssignmentRemoval(t *testing.T) {
	t.Setenv("TF_ACC", "0")
	t.Setenv("MS365_TEST_MODE", "true")
	_, tagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tagMock.CleanupMockState()
	config := func(assignments string) string {
		return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
 display_name = "Role Scope Tag Assignment Removal"
 description = "Assignment removal regression"
 %s
}`, assignments)
	}
	assigned := config(`assignments = [{type = "groupAssignmentTarget", group_id = "11111111-1111-1111-1111-111111111111"}]`)
	empty := config(`assignments = []`)
	omitted := config("")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: assigned, Check: resource.TestCheckResourceAttr(resourceType+".test", "assignments.#", "1")},
			{Config: empty, Check: resource.TestCheckResourceAttr(resourceType+".test", "assignments.#", "0")},
			{Config: empty, PlanOnly: true},
			{Config: assigned, Check: resource.TestCheckResourceAttr(resourceType+".test", "assignments.#", "1")},
			{Config: omitted, Check: resource.TestCheckNoResourceAttr(resourceType+".test", "assignments.#")},
			{Config: omitted, PlanOnly: true},
		},
	})
}

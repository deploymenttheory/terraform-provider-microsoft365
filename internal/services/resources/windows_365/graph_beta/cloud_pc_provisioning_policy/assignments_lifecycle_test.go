package graphBetaCloudPcProvisioningPolicy_test

import (
	"fmt"
	"strings"
	"testing"

	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_provisioning_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// testConfigFrontlineAssignment returns a sharedByEntraGroup policy with one Frontline group assignment.
func testConfigFrontlineAssignment(allotmentName string, licenceCount int) string {
	assignment := fmt.Sprintf(`provisioning_type = "sharedByEntraGroup"
  assignments = [
    {
      type                    = "groupAssignmentTarget"
      group_id                = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
      service_plan_id         = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
      allotment_display_name  = "%s"
      allotment_license_count = %d
    }
  ]
  display_name =`, allotmentName, licenceCount)
	return strings.Replace(testConfigMinimalWithResourceName("frontline"), "display_name =", assignment, 1)
}

// Regression test for #4928: Graph returns the assignment under an ID that differs from the group ID.
// State after create and update, import and a follow-up plan must all keep the target's group and Frontline fields.
func TestUnitResourceCloudPcProvisioningPolicy_FrontlineAssignmentLifecycle(t *testing.T) {
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	(&policyMocks.CloudPcProvisioningPolicyMock{}).RegisterMocks()
	name := "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.frontline"

	var steps []resource.TestStep
	for _, phase := range []struct {
		allotmentName string
		licenceCount  int
	}{{"example-allotment", 5}, {"example-allotment-updated", 7}} {
		config := testConfigFrontlineAssignment(phase.allotmentName, phase.licenceCount)
		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(name, "assignments.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
				"type":                    "groupAssignmentTarget",
				"group_id":                "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
				"service_plan_id":         "cccccccc-cccc-4ccc-8ccc-cccccccccccc",
				"allotment_display_name":  phase.allotmentName,
				"allotment_license_count": fmt.Sprintf("%d", phase.licenceCount),
			}),
		)
		steps = append(steps,
			resource.TestStep{Config: config, Check: check},
			resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true},
			resource.TestStep{Config: config, PlanOnly: true},
		)
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: steps})
}

package graphBetaCloudPcProvisioningPolicy_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_provisioning_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitResourceCloudPcProvisioningPolicy_UserExperienceLifecycle(t *testing.T) {
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	(&policyMocks.CloudPcProvisioningPolicyMock{}).RegisterMocks()
	name := "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.lifecycle"
	base := testConfigMinimalWithResourceName("lifecycle")
	var policyID string
	var steps []resource.TestStep
	for i, phase := range []struct{ config, value string }{{``, "cloudPc"},
		{`provisioning_type = "sharedByEntraGroup"
user_experience_type = "cloudApp"`, "cloudApp"},
		{`provisioning_type = "sharedByEntraGroup"`, "cloudApp"},
		{`provisioning_type = "sharedByEntraGroup"
user_experience_type = "cloudPc"`, "cloudPc"}} {
		config := strings.Replace(base, "display_name =", phase.config+"\ndisplay_name =", 1)
		checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(name, "user_experience_type", phase.value)}
		checks = append(checks, resource.TestCheckResourceAttrWith(name, "id", func(value string) error {
			if i == 2 && value != policyID {
				return fmt.Errorf("omitting the experience replaced the policy")
			}
			if (i == 1 || i == 3) && value == policyID {
				return fmt.Errorf("changing the experience did not replace the policy")
			}
			policyID = value
			return nil
		}))

		if i > 0 {
			config = strings.Replace(config, "Test Minimal Provisioning Policy", fmt.Sprintf("Updated policy %d", i), 1)
		}
		steps = append(steps, resource.TestStep{Config: config, Check: resource.ComposeTestCheckFunc(checks...)}, resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true}, resource.TestStep{Config: config, PlanOnly: true})
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: steps})
}
func TestUnitResourceCloudPcProvisioningPolicy_UserExperienceLifecycleInvalidProvisioningType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{{Config: strings.Replace(testConfigMinimalWithResourceName("invalid"), "display_name =", `provisioning_type = "dedicated"
user_experience_type = "cloudApp"`+"\ndisplay_name =", 1), ExpectError: regexp.MustCompile("cloudApp requires")}}})
}

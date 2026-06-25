package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaSPAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_token_lifetime_policy_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	testResource = graphBetaSPAssignment.ServicePrincipalTokenLifetimePolicyAssignmentTestResource{}
)

func TestAccResourceServicePrincipalTokenLifetimePolicyAssignment_01_Basic(t *testing.T) {
	spID := os.Getenv("ARM_SP_OBJECT_ID")
	if spID == "" {
		t.Skip("ARM_SP_OBJECT_ID not set, skipping acceptance test")
	}
	policyID := os.Getenv("ARM_TOKEN_LIFETIME_POLICY_ID")
	if policyID == "" {
		t.Skip("ARM_TOKEN_LIFETIME_POLICY_ID not set, skipping acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating SP token lifetime policy assignment")
				},
				Config: testAccConfigBasic(spID, policyID),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("SP token lifetime policy assignment", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".basic").ExistsInGraph(testResource),
					check.That(resourceType+".basic").Key("service_principal_id").HasValue(spID),
					check.That(resourceType+".basic").Key("token_lifetime_policy_id").HasValue(policyID),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing SP token lifetime policy assignment")
				},
				ResourceName:      resourceType + ".basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%s", spID, policyID),
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccConfigBasic(spID, policyID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment" "basic" {
  service_principal_id     = %q
  token_lifetime_policy_id = %q
}
`, spID, policyID)
}

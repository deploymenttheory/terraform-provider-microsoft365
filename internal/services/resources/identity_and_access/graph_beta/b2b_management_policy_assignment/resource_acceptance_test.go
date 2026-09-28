package graphBetaIdentityAndAccessB2bManagementPolicyAssignment_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaB2bManagementPolicyAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testResource = graphBetaB2bManagementPolicyAssignment.B2bManagementPolicyAssignmentTestResource{}
)

func TestAccResourceB2bManagementPolicyAssignment_01_ServicePrincipal(t *testing.T) {
	spID := os.Getenv("ARM_SP_OBJECT_ID")
	if spID == "" {
		t.Skip("ARM_SP_OBJECT_ID not set, skipping acceptance test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			10*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Applying B2B management policy to service principal")
				},
				Config: testAccConfigServicePrincipal(spID),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("directory_object_id").HasValue(spID),
					check.That(resourceType+".test").Key("directory_object_type").HasValue("servicePrincipal"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing B2B management policy assignment")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccConfigServicePrincipal(spID string) string {
	return acceptance.ConfiguredM365ProviderBlock(fmt.Sprintf(`
resource "random_uuid" "suffix" {}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "test" {
  display_name = "acc-test-b2b-assignment-${random_uuid.suffix.result}"
  definition   = ["{\"B2BManagementPolicy\":{\"InvitationsAllowedAndBlockedDomainsPolicy\":{\"BlockedDomains\":[\"example.com\"]}}}"]
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "test" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.test.id
  directory_object_id      = %q
}
`, spID))
}

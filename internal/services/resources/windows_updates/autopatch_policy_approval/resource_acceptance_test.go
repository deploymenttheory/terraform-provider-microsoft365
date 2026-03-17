package graphBetaWindowsUpdatesAutopatchPolicyApproval_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesPolicyApproval "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_policy_approval"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaWindowsUpdatesPolicyApproval.ResourceName
	testResource = graphBetaWindowsUpdatesPolicyApproval.WindowsUpdatesAutopatchPolicyApprovalTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsUpdatePolicyApproval_01_ApprovedAndImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create policy approval with status=approved")
				},
				Config: loadAcceptanceTestTerraform("01_approved.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("policy_id").Exists(),
					check.That(resourceType+".test").Key("catalog_entry_id").Exists(),
					check.That(resourceType+".test").Key("status").HasValue("approved"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import policy approval")
				},
				ResourceName: resourceType + ".test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test")
					}
					policyId := rs.Primary.Attributes["policy_id"]
					approvalId := rs.Primary.ID
					return policyId + "/" + approvalId, nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsUpdatePolicyApproval_02_ApprovedToSuspended(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create policy approval with status=approved")
				},
				Config: loadAcceptanceTestTerraform("01_approved.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("status").HasValue("approved"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update status to suspended")
					testlog.WaitForConsistency(fmt.Sprintf("%s (update)", resourceType), 5*time.Second)
					time.Sleep(5 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("02_suspended.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("status").HasValue("suspended"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

package graphBetaNetworkCloudFirewallPolicy_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloudMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy"

func TestUnitResourceNetworkCloudFirewallPolicy_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &cloudMocks.CloudFirewallMock{}
	mock.RegisterMocks()
	defer mock.CleanupMockState()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.ReplaceAll(config, "Managed by Terraform", "Updated description")
	updated = strings.ReplaceAll(updated, "enabled     = false", "enabled     = true")
	noDescription := strings.ReplaceAll(updated, `  description    = "Updated description"
`, "")
	noDescription = strings.ReplaceAll(noDescription, `  description = "Match example destination addresses"
`, "")
	var id string
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("id").Exists(), func(s *terraform.State) error {
			id = s.RootModule().Resources[resourceType+".test"].Primary.ID
			return nil
		})},
		{Config: updated, Check: func(s *terraform.State) error {
			if s.RootModule().Resources[resourceType+".test"].Primary.ID != id {
				return fmt.Errorf("update replaced ID")
			}
			return nil
		}},
		{Config: noDescription, Check: check.That(resourceType + ".test").Key("description").DoesNotExist()},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateKind: resource.ImportBlockWithResourceIdentity, ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop)}}},
		{Config: noDescription, PlanOnly: true},
	}})
}

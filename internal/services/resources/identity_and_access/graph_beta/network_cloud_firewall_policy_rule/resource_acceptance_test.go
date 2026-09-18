package graphBetaNetworkCloudFirewallPolicyRule_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy_rule"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNetworkCloudFirewallPolicyRule_01_Lifecycle(t *testing.T) {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource.tf")
	if err != nil {
		t.Fatal(err)
	}
	// Keep credentials in the environment rather than embedding them in test HCL.
	config = acceptance.ProviderConfig() + config
	updated := strings.ReplaceAll(config, "Managed by Terraform", "Updated description")
	updated = strings.ReplaceAll(updated, "enabled     = false", "enabled     = true")
	noDescription := strings.ReplaceAll(updated, `  description    = "Updated description"
`, "")
	noDescription = strings.ReplaceAll(noDescription, `  description = "Match example destination addresses"
`, "")
	var id string
	resource.Test(t, resource.TestCase{PreCheck: func() { mocks.TestAccPreCheck(t) }, ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories, CheckDestroy: destroy.CheckDestroyedAllFunc(cloud.NetworkCloudFirewallPolicyRuleTestResource{}, resourceType, 0), ExternalProviders: map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}}, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeTestCheckFunc(check.That(resourceType+".test").Key("id").Exists(), func(s *terraform.State) error {
			id = s.RootModule().Resources[resourceType+".test"].Primary.ID
			t.Logf("Acceptance created policy ID: %s; rule ID: %s", s.RootModule().Resources[resourceType+".test"].Primary.Attributes["policy_id"], id)
			return nil
		})},
		{Config: updated, Check: func(s *terraform.State) error {
			if s.RootModule().Resources[resourceType+".test"].Primary.ID != id {
				return fmt.Errorf("update replaced ID")
			}
			return nil
		}},
		{Config: noDescription, Check: check.That(resourceType + ".test").Key("description").DoesNotExist()},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true, ImportStateIdFunc: func(s *terraform.State) (string, error) {
			v := s.RootModule().Resources[resourceType+".test"].Primary
			return v.Attributes["policy_id"] + "/" + v.ID, nil
		}},
		{ResourceName: resourceType + ".test", ImportState: true, ImportStateKind: resource.ImportBlockWithResourceIdentity, ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop)}}},
		{Config: noDescription, PlanOnly: true},
	}})
}

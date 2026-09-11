package graphBetaNetworkThreatIntelligencePolicyRule_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_threat_intelligence_policy"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_threat_intelligence_policy_rule"
)

func TestAccResourceNetworkThreatIntelligencePolicyRule_01_Lifecycle(t *testing.T) {
	testResource := target.NetworkThreatIntelligencePolicyRuleTestResource{}
	config := func(file string) string {
		v, e := helpers.ParseHCLFile("tests/terraform/acceptance/" + file)
		if e != nil {
			t.Fatal(e)
		}
		return acceptance.ProviderConfig() + v
	}
	var savedID string
	sameID := func(s *terraform.State) error {
		id := s.RootModule().Resources[resourceType+".test"].Primary.ID
		t.Logf("Threat intelligence rule ID: %s; parent policy ID: %s", id,
			s.RootModule().Resources[resourceType+".test"].Primary.Attributes["threat_intelligence_policy_id"])
		if savedID == "" {
			savedID = id
		} else if id != savedID {
			return fmt.Errorf("in-place update replaced resource: %s -> %s", savedID, id)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(0,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: policy.ResourceName, TestResource: policy.NetworkThreatIntelligencePolicyTestResource{}},
		),
		Steps: []resource.TestStep{
			{
				Config: config("resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					sameID,
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					v := s.RootModule().Resources[resourceType+".test"]
					return v.Primary.Attributes["threat_intelligence_policy_id"] + "/" + v.Primary.ID, nil
				},
			},
			{
				Config: config("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					sameID,
					check.That(resourceType+".test").Key("enabled").HasValue("false"),
					check.That(resourceType+".test").Key("status").HasValue("disabled"),
					check.That(resourceType+".test").Key("priority").HasValue("65001"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("description").DoesNotExist(),
				),
			},
			{
				Config: config("resource_empty_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					sameID,
					check.That(resourceType+".test").Key("description").HasValue(""),
				),
			},
		},
	})
}

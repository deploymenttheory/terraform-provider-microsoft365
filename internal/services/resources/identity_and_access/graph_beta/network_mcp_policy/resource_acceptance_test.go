package graphBetaNetworkMCPPolicy_test

import (
	"fmt"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccResourceNetworkMCPPolicy_01_Lifecycle(t *testing.T) {
	testResource := target.NetworkMCPPolicyTestResource{}
	config := func(file string) string {
		v, e := helpers.ParseHCLFile("tests/terraform/acceptance/" + file)
		if e != nil {
			t.Fatal(e)
		}
		return acceptance.ConfiguredM365ProviderBlock(v)
	}
	var savedID string
	sameID := func(s *terraform.State) error {
		id := s.RootModule().Resources[resourceType+".test"].Primary.ID
		if savedID == "" {
			savedID = id
		} else if id != savedID {
			return fmt.Errorf("in-place update replaced resource: %s -> %s", savedID, id)
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { mocks.TestAccPreCheck(t) }, ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		CheckDestroy:      destroy.CheckDestroyedAllFunc(testResource, resourceType, 0),
		Steps: []resource.TestStep{
			{Config: config("resource.tf"), Check: resource.ComposeTestCheckFunc(sameID, check.That(resourceType+".test").ExistsInGraph(testResource))},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
			{Config: config("resource_updated.tf"), Check: resource.ComposeTestCheckFunc(sameID, check.That(resourceType+".test").Key("version").HasValue("1.0.0"), check.That(resourceType+".test").ExistsInGraph(testResource), check.That(resourceType+".test").Key("description").DoesNotExist())},
			{Config: config("resource_empty_description.tf"), Check: resource.ComposeTestCheckFunc(sameID, check.That(resourceType+".test").Key("description").HasValue(""))},
		}})
}

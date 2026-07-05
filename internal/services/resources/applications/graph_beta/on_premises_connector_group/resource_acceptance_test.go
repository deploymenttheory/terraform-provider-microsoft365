package graphBetaApplicationsOnPremisesConnectorGroup_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationsOnPremisesConnectorGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/on_premises_connector_group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceType = graphBetaApplicationsOnPremisesConnectorGroup.ResourceName
	testResource = graphBetaApplicationsOnPremisesConnectorGroup.OnPremisesConnectorGroupTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceConnectorGroup_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
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
					testlog.StepAction(resourceType, "Creating minimal connector group")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("name").MatchesRegex(regexp.MustCompile(`^acctest-connector-group-`)),
					check.That(resourceType+".minimal").Key("connector_group_type").HasValue("applicationProxy"),
					check.That(resourceType+".minimal").Key("is_default").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal connector group")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceConnectorGroup_02_WithRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
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
					testlog.StepAction(resourceType, "Creating connector group with region")
				},
				Config: loadAcceptanceTestTerraform("resource_with_region.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_region").ExistsInGraph(testResource),
					check.That(resourceType+".with_region").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".with_region").Key("name").MatchesRegex(regexp.MustCompile(`^acctest-connector-group-region-`)),
					check.That(resourceType+".with_region").Key("region").HasValue("nam"),
					check.That(resourceType+".with_region").Key("connector_group_type").HasValue("applicationProxy"),
					check.That(resourceType+".with_region").Key("is_default").HasValue("false"),
				),
			},
		},
	})
}

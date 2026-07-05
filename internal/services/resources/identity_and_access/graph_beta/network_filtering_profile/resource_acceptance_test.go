package graphBetaNetworkFilteringProfile_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkFilteringProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_profile"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaNetworkFilteringProfile.ResourceName

	// testResource is the test resource implementation for filtering profiles
	testResource = graphBetaNetworkFilteringProfile.NetworkFilteringProfileTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceNetworkFilteringProfile_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("name").IsNotEmpty(),
					check.That(resourceType+".test").Key("priority").HasValue("100"),
					check.That(resourceType+".test").Key("state").HasValue("enabled"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: loadAcceptanceTestTerraform("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("name").IsNotEmpty(),
					check.That(resourceType+".test").Key("description").IsNotEmpty(),
					check.That(resourceType+".test").Key("priority").HasValue("200"),
					check.That(resourceType+".test").Key("state").HasValue("disabled"),
				),
			},
		},
	})
}

func TestAccResourceNetworkFilteringProfile_02_EnabledState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_enabled_state.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".enabled").ExistsInGraph(testResource),
					check.That(resourceType+".enabled").Key("id").Exists(),
					check.That(resourceType+".enabled").Key("name").IsNotEmpty(),
					check.That(resourceType+".enabled").Key("description").IsNotEmpty(),
					check.That(resourceType+".enabled").Key("priority").HasValue("100"),
					check.That(resourceType+".enabled").Key("state").HasValue("enabled"),
					check.That(resourceType+".enabled").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccResourceNetworkFilteringProfile_03_DisabledState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_disabled_state.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".disabled").ExistsInGraph(testResource),
					check.That(resourceType+".disabled").Key("id").Exists(),
					check.That(resourceType+".disabled").Key("name").IsNotEmpty(),
					check.That(resourceType+".disabled").Key("description").IsNotEmpty(),
					check.That(resourceType+".disabled").Key("priority").HasValue("200"),
					check.That(resourceType+".disabled").Key("state").HasValue("disabled"),
					check.That(resourceType+".disabled").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccResourceNetworkFilteringProfile_04_MinimalConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
					testlog.WaitForConsistency("Microsoft Graph", 15*time.Second)
					time.Sleep(15 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_no_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("name").IsNotEmpty(),
					check.That(resourceType+".minimal").Key("priority").HasValue("300"),
					check.That(resourceType+".minimal").Key("state").HasValue("enabled"),
				),
			},
		},
	})
}

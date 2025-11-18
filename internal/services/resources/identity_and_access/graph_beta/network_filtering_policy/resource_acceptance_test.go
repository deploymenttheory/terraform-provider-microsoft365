package graphBetaNetworkFilteringPolicy_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkFilteringPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// Resource type name constructed from exported constants
	resourceType = constants.PROVIDER_NAME + "_" + graphBetaNetworkFilteringPolicy.ResourceName

	// testResource is the test resource implementation for filtering policies
	testResource = graphBetaNetworkFilteringPolicy.NetworkFilteringPolicyTestResource{}
)

func TestAccNetworkFilteringPolicyResource_Lifecycle(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccFilteringPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("action").HasValue("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      "microsoft365_graph_beta_identity_and_access_network_filtering_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: testAccFilteringPolicyConfig_updated(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("description").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.test").Key("action").HasValue("allow"),
				),
			},
		},
	})
}

func TestAccNetworkFilteringPolicyResource_BlockAction(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccFilteringPolicyConfig_blockAction(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").Key("name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").Key("description").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").Key("action").HasValue("block"),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.block").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccNetworkFilteringPolicyResource_AllowAction(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccFilteringPolicyConfig_allowAction(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").Key("name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").Key("description").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").Key("action").HasValue("allow"),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.allow").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccNetworkFilteringPolicyResource_MinimalConfiguration(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
					testlog.WaitForConsistency("Microsoft Graph", 15*time.Second)
					time.Sleep(15 * time.Second)
				},
				Config: testAccFilteringPolicyConfig_minimalNoDescription(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.minimal").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.minimal").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.minimal").Key("name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_identity_and_access_network_filtering_policy.minimal").Key("action").HasValue("block"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccFilteringPolicyConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccFilteringPolicyConfig_updated() string {
	config := mocks.LoadTerraformConfigFile("resource_updated.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccFilteringPolicyConfig_blockAction() string {
	config := mocks.LoadTerraformConfigFile("resource_block_action.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccFilteringPolicyConfig_allowAction() string {
	config := mocks.LoadTerraformConfigFile("resource_allow_action.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccFilteringPolicyConfig_minimalNoDescription() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal_no_description.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

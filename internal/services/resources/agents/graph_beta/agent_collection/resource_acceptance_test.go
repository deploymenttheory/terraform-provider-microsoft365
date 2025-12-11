package graphBetaAgentCollection_test

// currently disabled as the resource doesnt yet support DELETE
// however the resource has a delete method in the sdk
// it just doesnt work !

// import (
// 	"testing"
// 	"time"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
// 	graphBetaAgentCollection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_collection"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/terraform"
// )

// var (
// 	testResource = graphBetaAgentCollection.AgentCollectionTestResource{}
// )

// // TestAccAgentCollectionResource_Minimal tests creating an agent collection with minimal configuration
// func TestAccAgentCollectionResource_Minimal(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		CheckDestroy: destroy.CheckDestroyedAllFunc(
// 			testResource,
// 			resourceType,
// 			30*time.Second,
// 		),
// 		ExternalProviders: map[string]resource.ExternalProvider{
// 			"random": {
// 				Source:            "hashicorp/random",
// 				VersionConstraint: ">= 3.7.2",
// 			},
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				PreConfig: func() {
// 					testlog.StepAction(resourceType, "Creating agent collection with minimal configuration")
// 				},
// 				Config: testAccConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					func(_ *terraform.State) error {
// 						testlog.WaitForConsistency("agent collection", 15*time.Second)
// 						time.Sleep(15 * time.Second)
// 						return nil
// 					},
// 					check.That(resourceType+".test_minimal").ExistsInGraph(testResource),
// 					check.That(resourceType+".test_minimal").Key("id").Exists(),
// 					check.That(resourceType+".test_minimal").Key("display_name").Exists(),
// 					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
// 				),
// 			},
// 			{
// 				PreConfig: func() {
// 					testlog.StepAction(resourceType, "Importing agent collection")
// 				},
// 				ResourceName:      resourceType + ".test_minimal",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				ImportStateVerifyIgnore: []string{
// 					"timeouts",
// 				},
// 			},
// 		},
// 	})
// }

// // TestAccAgentCollectionResource_Maximal tests creating an agent collection with maximal configuration
// func TestAccAgentCollectionResource_Maximal(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		CheckDestroy: destroy.CheckDestroyedAllFunc(
// 			testResource,
// 			resourceType,
// 			30*time.Second,
// 		),
// 		ExternalProviders: map[string]resource.ExternalProvider{
// 			"random": {
// 				Source:            "hashicorp/random",
// 				VersionConstraint: ">= 3.7.2",
// 			},
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				PreConfig: func() {
// 					testlog.StepAction(resourceType, "Creating agent collection with maximal configuration")
// 				},
// 				Config: testAccConfigMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					func(_ *terraform.State) error {
// 						testlog.WaitForConsistency("agent collection", 15*time.Second)
// 						time.Sleep(15 * time.Second)
// 						return nil
// 					},
// 					check.That(resourceType+".test_maximal").ExistsInGraph(testResource),
// 					check.That(resourceType+".test_maximal").Key("id").Exists(),
// 					check.That(resourceType+".test_maximal").Key("display_name").Exists(),
// 					check.That(resourceType+".test_maximal").Key("owner_ids.#").HasValue("2"),
// 					check.That(resourceType+".test_maximal").Key("description").Exists(),
// 					check.That(resourceType+".test_maximal").Key("originating_store").HasValue("Deployment Theory"),
// 				),
// 			},
// 			{
// 				PreConfig: func() {
// 					testlog.StepAction(resourceType, "Importing agent collection")
// 				},
// 				ResourceName:      resourceType + ".test_maximal",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				ImportStateVerifyIgnore: []string{
// 					"timeouts",
// 				},
// 			},
// 		},
// 	})
// }

// func testAccConfigMinimal() string {
// 	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
// 	return acceptance.ConfiguredM365ProviderBlock(config)
// }

// func testAccConfigMaximal() string {
// 	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
// 	return acceptance.ConfiguredM365ProviderBlock(config)
// }

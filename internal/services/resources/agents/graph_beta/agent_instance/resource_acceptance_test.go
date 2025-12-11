package graphBetaAgentInstance_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentInstance "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_instance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	testResource = graphBetaAgentInstance.AgentInstanceTestResource{}
)

// TestAccAgentInstanceResource_Minimal tests creating an agent instance with minimal configuration
func TestAccAgentInstanceResource_Minimal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent instance with minimal configuration")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".test_minimal").Key("id").Exists(),
					check.That(resourceType+".test_minimal").Key("display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.protocol_version").HasValue("1.0"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.version").HasValue("1.0.1"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.supports_authenticated_extended_card").HasValue("false"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.push_notifications").HasValue("false"),
					check.That(resourceType+".test_minimal").Key("agent_card_manifest.capabilities.state_transition_history").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent instance")
				},
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccAgentInstanceResource_Maximal tests creating an agent instance with maximal configuration
func TestAccAgentInstanceResource_Maximal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent instance with maximal configuration")
				},
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".test_maximal").Key("id").Exists(),
					check.That(resourceType+".test_maximal").Key("display_name").Exists(),
					check.That(resourceType+".test_maximal").Key("originating_store").HasValue("Deployment Theory"),
					check.That(resourceType+".test_maximal").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("url").HasValue("https://servicedesk.deploymenttheory.com/api"),
					check.That(resourceType+".test_maximal").Key("preferred_transport").HasValue("HTTP+JSON"),
					check.That(resourceType+".test_maximal").Key("additional_interfaces.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.display_name").HasValue("IT Service Desk Agent"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.protocol_version").HasValue("1.0"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.version").HasValue("2.0.0"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.icon_url").HasValue("https://servicedesk.example.com/assets/agent-icon.png"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.documentation_url").HasValue("https://docs.example.com/servicedesk-agent"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.supports_authenticated_extended_card").HasValue("false"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.default_input_modes.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.default_output_modes.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.provider.organization").HasValue("Deployment Theory"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.provider.url").HasValue("https://www.deploymenttheory.com"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.push_notifications").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.state_transition_history").HasValue("false"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.capabilities.extensions.#").HasValue("1"),
					check.That(resourceType+".test_maximal").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent instance")
				},
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccAgentInstanceResource_UpdateMinimalToMaximal tests updating from minimal to maximal configuration
func TestAccAgentInstanceResource_UpdateMinimalToMaximal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Create with minimal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent instance with minimal configuration")
				},
				Config: testAccConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("id").Exists(),
					check.That(resourceType+".test_update").Key("display_name").Exists(),
					check.That(resourceType+".test_update").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("1.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("false"),
				),
			},
			// Step 2: Update to maximal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating agent instance to maximal configuration")
				},
				Config: testAccConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance update", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("id").Exists(),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("url").HasValue("https://updated-agent.example.com/api"),
					check.That(resourceType+".test_update").Key("preferred_transport").HasValue("HTTP+JSON"),
					check.That(resourceType+".test_update").Key("additional_interfaces.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("2.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.push_notifications").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
		},
	})
}

// TestAccAgentInstanceResource_UpdateMaximalToMinimal tests updating from maximal to minimal configuration
func TestAccAgentInstanceResource_UpdateMaximalToMinimal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Create with maximal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent instance with maximal configuration")
				},
				Config: testAccConfigUpdateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("id").Exists(),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("2"),
					check.That(resourceType+".test_update").Key("url").HasValue("https://updated-agent.example.com/api"),
					check.That(resourceType+".test_update").Key("additional_interfaces.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("true"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.skills.#").HasValue("1"),
				),
			},
			// Step 2: Update to minimal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating agent instance to minimal configuration")
				},
				Config: testAccConfigUpdateMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent instance update", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("id").Exists(),
					check.That(resourceType+".test_update").Key("originating_store").HasValue("Terraform"),
					check.That(resourceType+".test_update").Key("owner_ids.#").HasValue("1"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.version").HasValue("1.0.0"),
					check.That(resourceType+".test_update").Key("agent_card_manifest.capabilities.streaming").HasValue("false"),
				),
			},
		},
	})
}

func testAccConfigMinimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigMaximal() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigUpdateMinimal() string {
	config := mocks.LoadTerraformConfigFile("resource_update_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigUpdateMaximal() string {
	config := mocks.LoadTerraformConfigFile("resource_update_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

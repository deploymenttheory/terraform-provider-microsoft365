package graphBetaAgentIdentity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentity "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// accResourceType is the resource type name for acceptance tests
	accResourceType = graphBetaAgentIdentity.ResourceName

	// accTestResource is the test resource implementation for acceptance tests
	accTestResource = graphBetaAgentIdentity.AgentIdentityTestResource{}
)

// TestAccAgentIdentityResource_Minimal tests creating an agent identity with minimal configuration
func TestAccAgentIdentityResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			accResourceType,
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
					testlog.StepAction(accResourceType, "Creating agent identity with minimal configuration")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent identity", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(accResourceType+".test_minimal").ExistsInGraph(accTestResource),
					check.That(accResourceType+".test_minimal").Key("id").Exists(),
					check.That(accResourceType+".test_minimal").Key("display_name").Exists(),
					check.That(accResourceType+".test_minimal").Key("agent_identity_blueprint_id").Exists(),
					check.That(accResourceType+".test_minimal").Key("service_principal_type").HasValue("ServiceIdentity"),
					check.That(accResourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(accResourceType+".test_minimal").Key("sponsor_ids.#").HasValue("1"),
					check.That(accResourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Importing agent identity")
				},
				ResourceName: accResourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[accResourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", accResourceType+".test_minimal")
					}
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					return fmt.Sprintf("%s/%s", id, blueprintId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccAgentIdentityResource_WithTags tests creating an agent identity with tags
func TestAccAgentIdentityResource_WithTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			accResourceType,
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
					testlog.StepAction(accResourceType, "Creating agent identity with tags")
				},
				Config: testAccConfigWithTags(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent identity with tags", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(accResourceType+".test_with_tags").ExistsInGraph(accTestResource),
					check.That(accResourceType+".test_with_tags").Key("id").Exists(),
					check.That(accResourceType+".test_with_tags").Key("display_name").Exists(),
					check.That(accResourceType+".test_with_tags").Key("agent_identity_blueprint_id").Exists(),
					check.That(accResourceType+".test_with_tags").Key("service_principal_type").HasValue("ServiceIdentity"),
					check.That(accResourceType+".test_with_tags").Key("account_enabled").HasValue("true"),
					check.That(accResourceType+".test_with_tags").Key("tags.#").HasValue("3"),
					check.That(accResourceType+".test_with_tags").Key("sponsor_ids.#").HasValue("1"),
					check.That(accResourceType+".test_with_tags").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Importing agent identity with tags")
				},
				ResourceName: accResourceType + ".test_with_tags",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[accResourceType+".test_with_tags"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", accResourceType+".test_with_tags")
					}
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					return fmt.Sprintf("%s/%s", id, blueprintId), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccConfigMinimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigWithTags() string {
	config := mocks.LoadTerraformConfigFile("resource_with_tags.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

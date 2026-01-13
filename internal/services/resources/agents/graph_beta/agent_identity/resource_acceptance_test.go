package graphBetaAgentIdentity_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccAgentIdentityResource_Minimal tests creating an agent identity with minimal configuration
func TestAccAgentIdentityResource_Minimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating agent identity with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent identity", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".test_minimal").Key("id").Exists(),
					check.That(resourceType+".test_minimal").Key("display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("agent_identity_blueprint_id").Exists(),
					check.That(resourceType+".test_minimal").Key("service_principal_type").HasValue("ServiceIdentity"),
					check.That(resourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("sponsor_ids.#").HasValue("1"),
					check.That(resourceType+".test_minimal").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent identity")
				},
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s/%s:hard_delete=%s", id, blueprintId, hardDelete), nil
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
					testlog.StepAction(resourceType, "Creating agent identity with tags")
				},
				Config: loadAcceptanceTestTerraform("resource_with_tags.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent identity with tags", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_with_tags").ExistsInGraph(testResource),
					check.That(resourceType+".test_with_tags").Key("id").Exists(),
					check.That(resourceType+".test_with_tags").Key("display_name").Exists(),
					check.That(resourceType+".test_with_tags").Key("agent_identity_blueprint_id").Exists(),
					check.That(resourceType+".test_with_tags").Key("service_principal_type").HasValue("ServiceIdentity"),
					check.That(resourceType+".test_with_tags").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_with_tags").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_with_tags").Key("sponsor_ids.#").HasValue("1"),
					check.That(resourceType+".test_with_tags").Key("owner_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent identity with tags")
				},
				ResourceName: resourceType + ".test_with_tags",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_with_tags"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_with_tags")
					}
					id := rs.Primary.Attributes["id"]
					blueprintId := rs.Primary.Attributes["agent_identity_blueprint_id"]
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s/%s:hard_delete=%s", id, blueprintId, hardDelete), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

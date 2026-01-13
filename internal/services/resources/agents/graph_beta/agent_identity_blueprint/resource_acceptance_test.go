package graphBetaApplicationsAgentIdentityBlueprint_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAgentIdentityBlueprintResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating minimal agent identity blueprint")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("agent identity blueprint", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-agent-identity-blueprint-minimal-[a-z0-9]+$`)),
					check.That(resourceType+".test_minimal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgentIdentityBlueprintResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating maximal agent identity blueprint")
				},
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("agent identity blueprint", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-agent-identity-blueprint-maximal-[a-z0-9]+$`)),
					check.That(resourceType+".test_maximal").Key("description").HasValue("This is a test agent identity blueprint with all optional fields configured"),
					check.That(resourceType+".test_maximal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("terraform"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("test"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("agent"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName: resourceType + ".test_maximal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_maximal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_maximal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgentIdentityBlueprintResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating agent identity blueprint for update test")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("agent identity blueprint", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-agent-identity-blueprint-minimal-`)),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating agent identity blueprint")
				},
				Config: testAccConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("agent identity blueprint", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-agent-identity-blueprint-minimal-updated-`)),
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description for acceptance test"),
				),
			},
		},
	})
}

func testAccConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfigMaximal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfigMinimalUpdated() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal_updated.tf")
	if err != nil {
		panic(err)
	}
	return content
}

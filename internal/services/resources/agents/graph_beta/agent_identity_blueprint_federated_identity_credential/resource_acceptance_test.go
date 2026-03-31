package graphBetaAgentIdentityBlueprintFederatedIdentityCredential_test

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

func TestAccResourceAgentIdentityBlueprintFederatedIdentityCredential_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating minimal federated identity credential")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 1 - Waiting for consistency")
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 2 - Validating id")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 3 - Validating blueprint_id")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("blueprint_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 4 - Validating name")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-fic-minimal-[a-z0-9]+$`)),
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 5 - Validating issuer")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 6 - Validating subject")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("subject").MatchesRegex(regexp.MustCompile(`^repo:deploymenttheory/test-repo-[a-z0-9]+:environment:Production$`)),
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 7 - Validating audiences")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("audiences.#").HasValue("1"),
					func(s *terraform.State) error {
						t.Log("--- Step 1: All checks passed")
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_minimal"),
			},
		},
	})
}

func TestAccResourceAgentIdentityBlueprintFederatedIdentityCredential_02_Update(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating federated identity credential for update test")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 1 - Waiting for consistency")
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					func(s *terraform.State) error {
						t.Log("--- Step 1: Check 2 - Validating id")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					func(s *terraform.State) error {
						t.Log("--- Step 1: All checks passed")
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating federated identity credential")
				},
				Config: testAccConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("--- Step 2: Check 1 - Waiting for consistency")
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					func(s *terraform.State) error {
						t.Log("--- Step 2: Check 2 - Validating description")
						return nil
					},
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description for acceptance test"),
					func(s *terraform.State) error {
						t.Log("--- Step 2: All checks passed")
						return nil
					},
				),
			},
		},
	})
}

func testAccImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		blueprintID := rs.Primary.Attributes["blueprint_id"]
		credentialID := rs.Primary.ID

		return fmt.Sprintf("%s/%s", blueprintID, credentialID), nil
	}
}

func testAccConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
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

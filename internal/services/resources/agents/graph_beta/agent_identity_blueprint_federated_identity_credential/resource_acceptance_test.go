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
	return `
resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-fic-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-fic-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-fic-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_2" {
  display_name        = "acc-test-fic-user2-${random_string.test_id.result}"
  user_principal_name = "acc-test-fic-user2-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-fic-user2-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_dependency" {
  display_name     = "acc-test-agent-identity-blueprint-fic-dependency-${random_string.test_id.result}"
  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
    microsoft365_graph_beta_users_user.dependency_user_2.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
    microsoft365_graph_beta_users_user.dependency_user_2.id,
  ]
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "test_minimal" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_dependency.id
  name         = "acc-test-fic-minimal-${random_string.test_id.result}"
  issuer       = "https://token.actions.githubusercontent.com"
  subject      = "repo:deploymenttheory/test-repo-${random_string.test_id.result}:environment:Production"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Updated description for acceptance test"
}
`
}

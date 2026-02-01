package graphBetaApplicationFederatedIdentityCredential_test

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

func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceApplicationFederatedIdentityCredential_01_Minimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Step 1: Creating minimal federated identity credential")
				},
				Config: loadAccTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("application_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-fic-minimal-[a-z0-9]+$`)),
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_minimal").Key("subject").MatchesRegex(regexp.MustCompile(`^repo:deploymenttheory/test-repo-[a-z0-9]+:environment:Production$`)),
					check.That(resourceType+".test_minimal").Key("audiences.#").HasValue("1"),
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

func TestAccResourceApplicationFederatedIdentityCredential_02_Maximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Step 1: Creating maximal federated identity credential")
				},
				Config: loadAccTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("application_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-fic-maximal-[a-z0-9]+$`)),
					check.That(resourceType+".test_maximal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_maximal").Key("subject").MatchesRegex(regexp.MustCompile(`^repo:deploymenttheory/test-repo-[a-z0-9]+:environment:Production$`)),
					check.That(resourceType+".test_maximal").Key("description").HasValue("Federated credential scenario - GitHub Actions with all optional fields configured"),
					check.That(resourceType+".test_maximal").Key("audiences.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_maximal"),
			},
		},
	})
}

func TestAccResourceApplicationFederatedIdentityCredential_03_Update(t *testing.T) {
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
				Config: loadAccTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating federated identity credential")
				},
				Config: testAccConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("federated identity credential", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description for acceptance test"),
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

		applicationID := rs.Primary.Attributes["application_id"]
		credentialID := rs.Primary.ID

		return fmt.Sprintf("%s/%s", applicationID, credentialID), nil
	}
}

func testAccConfigMinimalUpdated() string {
	return `
resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app" {
  display_name = "acc-test-app-fic-${random_string.test_id.result}"
  description  = "Application for federated identity credential acceptance test"
  hard_delete  = true
}

resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "test_minimal" {
  application_id = microsoft365_graph_beta_applications_application.test_app.id
  name           = "acc-test-fic-minimal-${random_string.test_id.result}"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:deploymenttheory/test-repo-${random_string.test_id.result}:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]
  description    = "Updated description for acceptance test"
}
`
}

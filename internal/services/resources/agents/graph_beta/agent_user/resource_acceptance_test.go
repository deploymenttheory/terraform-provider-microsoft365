package graphBetaAgentUser_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
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

// TestAccAgentUserResource_Minimal tests creating an agent user with minimal configuration
func TestAccAgentUserResource_Minimal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: ">= 0.9.0",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent user with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent user", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".test_minimal").Key("id").Exists(),
					check.That(resourceType+".test_minimal").Key("display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("agent_identity_id").Exists(),
					check.That(resourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("user_principal_name").Exists(),
					check.That(resourceType+".test_minimal").Key("mail_nickname").Exists(),
					check.That(resourceType+".test_minimal").Key("sponsor_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent user")
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
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccAgentUserResource_Maximal tests creating an agent user with all optional fields
func TestAccAgentUserResource_Maximal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: ">= 0.9.0",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating agent user with maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("agent user with all fields", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".test_maximal").Key("id").Exists(),
					check.That(resourceType+".test_maximal").Key("display_name").Exists(),
					check.That(resourceType+".test_maximal").Key("agent_identity_id").Exists(),
					check.That(resourceType+".test_maximal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("user_principal_name").Exists(),
					check.That(resourceType+".test_maximal").Key("mail_nickname").Exists(),
					check.That(resourceType+".test_maximal").Key("sponsor_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("given_name").HasValue("Agent"),
					check.That(resourceType+".test_maximal").Key("surname").HasValue("User"),
					check.That(resourceType+".test_maximal").Key("job_title").HasValue("AI Agent"),
					check.That(resourceType+".test_maximal").Key("department").HasValue("Engineering"),
					check.That(resourceType+".test_maximal").Key("company_name").HasValue("Contoso"),
					check.That(resourceType+".test_maximal").Key("office_location").HasValue("Building A"),
					check.That(resourceType+".test_maximal").Key("city").HasValue("Seattle"),
					check.That(resourceType+".test_maximal").Key("state").HasValue("WA"),
					check.That(resourceType+".test_maximal").Key("country").HasValue("US"),
					check.That(resourceType+".test_maximal").Key("postal_code").HasValue("98101"),
					check.That(resourceType+".test_maximal").Key("street_address").HasValue("123 Main Street"),
					check.That(resourceType+".test_maximal").Key("usage_location").HasValue("US"),
					check.That(resourceType+".test_maximal").Key("preferred_language").HasValue("en-US"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing agent user with all fields")
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
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

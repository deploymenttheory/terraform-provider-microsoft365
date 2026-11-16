package graphBetaIdentityAndAccessRoleDefinition_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIdentityAndAccessRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/role_definition"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessRoleDefinition.ResourceName

	// testResource is the test resource implementation for custom directory role definitions
	testResource = graphBetaIdentityAndAccessRoleDefinition.RoleDefinitionTestResource{}

	// externalProviders are the providers the acceptance configs use for unique names and template IDs
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// waitForConsistency gives Microsoft Entra replicas time to converge before checks use a separate client
func waitForConsistency(_ *terraform.State) error {
	testlog.WaitForConsistency("role definition", 30*time.Second)
	time.Sleep(30 * time.Second)
	return nil
}

// Minimal: required attributes and one permission block
func TestAccResourceRoleDefinition_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal custom directory role")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-role-definition-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("description").DoesNotExist(),
					check.That(resourceType+".minimal").Key("version").DoesNotExist(),
					check.That(resourceType+".minimal").Key("template_id").MatchesOtherKey(check.That(resourceType+".minimal").Key("id")),
					check.That(resourceType+".minimal").Key("is_built_in").HasValue("false"),
					check.That(resourceType+".minimal").Key("resource_scopes.*").ContainsTypeSetElement("/"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/applications/standard/read"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal custom directory role")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal: every optional attribute, a custom template ID and two permission blocks
func TestAccResourceRoleDefinition_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal custom directory role")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").MatchesOtherKey(check.That("random_uuid.template_id").Key("result")),
					check.That(resourceType+".maximal").Key("template_id").MatchesOtherKey(check.That("random_uuid.template_id").Key("result")),
					check.That(resourceType+".maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-role-definition-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("description").HasValue("Acceptance test custom role with every optional field set"),
					check.That(resourceType+".maximal").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("version").HasValue("1"),
					check.That(resourceType+".maximal").Key("role_permissions.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/groups/standard/read"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/bitlockerKeys/metadata/read"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal custom directory role")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Minimal to maximal and back: in-place updates, including clearing description and version
func TestAccResourceRoleDefinition_03_MinimalToMaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal custom directory role")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating custom directory role to maximal in place")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-role-definition-updated-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("description").HasValue("Acceptance test custom role updated in place"),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".minimal").Key("version").HasValue("2"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("2"),
					check.That(resourceType+".minimal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/devices/standard/read"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating custom directory role back to minimal in place")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-role-definition-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("description").DoesNotExist(),
					check.That(resourceType+".minimal").Key("version").DoesNotExist(),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing updated custom directory role")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Unknown action: rejected against the live microsoft.directory action catalogue before any write
func TestAccResourceRoleDefinition_04_ValidationErrorUnknownAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating custom directory role with an unknown action")
				},
				Config:      loadAcceptanceTestTerraform("resource_invalid_action.tf"),
				ExpectError: regexp.MustCompile(`resource actions not found in the microsoft\.directory namespace:\s+microsoft\.directory/applications/basic/read`),
			},
		},
	})
}

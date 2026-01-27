package graphBetaRoleDefinition_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_definition"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// testResource is the test resource implementation for role definitions
	testResource = graphBetaRoleDefinition.RoleDefinitionTestResource{}

	// resourceType is the Terraform resource type name
	resourceType = graphBetaRoleDefinition.ResourceName
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceRoleDefinition_01_Lifecycle(t *testing.T) {
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
			// Create with minimal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating role definition with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role definition", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-role-definition-minimal"),
					check.That(resourceType+".test").Key("description").HasValue(""),
					check.That(resourceType+".test").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".test").Key("is_built_in").Exists(),
					check.That(resourceType+".test").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.#").HasValue("2"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.*").ContainsTypeSetElement("Microsoft.Intune_ManagedDevices_Read"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.*").ContainsTypeSetElement("Microsoft.Intune_ManagedDevices_Update"),
				),
			},
			// ImportState testing
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing role definition")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating role definition to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role definition", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-role-definition-maximal"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That(resourceType+".test").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".test").Key("is_built_in").Exists(),
					check.That(resourceType+".test").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".test").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".test").Key("role_scope_tag_ids.*").ContainsTypeSetElement("1"),
					check.That(resourceType+".test").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.#").Exists(),
				),
			},
		},
	})
}

func TestAccResourceRoleDefinition_02_Description(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating role definition with description")
				},
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role definition", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".description").ExistsInGraph(testResource),
					check.That(resourceType+".description").Key("id").Exists(),
					check.That(resourceType+".description").Key("display_name").HasValue("acc-test-role-definition-description"),
					check.That(resourceType+".description").Key("description").HasValue("This is a test role definition with description"),
					check.That(resourceType+".description").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".description").Key("is_built_in").Exists(),
					check.That(resourceType+".description").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".description").Key("role_permissions.0.allowed_resource_actions.#").HasValue("1"),
				),
			},
		},
	})
}

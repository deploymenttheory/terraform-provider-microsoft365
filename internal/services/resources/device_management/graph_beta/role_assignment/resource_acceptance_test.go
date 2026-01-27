package graphBetaRoleDefinitionAssignment_test

import (
	"fmt"
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
	graphBetaRoleAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// resourceType is the Terraform resource type name
	resourceType = graphBetaRoleAssignment.ResourceName

	// testResource is the test resource implementation for role assignments
	testResource = graphBetaRoleAssignment.RoleAssignmentTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceRoleAssignment_01_Lifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating role assignment with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("display_name").Exists(),
					check.That(resourceType+".test").Key("role_definition_id").HasValue("0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					check.That(resourceType+".test").Key("members.#").HasValue("1"),
					check.That(resourceType+".test").Key("scope_configuration.#").HasValue("1"),
					check.That(resourceType+".test").Key("scope_configuration.0.type").HasValue("AllLicensedUsers"),
				),
			},
			// ImportState testing
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing role assignment")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("not found: %s.test", resourceType)
					}
					id := rs.Primary.ID
					roleDefId := rs.Primary.Attributes["role_definition_id"]
					compositeId := fmt.Sprintf("%s/%s", id, roleDefId)
					return compositeId, nil
				},
			},
			// Update to maximal configuration
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating role assignment to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").Exists(),
					check.That(resourceType+".test").Key("role_definition_id").HasValue("0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
					check.That(resourceType+".test").Key("scope_configuration.0.type").HasValue("AllDevices"),
				),
			},
		},
	})
}

func TestAccResourceRoleAssignment_02_ResourceScopes(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating role assignment with resource scopes")
				},
				Config: loadAcceptanceTestTerraform("resource_resource_scopes.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".resource_scopes").ExistsInGraph(testResource),
					check.That(resourceType+".resource_scopes").Key("id").Exists(),
					check.That(resourceType+".resource_scopes").Key("scope_configuration.0.type").HasValue("ResourceScopes"),
					check.That(resourceType+".resource_scopes").Key("scope_configuration.0.resource_scopes.#").HasValue("2"),
					check.That(resourceType+".resource_scopes").Key("members.#").HasValue("2"),
				),
			},
		},
	})
}

func TestAccResourceRoleAssignment_03_AllDevicesScope(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating role assignment with AllDevices scope")
				},
				Config: loadAcceptanceTestTerraform("resource_all_devices.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".all_devices").ExistsInGraph(testResource),
					check.That(resourceType+".all_devices").Key("id").Exists(),
					check.That(resourceType+".all_devices").Key("scope_configuration.0.type").HasValue("AllDevices"),
					check.That(resourceType+".all_devices").Key("members.#").HasValue("2"),
				),
			},
		},
	})
}

func TestAccResourceRoleAssignment_04_AllUsersScope(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating role assignment with AllLicensedUsers scope")
				},
				Config: loadAcceptanceTestTerraform("resource_all_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("role assignment", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".all_users").ExistsInGraph(testResource),
					check.That(resourceType+".all_users").Key("id").Exists(),
					check.That(resourceType+".all_users").Key("scope_configuration.0.type").HasValue("AllLicensedUsers"),
					check.That(resourceType+".all_users").Key("members.#").HasValue("2"),
				),
			},
		},
	})
}

package graphBetaRoleDefinition_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_role_definition"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testResource = graphBetaRoleDefinition.RoleDefinitionTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceCloudPcRoleDefinition_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("--- Creating role definition with minimal configuration...")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-cloud-pc-role-definition-minimal"),
					check.That(resourceType+".test").Key("description").HasValue(""),
					check.That(resourceType+".test").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".test").Key("is_built_in").Exists(),
					check.That(resourceType+".test").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					t.Log("--- Waiting 15s for resource to achieve eventual consistency...")
					time.Sleep(15 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					t.Log("--- Importing role definition...")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					t.Log("--- Updating role definition to maximal configuration...")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-cloud-pc-role-definition-maximal"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.#").HasValue("92"),
					check.That(resourceType+".test").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					t.Log("--- Waiting 15s for resource to achieve eventual consistency...")
					time.Sleep(15 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
			},
		},
	})
}

func TestAccResourceCloudPcRoleDefinition_02_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("--- Creating role definition with description...")
				},
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".description").Key("id").Exists(),
					check.That(resourceType+".description").Key("display_name").HasValue("acc-test-cloud-pc-role-definition-description"),
					check.That(resourceType+".description").Key("description").HasValue("This is a test role definition with description"),
				),
			},
			{
				PreConfig: func() {
					t.Log("--- Waiting 15s for resource to achieve eventual consistency...")
					time.Sleep(15 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
			},
		},
	})
}

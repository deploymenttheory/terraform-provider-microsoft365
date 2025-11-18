package graphBetaRoleScopeTag_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_scope_tag"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type names constructed from exported constants
	resourceType      = constants.PROVIDER_NAME + "_" + graphBetaRoleScopeTag.ResourceName
	groupResourceType = constants.PROVIDER_NAME + "_" + graphBetaGroup.ResourceName

	// testResource is the test resource implementation for role scope tags
	testResource = graphBetaRoleScopeTag.RoleScopeTagTestResource{}

	// groupTestResource is the test resource implementation for groups (used when testing dependencies)
	groupTestResource = graphBetaGroup.GroupTestResource{}
)

func TestAccRoleScopeTagResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			15*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: resourceType,
				TestResource: testResource,
			},
			destroy.ResourceTypeMapping{
				ResourceType: groupResourceType,
				TestResource: groupTestResource,
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccRoleScopeTagConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.test").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.test").Key("id").Exists(),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.test").Key("display_name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.test").Key("description").IsNotEmpty(),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.test").Key("is_built_in").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      "microsoft365_graph_beta_device_management_role_scope_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: testAccRoleScopeTagConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("Microsoft Entra ID", 60*time.Second)
						time.Sleep(60 * time.Second)
						return nil
					},
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.test", "description"),
				),
			},
		},
	})
}

func TestAccRoleScopeTagResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			0,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccRoleScopeTagConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.description").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.description").Key("id").Exists(),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.description").Key("display_name").IsNotEmpty(),
					check.That("microsoft365_graph_beta_device_management_role_scope_tag.description").Key("description").IsNotEmpty(),
				),
			},
		},
	})
}

func TestAccRoleScopeTagResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			15*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: resourceType,
				TestResource: testResource,
			},
			destroy.ResourceTypeMapping{
				ResourceType: groupResourceType,
				TestResource: groupTestResource,
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccRoleScopeTagConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("Microsoft Entra ID", 60*time.Second)
						time.Sleep(60 * time.Second)
						return nil
					},
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.assignments", "assignments.#", "2"),
				),
			},
		},
	})
}

// Test configuration functions
func testAccRoleScopeTagConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleScopeTagConfig_maximal() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

func testAccRoleScopeTagConfig_description() string {
	config := mocks.LoadTerraformConfigFile("resource_description.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccRoleScopeTagConfig_assignments() string {
	dependencies := mocks.LoadTerraformConfigFile("resource_dependencies.tf")
	config := mocks.LoadTerraformConfigFile("resource_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(dependencies + "\n" + config)
}

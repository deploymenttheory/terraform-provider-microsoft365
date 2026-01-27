package graphBetaRoleScopeTag_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_scope_tag"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType      = graphBetaRoleScopeTag.ResourceName
	groupResourceType = graphBetaGroup.ResourceName
	testResource      = graphBetaRoleScopeTag.RoleScopeTagTestResource{}
	groupTestResource = graphBetaGroup.GroupTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceRoleScopeTag_01_Lifecycle(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".test").Key("description").IsNotEmpty(),
					check.That(resourceType+".test").Key("is_built_in").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("Microsoft Entra ID", 60*time.Second)
						time.Sleep(60 * time.Second)
						return nil
					},
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "display_name"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "description"),
				),
			},
		},
	})
}

func TestAccResourceRoleScopeTag_02_Description(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".description").ExistsInGraph(testResource),
					check.That(resourceType+".description").Key("id").Exists(),
					check.That(resourceType+".description").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".description").Key("description").IsNotEmpty(),
				),
			},
		},
	})
}

func TestAccResourceRoleScopeTag_03_Assignments(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("Microsoft Entra ID", 60*time.Second)
						time.Sleep(60 * time.Second)
						return nil
					},
					resource.TestCheckResourceAttrSet(resourceType+".assignments", "id"),
					resource.TestCheckResourceAttrSet(resourceType+".assignments", "description"),
					resource.TestCheckResourceAttr(resourceType+".assignments", "assignments.#", "2"),
				),
			},
		},
	})
}

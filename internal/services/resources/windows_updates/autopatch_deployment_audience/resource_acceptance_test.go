package graphBetaWindowsUpdatesAutopatchDeploymentAudience_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	graphBetaWindowsAutopatchDeploymentAudience "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_audience"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// loadAcceptanceTestTerraform loads an acceptance test HCL config from the acceptance test directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

const resourceType = graphBetaWindowsAutopatchDeploymentAudience.ResourceName

var testResource = graphBetaWindowsAutopatchDeploymentAudience.WindowsUpdateDeploymentAudienceTestResource{}

// Test 001: Basic audience creation (no members or exclusions)
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestAccResourceWindowsUpdateDeploymentAudience_01_BasicAudience(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsAutopatchDeploymentAudience.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating basic audience with no members")
				},
				Config: loadAcceptanceTestTerraform("01_basic_audience.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("member_type").HasValue("azureADDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Audience with updatableAssetGroup members (Entra groups)
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x2)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestAccResourceWindowsUpdateDeploymentAudience_02_WithMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsAutopatchDeploymentAudience.ResourceName,
				TestResource: graphBetaWindowsAutopatchDeploymentAudience.WindowsUpdateDeploymentAudienceTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience with 2 updatableAssetGroup members")
				},
				Config: loadAcceptanceTestTerraform("02_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("member_type").HasValue("updatableAssetGroup"),
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				// The /members endpoint returns the enrolled devices within each group, not the
				// group objects themselves. Since the test groups are empty (no enrolled WU devices),
				// the API returns [] on import read. members/exclusions cannot be round-tripped
				// via import for updatableAssetGroup when backing groups contain no enrolled assets.
				ImportStateVerifyIgnore: []string{"member_type", "members", "exclusions"},
			},
		},
	})
}

// Test 003: Audience with updatableAssetGroup members and exclusions
//
// API calls exercised:
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x2, addExclusions x1)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestAccResourceWindowsUpdateDeploymentAudience_03_WithExclusions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsAutopatchDeploymentAudience.ResourceName,
				TestResource: graphBetaWindowsAutopatchDeploymentAudience.WindowsUpdateDeploymentAudienceTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience with 2 members and 1 exclusion")
				},
				Config: loadAcceptanceTestTerraform("03_with_exclusions.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members and exclusions", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("member_type").HasValue("updatableAssetGroup"),
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
					check.That(resourceType+".test").Key("exclusions.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"member_type", "members", "exclusions"},
			},
		},
	})
}

// Test 004: Lifecycle - add a member via updateAudience (diff-based update)
//
// Step 1: Create audience with 2 updatableAssetGroup members.
// Step 2: Add a third member - exercises the diff calculation in constructUpdateMembersRequest
//         which sends only the delta (addMembers=[group3]) to updateAudience.
//
// API calls exercised (step 1):
//   - POST /admin/windows/updates/deploymentAudiences
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x2)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//
// API calls exercised (step 2):
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (addMembers x1)
//   - GET  /admin/windows/updates/deploymentAudiences/{id}
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/members
//   - GET  /admin/windows/updates/deploymentAudiences/{id}/exclusions
//
// API calls exercised (destroy):
//   - POST /admin/windows/updates/deploymentAudiences/{id}/microsoft.graph.windowsUpdates.updateAudience (removeMembers x3)
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
func TestAccResourceWindowsUpdateDeploymentAudience_04_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsAutopatchDeploymentAudience.ResourceName,
				TestResource: graphBetaWindowsAutopatchDeploymentAudience.WindowsUpdateDeploymentAudienceTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience with 2 initial members")
				},
				Config: loadAcceptanceTestTerraform("04_lifecycle_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("member_type").HasValue("updatableAssetGroup"),
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Adding a third member via diff-based updateAudience")
				},
				Config: loadAcceptanceTestTerraform("04_lifecycle_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("member_type").HasValue("updatableAssetGroup"),
					check.That(resourceType+".test").Key("members.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 3: Import state verification")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"member_type", "members", "exclusions"},
			},
		},
	})
}

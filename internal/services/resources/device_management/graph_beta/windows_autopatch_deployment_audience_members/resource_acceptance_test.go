package graphBetaWindowsAutopatchDeploymentAudienceMembers_test

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
	WindowsAutopatchDeploymentResourceAudience "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopatch_deployment_audience"
	WindowsAutopatchDeploymentResourceAudienceMembers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopatch_deployment_audience_members"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to load acceptance test configs
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

const resourceType = WindowsAutopatchDeploymentResourceAudienceMembers.ResourceName

var testResource = WindowsAutopatchDeploymentResourceAudienceMembers.WindowsUpdateDeploymentAudienceMembersTestResource{}

// Test 001: Basic members with groups
func TestAccResourceWindowsUpdateDeploymentAudienceMembers_01_BasicMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudienceMembers.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudienceMembers.WindowsUpdateDeploymentAudienceMembersTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudience.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudience.WindowsUpdateDeploymentAudienceTestResource{},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience members with groups")
				},
				Config: loadAcceptanceTestTerraform("01_basic_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_updatableAssetGroup$`)),
					check.That(resourceType+".test").Key("audience_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
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
			},
		},
	})
}

// Test 002: Members with exclusions
func TestAccResourceWindowsUpdateDeploymentAudienceMembers_02_MembersWithExclusions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudienceMembers.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudienceMembers.WindowsUpdateDeploymentAudienceMembersTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudience.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudience.WindowsUpdateDeploymentAudienceTestResource{},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience members with exclusions")
				},
				Config: loadAcceptanceTestTerraform("02_members_with_exclusions.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_updatableAssetGroup$`)),
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
			},
		},
	})
}

// Test 003: Lifecycle - add and remove members
func TestAccResourceWindowsUpdateDeploymentAudienceMembers_03_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudienceMembers.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudienceMembers.WindowsUpdateDeploymentAudienceMembersTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: WindowsAutopatchDeploymentResourceAudience.ResourceName,
				TestResource: WindowsAutopatchDeploymentResourceAudience.WindowsUpdateDeploymentAudienceTestResource{},
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
					testlog.StepAction(resourceType, "Step 1: Creating audience with initial members")
				},
				Config: loadAcceptanceTestTerraform("03_lifecycle_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Adding more members")
				},
				Config: loadAcceptanceTestTerraform("03_lifecycle_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("audience members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
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
			},
		},
	})
}

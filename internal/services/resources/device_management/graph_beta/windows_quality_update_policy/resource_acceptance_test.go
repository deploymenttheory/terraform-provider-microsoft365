package graphBetaWindowsQualityUpdatePolicy_test

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
	graphBetaWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_policy"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ============================================================================
// Test Strategy & Timing Considerations
// ============================================================================
//
// These acceptance tests interact with Microsoft Graph API which has eventual
// consistency characteristics. For tests with assignments:
//
// - Groups with hard_delete=true require 30 seconds before CheckDestroy to
//   allow the two-phase deletion (soft → hard) to propagate through Microsoft
//   Graph's backend.
// ============================================================================

const resourceType = graphBetaWindowsQualityUpdatePolicy.ResourceName

var testResource = graphBetaWindowsQualityUpdatePolicy.WindowsQualityUpdatePolicyTestResource{}

// loadAcceptanceTestTerraform loads test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// loadAcceptanceTestTerraformWithDeps loads test config with standard device management dependencies
func loadAcceptanceTestTerraformWithDeps(filename string) string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic("failed to load groups config: " + err.Error())
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		panic("failed to load role scope tags config: " + err.Error())
	}
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + config)
}

// Test 001: Scenario 1 - Minimal configuration without assignments
func TestAccWindowsQualityUpdatePolicyResource_001_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsQualityUpdatePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("display_name").Exists(),
					check.That(resourceType+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Scenario 2 - Maximal configuration without assignments
func TestAccWindowsQualityUpdatePolicyResource_002_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsQualityUpdatePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("display_name").Exists(),
					check.That(resourceType+".test_002").Key("description").Exists(),
					check.That(resourceType+".test_002").Key("hotpatch_enabled").HasValue("true"),
					check.That(resourceType+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:      resourceType + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Scenario 3 - Lifecycle from minimal to maximal
func TestAccWindowsQualityUpdatePolicyResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsQualityUpdatePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("description").Exists(),
					check.That(resourceType+".test_003").Key("hotpatch_enabled").HasValue("true"),
					check.That(resourceType+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Scenario 4 - Lifecycle from maximal to minimal
func TestAccWindowsQualityUpdatePolicyResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsQualityUpdatePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("description").Exists(),
					check.That(resourceType+".test_004").Key("hotpatch_enabled").HasValue("true"),
					check.That(resourceType+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: Scenario 5 - Minimal assignments
func TestAccWindowsQualityUpdatePolicyResource_005_AssignmentsMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsQualityUpdatePolicy.ResourceName,
				TestResource: graphBetaWindowsQualityUpdatePolicy.WindowsQualityUpdatePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraformWithDeps("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_005").Key("display_name").Exists(),
					check.That(resourceType+".test_005").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 006: Scenario 6 - Maximal assignments
func TestAccWindowsQualityUpdatePolicyResource_006_AssignmentsMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsQualityUpdatePolicy.ResourceName,
				TestResource: graphBetaWindowsQualityUpdatePolicy.WindowsQualityUpdatePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraformWithDeps("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_006").Key("display_name").Exists(),
					check.That(resourceType+".test_006").Key("assignments.#").HasValue("3"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 007: Scenario 7 - Assignments lifecycle minimal to maximal
func TestAccWindowsQualityUpdatePolicyResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsQualityUpdatePolicy.ResourceName,
				TestResource: graphBetaWindowsQualityUpdatePolicy.WindowsQualityUpdatePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraformWithDeps("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("display_name").Exists(),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadAcceptanceTestTerraformWithDeps("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("display_name").Exists(),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("3"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_007",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 008: Scenario 8 - Assignments lifecycle maximal to minimal
func TestAccWindowsQualityUpdatePolicyResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsQualityUpdatePolicy.ResourceName,
				TestResource: graphBetaWindowsQualityUpdatePolicy.WindowsQualityUpdatePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraformWithDeps("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("3"),
				),
			},
			{
				Config: loadAcceptanceTestTerraformWithDeps("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows quality update policy assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      resourceType + ".test_008",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

package graphBetaWindowsRemediationScript_test

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
	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/assignment_filter"
	graphBetaWindowsRemediationScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_remediation_script"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ============================================================================
// Test Strategy & Timing Considerations
// ============================================================================
//
// These acceptance tests interact with Microsoft Graph API which has eventual
// consistency characteristics that require strategic timing to ensure reliable
// test execution. The following explains the test flow and why specific waits
// are necessary:
//
// ## Multi-Step Lifecycle Tests (007, 008)
//
// For tests that modify assignments across multiple steps:
//
// 1. ✅ Complete Step 1 with initial assignment configuration
// 2. ⏰ Wait 20 seconds for consistency
//    WHY: When assignments are added/removed from a script, the Microsoft
//    Graph API needs time to propagate these changes. Without this wait,
//    Terraform's refresh between steps may detect false drift (resources
//    appearing as deleted when they still exist). This prevents spurious
//    "resource has been deleted outside of Terraform" errors.
//
// 3. ✅ Complete Step 2 with modified assignments (no drift errors)
// 4. ⏰ Wait 60 seconds before CheckDestroy
//    WHY: Groups with hard_delete=true require a two-phase deletion:
//    - Soft delete: Resource moves to "deleted items" collection
//    - Hard delete: Resource is permanently removed
//    The hard delete operation can take 60-90 seconds to fully propagate
//    through Microsoft Graph's backend. Without this wait, CheckDestroy
//    may find resources still in the deleted items collection, causing
//    false test failures.
//
// 5. ✅ Verify all resources are properly destroyed
//
// ## Assignment Filter Delete Timing
//
// Assignment filters have an additional 10-second pre-delete wait:
//
// WHY: Assignment filters are locked when actively referenced by remediation
// script assignments. Even after the script is deleted, there's a brief
// window where the filter remains locked in the backend. The 10-second pause
// allows the backend to release the lock before attempting deletion, preventing
// 500 Internal Server Error responses.
//
// ## Idempotent Delete Operations
//
// All delete operations (groups, filters, scripts) are idempotent:
//
// WHY: Due to eventual consistency and test retries, resources may already be
// deleted when Terraform attempts to delete them (404 Not Found). Treating
// 404 as success ensures tests can be re-run without manual cleanup and
// handles backend timing variations gracefully.
// ============================================================================

// Helper function to load acceptance test configs
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

const resourceType = graphBetaWindowsRemediationScript.ResourceName

var testResource = graphBetaWindowsRemediationScript.WindowsRemediationScriptTestResource{}

// Test 001: Scenario 1 - Minimal configuration without assignments
func TestAccResourceWindowsRemediationScript_01_Scenario_Minimal(t *testing.T) {
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
			graphBetaWindowsRemediationScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-remediation-script-001-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".test_001").Key("description").HasValue("Scenario 1: Minimal configuration without assignments"),
					check.That(resourceType+".test_001").Key("publisher").HasValue("Terraform Provider Test"),
					check.That(resourceType+".test_001").Key("run_as_account").HasValue("system"),
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
func TestAccResourceWindowsRemediationScript_02_Scenario_Maximal(t *testing.T) {
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
			graphBetaWindowsRemediationScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-remediation-script-002-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".test_002").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_002").Key("run_as_32_bit").HasValue("true"),
					check.That(resourceType+".test_002").Key("enforce_signature_check").HasValue("true"),
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
func TestAccResourceWindowsRemediationScript_03_Lifecycle_MinimalToMaximal(t *testing.T) {
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
			graphBetaWindowsRemediationScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("run_as_account").HasValue("system"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_003").Key("run_as_32_bit").HasValue("true"),
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
func TestAccResourceWindowsRemediationScript_04_Lifecycle_MaximalToMinimal(t *testing.T) {
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
			graphBetaWindowsRemediationScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("run_as_account").HasValue("user"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("run_as_account").HasValue("system"),
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
func TestAccResourceWindowsRemediationScript_05_AssignmentsMinimal(t *testing.T) {
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
			graphBetaWindowsRemediationScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_005").Key("assignments.#").HasValue("1"),
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
func TestAccResourceWindowsRemediationScript_06_AssignmentsMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			60*time.Second, // Increased wait time for groups hard delete to propagate (can take 60-90s)
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsRemediationScript.ResourceName,
				TestResource: graphBetaWindowsRemediationScript.WindowsRemediationScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaAssignmentFilter.ResourceName,
				TestResource: graphBetaAssignmentFilter.AssignmentFilterTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_006").Key("assignments.#").HasValue("5"),
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
func TestAccResourceWindowsRemediationScript_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
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
			60*time.Second, // Increased wait time for groups hard delete to propagate (can take 60-90s)
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsRemediationScript.ResourceName,
				TestResource: graphBetaWindowsRemediationScript.WindowsRemediationScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaAssignmentFilter.ResourceName,
				TestResource: graphBetaAssignmentFilter.AssignmentFilterTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows remediation script assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("5"),
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
func TestAccResourceWindowsRemediationScript_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
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
			60*time.Second, // Increased wait time for groups hard delete to propagate (can take 60-90s)
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaWindowsRemediationScript.ResourceName,
				TestResource: graphBetaWindowsRemediationScript.WindowsRemediationScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaAssignmentFilter.ResourceName,
				TestResource: graphBetaAssignmentFilter.AssignmentFilterTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("5"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows remediation script assignments", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("1"),
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

// Test 009: Scenario 9 - Validation errors
// testAccCheckWindowsRemediationScriptDestroy verifies that Windows remediation scripts have been destroyed

package graphBetaMacOSPlatformScript_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMacOSPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_platform_script"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaMacOSPlatformScript.ResourceName

	// testResource is the test resource implementation for macOS platform scripts
	testResource = graphBetaMacOSPlatformScript.MacOSPlatformScriptTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Scenario 01: Minimal macOS Platform Script
func TestAccMacOSPlatformScriptResource_Scenario01_Minimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating minimal macOS platform script")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_01_minimal_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("acc-test-minimal-macos-script"),
					check.That(resourceType+".minimal").Key("file_name").HasValue("minimal_test.sh"),
					check.That(resourceType+".minimal").Key("run_as_account").HasValue("system"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal macOS platform script")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 02: Maximal macOS Platform Script
func TestAccMacOSPlatformScriptResource_Scenario02_Maximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating maximal macOS platform script with all features")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_02_maximal_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("display_name").HasValue("acc-test-maximal-macos-script"),
					check.That(resourceType+".maximal").Key("description").HasValue("Comprehensive macOS platform script with all features enabled for unit testing"),
					check.That(resourceType+".maximal").Key("file_name").HasValue("maximal_test.sh"),
					check.That(resourceType+".maximal").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("1"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("2"),
					check.That(resourceType+".maximal").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".maximal").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".maximal").Key("retry_count").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal macOS platform script")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 03: Minimal to Maximal Update
func TestAccMacOSPlatformScriptResource_Scenario03_MinimalToMaximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating minimal macOS platform script for update test")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("acc-test-update-test-script"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test.sh"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("system"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("acc-test-update-test-script-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test_maximal.sh"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".update_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".update_test").Key("execution_frequency").HasValue("PT12H"),
					check.That(resourceType+".update_test").Key("retry_count").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing updated macOS platform script")
				},
				ResourceName:            resourceType + ".update_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 04: Maximal to Minimal Update
func TestAccMacOSPlatformScriptResource_Scenario04_MaximalToMinimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating maximal macOS platform script for downgrade test")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".downgrade_test").ExistsInGraph(testResource),
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("display_name").HasValue("acc-test-downgrade-test-script"),
					check.That(resourceType+".downgrade_test").Key("description").HasValue("Initial maximal configuration for downgrade testing"),
					check.That(resourceType+".downgrade_test").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".downgrade_test").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(resourceType+".downgrade_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".downgrade_test").Key("retry_count").HasValue("5"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Downgrading to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".downgrade_test").ExistsInGraph(testResource),
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("display_name").HasValue("acc-test-downgrade-test-script-minimal"),
					check.That(resourceType+".downgrade_test").Key("file_name").HasValue("downgrade_test_minimal.sh"),
					check.That(resourceType+".downgrade_test").Key("run_as_account").HasValue("system"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing downgraded macOS platform script")
				},
				ResourceName:            resourceType + ".downgrade_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 05: No Assignments to Minimal Assignment
func TestAccMacOSPlatformScriptResource_Scenario05_NoAssignmentsToMinimalAssignment(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating macOS platform script without assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".add_minimal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".add_minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_minimal_assignment").Key("display_name").HasValue("acc-test-add-minimal-assignment"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adding single assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".add_minimal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".add_minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_minimal_assignment").Key("display_name").HasValue("acc-test-add-minimal-assignment"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_minimal_assignment", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing macOS platform script with assignment")
				},
				ResourceName:            resourceType + ".add_minimal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 06: No Assignments to Maximal Assignments
func TestAccMacOSPlatformScriptResource_Scenario06_NoAssignmentsToMaximalAssignments(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating macOS platform script without assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".add_maximal_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".add_maximal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_maximal_assignments").Key("display_name").HasValue("acc-test-add-maximal-assignments"),
					check.That(resourceType+".add_maximal_assignments").Key("assignments.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adding all 4 assignment types")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".add_maximal_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".add_maximal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_maximal_assignments").Key("display_name").HasValue("acc-test-add-maximal-assignments"),
					check.That(resourceType+".add_maximal_assignments").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing macOS platform script with maximal assignments")
				},
				ResourceName:            resourceType + ".add_maximal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 07: Minimal to Maximal Assignments
func TestAccMacOSPlatformScriptResource_Scenario07_MinimalToMaximalAssignments(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating macOS platform script with single assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_update").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_update").Key("display_name").HasValue("acc-test-assignment-update"),
					check.That(resourceType+".assignment_update").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to all 4 assignment types")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_update").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_update").Key("display_name").HasValue("acc-test-assignment-update"),
					check.That(resourceType+".assignment_update").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing macOS platform script with updated assignments")
				},
				ResourceName:            resourceType + ".assignment_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

// Scenario 08: Maximal to Minimal Assignments
func TestAccMacOSPlatformScriptResource_Scenario08_MaximalToMinimalAssignments(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating macOS platform script with all 4 assignment types")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_downgrade").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_downgrade").Key("display_name").HasValue("acc-test-assignment-downgrade"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Downgrading to single assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_downgrade").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_downgrade").Key("display_name").HasValue("acc-test-assignment-downgrade"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing macOS platform script with downgraded assignments")
				},
				ResourceName:            resourceType + ".assignment_downgrade",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "script_content"},
			},
		},
	})
}

package graphBetaMacOSPlatformScript_test

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
	graphBetaMacOSPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_platform_script"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
func TestAccResourceMacOSPlatformScript_01_Minimal(t *testing.T) {
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
					check.That(resourceType+".minimal").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
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
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 02: Maximal macOS Platform Script
func TestAccResourceMacOSPlatformScript_02_Maximal(t *testing.T) {
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
					check.That(resourceType+".maximal").Key("script_content").HasValue("#!/bin/bash\necho 'Max Test'\nexit 0"),
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
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 03: Minimal to Maximal Update
func TestAccResourceMacOSPlatformScript_03_MinimalToMaximal(t *testing.T) {
	var scriptID string
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
				Config:           loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".update_test", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".update_test", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".update_test").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test.sh"),
					check.That(resourceType+".update_test").Key("execution_frequency").HasValue("PT0S"),
					check.That(resourceType+".update_test").Key("block_execution_notifications").HasValue("false"),
					check.That(resourceType+".update_test").Key("retry_count").HasValue("0"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("acc-test-update-test-script"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("system"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".update_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to maximal configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".update_test", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".update_test", "id", &scriptID),
					check.That(resourceType+".update_test").Key("script_content").HasValue("#!/bin/bash\necho 'Max Test'\nexit 0"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test_maximal.sh"),
					check.That(resourceType+".update_test").Key("execution_frequency").HasValue("PT12H"),
					check.That(resourceType+".update_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".update_test").Key("retry_count").HasValue("2"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("acc-test-update-test-script-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".update_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Returning to initial configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".update_test", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".update_test", "id", &scriptID),
					check.That(resourceType+".update_test").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test.sh"),
					check.That(resourceType+".update_test").Key("execution_frequency").HasValue("PT0S"),
					check.That(resourceType+".update_test").Key("block_execution_notifications").HasValue("false"),
					check.That(resourceType+".update_test").Key("retry_count").HasValue("0"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("acc-test-update-test-script"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("system"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".update_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"), PlanOnly: true},
		},
	})
}

// Scenario 04: Maximal to Minimal Update
func TestAccResourceMacOSPlatformScript_04_MaximalToMinimal(t *testing.T) {
	var scriptID string
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
				Config:           loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".downgrade_test", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".downgrade_test", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".downgrade_test").Key("script_content").HasValue("#!/bin/bash\necho 'Max Test'\nexit 0"),
					check.That(resourceType+".downgrade_test").Key("file_name").HasValue("downgrade_test.sh"),
					check.That(resourceType+".downgrade_test").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".downgrade_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("retry_count").HasValue("5"),
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
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".downgrade_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Downgrading to minimal configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".downgrade_test", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".downgrade_test", "id", &scriptID),
					check.That(resourceType+".downgrade_test").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".downgrade_test").Key("file_name").HasValue("downgrade_test_minimal.sh"),
					check.That(resourceType+".downgrade_test").Key("execution_frequency").HasValue("PT0S"),
					check.That(resourceType+".downgrade_test").Key("block_execution_notifications").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("retry_count").HasValue("0"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".downgrade_test").ExistsInGraph(testResource),
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("display_name").HasValue("acc-test-downgrade-test-script-minimal"),
					check.That(resourceType+".downgrade_test").Key("run_as_account").HasValue("system"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".downgrade_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Restoring maximal configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".downgrade_test", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".downgrade_test", "id", &scriptID),
					check.That(resourceType+".downgrade_test").Key("script_content").HasValue("#!/bin/bash\necho 'Max Test'\nexit 0"),
					check.That(resourceType+".downgrade_test").Key("file_name").HasValue("downgrade_test.sh"),
					check.That(resourceType+".downgrade_test").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".downgrade_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("retry_count").HasValue("5"),
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
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".downgrade_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"), PlanOnly: true},
		},
	})
}

// Scenario 05: No Assignments to Minimal Assignment
func TestAccResourceMacOSPlatformScript_05_NoAssignmentsToMinimalAssignment(t *testing.T) {
	var scriptID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: graphBetaGroup.ResourceName, TestResource: graphBetaGroup.GroupTestResource{}},
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
					testlog.StepAction(resourceType, "Creating macOS platform script without assignments")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_minimal_assignment", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".add_minimal_assignment", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".add_minimal_assignment").Key("script_content").HasValue("#!/bin/bash\necho 'Test'\nexit 0"),
					check.That(resourceType+".add_minimal_assignment").Key("file_name").HasValue("test_minimal_assignment.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_minimal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adding single assignment")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_minimal_assignment", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".add_minimal_assignment", "id", &scriptID),
					check.That(resourceType+".add_minimal_assignment").Key("script_content").HasValue("#!/bin/bash\necho 'Test'\nexit 0"),
					check.That(resourceType+".add_minimal_assignment").Key("file_name").HasValue("test_minimal_assignment.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_minimal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Returning to initial configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_minimal_assignment", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".add_minimal_assignment", "id", &scriptID),
					check.That(resourceType+".add_minimal_assignment").Key("script_content").HasValue("#!/bin/bash\necho 'Test'\nexit 0"),
					check.That(resourceType+".add_minimal_assignment").Key("file_name").HasValue("test_minimal_assignment.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_minimal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"), PlanOnly: true},
		},
	})
}

// Scenario 06: No Assignments to Maximal Assignments
func TestAccResourceMacOSPlatformScript_06_NoAssignmentsToMaximalAssignments(t *testing.T) {
	var scriptID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: graphBetaGroup.ResourceName, TestResource: graphBetaGroup.GroupTestResource{}},
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
					testlog.StepAction(resourceType, "Creating macOS platform script without assignments")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_maximal_assignments", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".add_maximal_assignments", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".add_maximal_assignments").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".add_maximal_assignments").Key("file_name").HasValue("test_maximal_assignments.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_maximal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adding all 4 assignment types")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_maximal_assignments", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".add_maximal_assignments", "id", &scriptID),
					check.That(resourceType+".add_maximal_assignments").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".add_maximal_assignments").Key("file_name").HasValue("test_maximal_assignments.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_maximal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Returning to initial configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".add_maximal_assignments", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".add_maximal_assignments", "id", &scriptID),
					check.That(resourceType+".add_maximal_assignments").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".add_maximal_assignments").Key("file_name").HasValue("test_maximal_assignments.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".add_maximal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"), PlanOnly: true},
		},
	})
}

// Scenario 07: Minimal to Maximal Assignments
func TestAccResourceMacOSPlatformScript_07_MinimalToMaximalAssignments(t *testing.T) {
	var scriptID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: graphBetaGroup.ResourceName, TestResource: graphBetaGroup.GroupTestResource{}},
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
					testlog.StepAction(resourceType, "Creating macOS platform script with single assignment")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_update", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".assignment_update", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".assignment_update").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_update").Key("file_name").HasValue("test_assignment_update.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to all 4 assignment types")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_update", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".assignment_update", "id", &scriptID),
					check.That(resourceType+".assignment_update").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_update").Key("file_name").HasValue("test_assignment_update.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Returning to initial configuration")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_update", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".assignment_update", "id", &scriptID),
					check.That(resourceType+".assignment_update").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_update").Key("file_name").HasValue("test_assignment_update.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"), PlanOnly: true},
		},
	})
}

// Scenario 08: Maximal to Minimal Assignments
func TestAccResourceMacOSPlatformScript_08_MaximalToMinimalAssignments(t *testing.T) {
	var scriptID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: graphBetaGroup.ResourceName, TestResource: graphBetaGroup.GroupTestResource{}},
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
					testlog.StepAction(resourceType, "Creating macOS platform script with all 4 assignment types")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_downgrade", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".assignment_downgrade", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".assignment_downgrade").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_downgrade").Key("file_name").HasValue("test_assignment_downgrade.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_downgrade",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Downgrading to single assignment")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_02_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_downgrade", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".assignment_downgrade", "id", &scriptID),
					check.That(resourceType+".assignment_downgrade").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_downgrade").Key("file_name").HasValue("test_assignment_downgrade.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_downgrade",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_02_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing all assignments")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_03_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_downgrade", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".assignment_downgrade", "id", &scriptID),
					check.That(resourceType+".assignment_downgrade").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_downgrade").Key("file_name").HasValue("test_assignment_downgrade.sh"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macos platform script", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_downgrade").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_downgrade").Key("display_name").HasValue("acc-test-assignment-downgrade"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_downgrade",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_03_001.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Restoring maximal assignments")
				},
				Config:           loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".assignment_downgrade", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".assignment_downgrade", "id", &scriptID),
					check.That(resourceType+".assignment_downgrade").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".assignment_downgrade").Key("file_name").HasValue("test_assignment_downgrade.sh"),
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
					testlog.StepAction(resourceType, "Verifying import")
				},
				ResourceName:            resourceType + ".assignment_downgrade",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"), PlanOnly: true},
		},
	})
}

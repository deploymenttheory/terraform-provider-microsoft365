package graphBetaAppControlForBusinessBuiltInControls_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAppControlForBusinessBuiltInControls "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testResource = graphBetaAppControlForBusinessBuiltInControls.AppControlForBusinessBuiltInControlsTestResource{}

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// Test 001: Audit Mode
func TestAccAppControlForBusinessBuiltInControlsResource_001_AuditMode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating app control policy in audit mode")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_audit_mode.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-app-control-audit-mode-[a-z0-9]{8}$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("enable_app_control").HasValue("audit"),
				),
			},
		},
	})
}

// Test 002: Enforce Mode
func TestAccAppControlForBusinessBuiltInControlsResource_002_EnforceMode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating app control policy in enforce mode")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_enforce_mode.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-app-control-enforce-mode-[a-z0-9]{8}$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("enable_app_control").HasValue("enforce"),
				),
			},
		},
	})
}

// Test 003: Minimal Configuration
func TestAccAppControlForBusinessBuiltInControlsResource_003_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating minimal app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-app-control-for-business-built-in-controls-minimal-[a-z0-9]{8}$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("enable_app_control").HasValue("audit"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("role_scope_tag_ids.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Importing minimal app control policy")
				},
				ResourceName:            graphBetaAppControlForBusinessBuiltInControls.ResourceName + ".advanced",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 004: Maximal Configuration with Additional Rules
func TestAccAppControlForBusinessBuiltInControlsResource_004_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating maximal app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-app-control-for-business-built-in-controls-maximal-[a-z0-9]{8}$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("assignments.#").HasValue("3"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("assignments.0.type").HasValue("allLicensedUsersAssignmentTarget"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("assignments.1.type").HasValue("groupAssignmentTarget"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("assignments.2.type").HasValue("allDevicesAssignmentTarget"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Importing maximal app control policy")
				},
				ResourceName:            graphBetaAppControlForBusinessBuiltInControls.ResourceName + ".advanced",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 005: Lifecycle - Minimal to Maximal
func TestAccAppControlForBusinessBuiltInControlsResource_005_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating minimal lifecycle app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("enable_app_control").HasValue("audit"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Updating to maximal lifecycle app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 006: Lifecycle - Maximal to Minimal (Downgrade)
func TestAccAppControlForBusinessBuiltInControlsResource_006_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating maximal downgrade app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Downgrading to minimal app control policy")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 007: Assignments Lifecycle - Minimal to Maximal
func TestAccAppControlForBusinessBuiltInControlsResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating app control policy with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_assignments_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("assignments.#").HasValue("1"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("assignments.0.type").HasValue("allLicensedUsersAssignmentTarget"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Updating to maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_assignments_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("assignments.#").HasValue("3"),
				),
			},
		},
	})
}

// Test 008: Assignments Lifecycle - Maximal to Minimal (Downgrade)
func TestAccAppControlForBusinessBuiltInControlsResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaAppControlForBusinessBuiltInControls.ResourceName,
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
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Creating app control policy with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_assignments_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaAppControlForBusinessBuiltInControls.ResourceName, "Downgrading to minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_acfb_built_in_controls_assignments_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("App control policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").ExistsInGraph(testResource),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("assignments.#").HasValue("1"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("assignments.0.type").HasValue("allLicensedUsersAssignmentTarget"),
				),
			},
		},
	})
}

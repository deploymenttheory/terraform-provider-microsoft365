package graphBetaGroupPolicyDefinition_test

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
	graphBetaGroupPolicyDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_definition"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ============================================================================
// Test Strategy & Timing Considerations
// ============================================================================
//
// These acceptance tests interact with Microsoft Graph API which has eventual
// consistency characteristics that require strategic timing to ensure reliable
// test execution.
//
// ## Group Policy Definition Testing
//
// Group policy definitions are managed through the updateDefinitionValues
// endpoint and have dependencies on:
// - Group Policy Configurations (parent resource)
// - Group Policy Definitions (templates in Microsoft Graph)
// - Group Policy Presentations (individual settings with different types)
//
// ## Presentation Types Tested
//
// This resource supports multiple presentation types, each tested individually:
// 1. CheckBox (Boolean) - true/false values
// 2. TextBox (Text) - string values
// 3. DecimalTextBox (Numeric) - integer values
// 4. MultiTextBox (Multi-line) - newline-separated strings
// 5. DropdownList (Select) - predefined option values
//
// ## Resource Dependencies
//
// The test creates a group policy configuration first, then creates definitions
// within that configuration. During cleanup, the definition is deleted first,
// followed by the configuration.
//
// ## Timing Considerations
//
// 1. **Read After Write:** After creating or updating a definition, there's
//    a brief delay before the value is consistently readable from the API.
//    The ReadWithRetry mechanism handles this automatically.
//
// 2. **CheckDestroy Wait:** A 30-second wait is used before CheckDestroy to
//    ensure the resource deletion has propagated through the backend.
//
// 3. **Type Transitions:** When changing presentation types between test steps,
//    a brief delay ensures the previous state is fully committed before the
//    new configuration is applied.
//
// ============================================================================

// Helper function to load acceptance test configs
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

const resourceType = graphBetaGroupPolicyDefinition.ResourceName

var testResource = graphBetaGroupPolicyDefinition.GroupPolicyDefinitionTestResource{}

// Test 001: Boolean (CheckBox) - Minimal configuration
func TestAccGroupPolicyDefinitionResource_001_Boolean_Minimal(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating Boolean Minimal Configuration")
				},
				Config: loadAcceptanceTestTerraform("001_boolean_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").Exists(),
					check.That(resourceType+".test_001").Key("group_policy_configuration_id").Exists(),
					check.That(resourceType+".test_001").Key("policy_name").HasValue("Remove Default Microsoft Store packages from the system."),
					check.That(resourceType+".test_001").Key("class_type").HasValue("machine"),
					check.That(resourceType+".test_001").Key("category_path").HasValue("\\Windows Components\\App Package Deployment"),
					check.That(resourceType+".test_001").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_001").Key("values.#").HasValue("2"),
					// Verify IDs are populated
					check.That(resourceType+".test_001").Key("values.0.id").Exists(),
					check.That(resourceType+".test_001").Key("values.1.id").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing Boolean Minimal Configuration")
				},
				ResourceName:      resourceType + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Boolean (CheckBox) - Maximal configuration (25 values)
func TestAccGroupPolicyDefinitionResource_002_Boolean_Maximal(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating Boolean Maximal Configuration (25 values)")
				},
				Config: loadAcceptanceTestTerraform("002_boolean_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").Exists(),
					check.That(resourceType+".test_002").Key("policy_name").HasValue("Remove Default Microsoft Store packages from the system."),
					check.That(resourceType+".test_002").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_002").Key("values.#").HasValue("25"),
					// Verify all IDs are populated
					check.That(resourceType+".test_002").Key("values.0.id").Exists(),
					check.That(resourceType+".test_002").Key("values.24.id").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing Boolean Maximal Configuration")
				},
				ResourceName:      resourceType + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: TextBox - Single value
func TestAccGroupPolicyDefinitionResource_003_TextBox(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating TextBox Configuration")
				},
				Config: loadAcceptanceTestTerraform("003_textbox.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").Exists(),
					check.That(resourceType+".test_003").Key("policy_name").HasValue("Browsing Data Lifetime Settings"),
					check.That(resourceType+".test_003").Key("class_type").HasValue("machine"),
					check.That(resourceType+".test_003").Key("category_path").HasValue("\\Microsoft Edge"),
					check.That(resourceType+".test_003").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_003").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_003").Key("values.0.id").Exists(),
					check.That(resourceType+".test_003").Key("values.0.label").HasValue("Browsing Data Lifetime Settings"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing TextBox Configuration")
				},
				ResourceName:      resourceType + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: DecimalTextBox - Numeric value
func TestAccGroupPolicyDefinitionResource_004_Decimal(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating DecimalTextBox Configuration")
				},
				Config: loadAcceptanceTestTerraform("004_decimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").Exists(),
					check.That(resourceType+".test_004").Key("policy_name").HasValue("Configure time out for detections in non-critical failed state"),
					check.That(resourceType+".test_004").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_004").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_004").Key("values.0.id").Exists(),
					check.That(resourceType+".test_004").Key("values.0.value").HasValue("7200"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing DecimalTextBox Configuration")
				},
				ResourceName:      resourceType + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: MultiTextBox - Multi-line value
func TestAccGroupPolicyDefinitionResource_005_MultiText(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating MultiTextBox Configuration")
				},
				Config: loadAcceptanceTestTerraform("005_multitext.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").Key("id").Exists(),
					check.That(resourceType+".test_005").Key("policy_name").HasValue("Dev drive filter attach policy"),
					check.That(resourceType+".test_005").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_005").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_005").Key("values.0.id").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing MultiTextBox Configuration")
				},
				ResourceName:      resourceType + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 006: DropdownList - Select option
func TestAccGroupPolicyDefinitionResource_006_Dropdown(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Creating DropdownList Configuration")
				},
				Config: loadAcceptanceTestTerraform("006_dropdown.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_006").Key("id").Exists(),
					check.That(resourceType+".test_006").Key("policy_name").HasValue("Navigate windows and frames across different domains"),
					check.That(resourceType+".test_006").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_006").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_006").Key("values.0.id").Exists(),
					check.That(resourceType+".test_006").Key("values.0.value").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Importing DropdownList Configuration")
				},
				ResourceName:      resourceType + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 007: Lifecycle - Type transitions through all presentation types
func TestAccGroupPolicyDefinitionResource_007_Lifecycle_TypeTransitions(t *testing.T) {
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
			graphBetaGroupPolicyDefinition.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			// Step 1: Start with Boolean
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Step 1: Creating Boolean Configuration")
				},
				Config: loadAcceptanceTestTerraform("007_lifecycle_step_1_boolean.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").Exists(),
					check.That(resourceType+".test_007").Key("policy_name").HasValue("Remove Default Microsoft Store packages from the system."),
					check.That(resourceType+".test_007").Key("values.#").HasValue("2"),
				),
			},
			// Step 2: Transition to TextBox
			{
				PreConfig: func() {
					time.Sleep(3 * time.Second) // Pause to avoid API throttling between lifecycle transitions
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Step 2: Transitioning to TextBox (Destroy + Create)")
				},
				Config: loadAcceptanceTestTerraform("007_lifecycle_step_2_textbox.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").Exists(),
					check.That(resourceType+".test_007").Key("policy_name").HasValue("Browsing Data Lifetime Settings"),
					check.That(resourceType+".test_007").Key("values.#").HasValue("1"),
				),
			},
			// Step 3: Transition to Decimal
			{
				PreConfig: func() {
					time.Sleep(3 * time.Second) // Pause to avoid API throttling between lifecycle transitions
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Step 3: Transitioning to Decimal (Destroy + Create)")
				},
				Config: loadAcceptanceTestTerraform("007_lifecycle_step_3_decimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").Exists(),
					check.That(resourceType+".test_007").Key("policy_name").HasValue("Configure time out for detections in non-critical failed state"),
					check.That(resourceType+".test_007").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_007").Key("values.0.value").HasValue("7200"),
				),
			},
			// Step 4: Transition to MultiText
			{
				PreConfig: func() {
					time.Sleep(3 * time.Second) // Pause to avoid API throttling between lifecycle transitions
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Step 4: Transitioning to MultiText (Destroy + Create)")
				},
				Config: loadAcceptanceTestTerraform("007_lifecycle_step_4_multitext.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").Exists(),
					check.That(resourceType+".test_007").Key("policy_name").HasValue("Dev drive filter attach policy"),
					check.That(resourceType+".test_007").Key("values.#").HasValue("1"),
				),
			},
			// Step 5: Transition to Dropdown
			{
				PreConfig: func() {
					time.Sleep(3 * time.Second) // Pause to avoid API throttling between lifecycle transitions
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Step 5: Transitioning to Dropdown (Destroy + Create)")
				},
				Config: loadAcceptanceTestTerraform("007_lifecycle_step_5_dropdown.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").Exists(),
					check.That(resourceType+".test_007").Key("policy_name").HasValue("Navigate windows and frames across different domains"),
					check.That(resourceType+".test_007").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_007").Key("values.0.value").HasValue("1"),
				),
			},
			// Final import test
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaGroupPolicyDefinition.ResourceName, "Lifecycle Final Step: Importing Dropdown Configuration")
				},
				ResourceName:      resourceType + ".test_007",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 008: Error Scenarios - placeholder for error validation tests
// TODO: Implement error scenario tests
func TestAccGroupPolicyDefinitionResource_008_ErrorScenarios(t *testing.T) {
	t.Skip("Error scenario tests not yet implemented")
}

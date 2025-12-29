package graphBetaGroupPolicyBooleanValue_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyBooleanValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_boolean_value"
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
// ## Group Policy Boolean Value Testing
//
// Group policy boolean values are managed through the updateDefinitionValues
// endpoint and have dependencies on:
// - Group Policy Configurations (parent resource)
// - Group Policy Definitions (templates in Microsoft Graph)
// - Group Policy Presentations (individual checkboxes/settings)
//
// ## Resource Dependencies
//
// The test creates a group policy configuration first, then creates boolean
// values within that configuration. During cleanup, the boolean value is
// deleted first, followed by the configuration.
//
// ## Timing Considerations
//
// 1. **Read After Write:** After creating or updating a boolean value, there's
//    a brief delay before the value is consistently readable from the API.
//    The ReadWithRetry mechanism handles this automatically.
//
// 2. **CheckDestroy Wait:** A 30-second wait is used before CheckDestroy to
//    ensure the resource deletion has propagated through the backend.
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

const resourceType = graphBetaGroupPolicyBooleanValue.ResourceName

var testResource = graphBetaGroupPolicyBooleanValue.GroupPolicyBooleanValueTestResource{}

// Test 001: Scenario 1 - Minimal configuration
func TestAccGroupPolicyBooleanValueResource_001_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaGroupPolicyBooleanValue.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").Exists(),
					check.That(resourceType+".test_001").Key("group_policy_configuration_id").Exists(),
					check.That(resourceType+".test_001").Key("policy_name").HasValue("Allow Cloud Policy Management"),
					check.That(resourceType+".test_001").Key("class_type").HasValue("machine"),
					check.That(resourceType+".test_001").Key("category_path").HasValue("\\FSLogix\\Profile Containers"),
					check.That(resourceType+".test_001").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_001").Key("values.#").HasValue("1"),
					check.That(resourceType+".test_001").Key("values.0.value").HasValue("true"),
					check.That(resourceType+".test_001").Key("group_policy_definition_value_id").Exists(),
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

// Test 002: Scenario 2 - Maximal configuration
func TestAccGroupPolicyBooleanValueResource_002_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaGroupPolicyBooleanValue.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").Exists(),
					check.That(resourceType+".test_002").Key("policy_name").HasValue("Enable Profile Containers"),
					check.That(resourceType+".test_002").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_002").Key("values.#").HasValue("3"),
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
func TestAccGroupPolicyBooleanValueResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaGroupPolicyBooleanValue.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").Exists(),
					check.That(resourceType+".test_003").Key("enabled").HasValue("false"),
					check.That(resourceType+".test_003").Key("values.#").HasValue("1"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").Exists(),
					check.That(resourceType+".test_003").Key("policy_name").HasValue("Enable Profile Containers"),
					check.That(resourceType+".test_003").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_003").Key("values.#").HasValue("3"),
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
func TestAccGroupPolicyBooleanValueResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaGroupPolicyBooleanValue.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").Exists(),
					check.That(resourceType+".test_004").Key("enabled").HasValue("true"),
					check.That(resourceType+".test_004").Key("values.#").HasValue("3"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").Exists(),
					check.That(resourceType+".test_004").Key("policy_name").HasValue("Allow Cloud Policy Management"),
					check.That(resourceType+".test_004").Key("enabled").HasValue("false"),
					check.That(resourceType+".test_004").Key("values.#").HasValue("1"),
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


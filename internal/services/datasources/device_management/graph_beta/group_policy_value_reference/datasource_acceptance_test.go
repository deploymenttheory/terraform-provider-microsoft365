package graphBetaDeviceManagementGroupPolicyValueReference_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGroupPolicyValueReference "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/group_policy_value_reference"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper function to load acceptance test configs
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// Test 01: Single Machine Policy
func TestAccDatasourceGroupPolicyValueReference_01_SingleMachinePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_single_machine_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Prohibit removal of updates"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("1"),
					// Verify definition structure
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.display_name").HasValue("Prohibit removal of updates"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.class_type").HasValue("machine"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.category_path").Exists(),
				),
			},
		},
	})
}

// Test 02: Multiple Definitions - Policy with multiple variants (user variants)
func TestAccDatasourceGroupPolicyValueReference_02_MultipleDefinitions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_multiple_definitions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Action to take on Microsoft Edge startup"),
					// This policy has multiple variants across different Edge versions
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.#"),
					// Verify first definition structure
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.display_name").HasValue("Action to take on Microsoft Edge startup"),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.class_type"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.category_path").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.policy_type").Exists(),
					// Verify presentations exist
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.presentations.#").Exists(),
				),
			},
		},
	})
}

// Test 03: No Results - Policy that doesn't exist returns empty definitions
func TestAccDatasourceGroupPolicyValueReference_03_NoResults(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_no_results.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("This Policy Does Not Exist At All"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 04: Boolean/CheckBox Presentation Type
func TestAccDatasourceGroupPolicyValueReference_04_BooleanPresentation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_boolean_presentation.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Remove Default Microsoft Store packages from the system."),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("1"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.class_type").HasValue("machine"),
					// Verify presentations with checkbox type
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.presentations.#").Exists(),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.id"),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.label"),
				),
			},
		},
	})
}

// Test 05: Text/TextBox Presentation Type
func TestAccDatasourceGroupPolicyValueReference_05_TextPresentation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_text_presentation.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Browsing Data Lifetime Settings"),
					// This policy may have multiple variants
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.#"),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.class_type"),
					// Verify presentations with textbox type
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.presentations.#").Exists(),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.id"),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.label"),
				),
			},
		},
	})
}

// Test 06: ListBox Presentation Type (Array of Strings)
func TestAccDatasourceGroupPolicyValueReference_06_ListBoxPresentation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("06_listbox_presentation.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Configure list of Enhanced Storage devices usable on your computer"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("1"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.class_type").HasValue("machine"),
					// Verify presentations with listbox type
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.presentations.#").Exists(),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.id"),
					resource.TestCheckResourceAttrSet("data."+utilityGroupPolicyValueReference.DataSourceName+".test", "definitions.0.presentations.0.label"),
				),
			},
		},
	})
}

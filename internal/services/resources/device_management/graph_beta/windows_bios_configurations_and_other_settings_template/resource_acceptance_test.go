package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate_test

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
	graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_bios_configurations_and_other_settings_template"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.ResourceName

var testResource = graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.WindowsBiosConfigurationsAndOtherSettingsTemplateTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
}

// checkDestroyBase verifies the template itself is gone.
func checkDestroyBase() resource.TestCheckFunc {
	return destroy.CheckDestroyedAllFunc(
		testResource,
		graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.ResourceName,
		30*time.Second,
	)
}

// checkDestroyWithGroups also verifies the Entra ID groups created as assignment targets are gone.
func checkDestroyWithGroups() resource.TestCheckFunc {
	return destroy.CheckDestroyedTypesFunc(
		30*time.Second,
		destroy.ResourceTypeMapping{
			ResourceType: graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.ResourceName,
			TestResource: graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.WindowsBiosConfigurationsAndOtherSettingsTemplateTestResource{},
		},
		destroy.ResourceTypeMapping{
			ResourceType: graphBetaGroup.ResourceName,
			TestResource: graphBetaGroup.GroupTestResource{},
		},
	)
}

// waitForAssignmentConsistency absorbs Microsoft Graph's eventual consistency on assignment reads.
func waitForAssignmentConsistency() resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		testlog.WaitForConsistency("windows bios configuration assignments", 20*time.Second)
		time.Sleep(20 * time.Second)
		return nil
	}
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_01_Scenario_Minimal(t *testing.T) {
	name := resourceType + ".test_001"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyBase(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "create minimal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("display_name").Exists(),
					check.That(name).Key("file_name").HasValue("test-bios-001.cctk"),
					check.That(name).Key("configuration_file_content").Exists(),
					check.That(name).Key("hardware_configuration_format").HasValue("dell"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(name).Key("version").Exists(),
					check.That(name).Key("created_date_time").Exists(),
					check.That(name).Key("last_modified_date_time").Exists(),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_02_Scenario_Maximal(t *testing.T) {
	name := resourceType + ".test_002"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyBase(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "create maximal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("description").HasValue("Maximal BIOS configuration template"),
					check.That(name).Key("file_name").HasValue("test-bios-002.cctk"),
					check.That(name).Key("hardware_configuration_format").HasValue("dell"),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	name := resourceType + ".test_003"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyBase(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "create minimal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(name).Key("per_device_password_disabled").HasValue("false"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "promote to maximal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("description").HasValue("Promoted to a maximal BIOS configuration template"),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	name := resourceType + ".test_004"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyBase(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "create maximal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("description").HasValue("Starts as a maximal BIOS configuration template"),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "reduce to minimal bios configuration template") },
				Config:    loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("description").HasValue("Reduced to a minimal BIOS configuration template"),
					check.That(name).Key("per_device_password_disabled").HasValue("false"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_05_AssignmentsMinimal(t *testing.T) {
	name := resourceType + ".test_005"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyWithGroups(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "create bios configuration template with a single assignment")
				},
				Config: loadAcceptanceTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("assignments.#").HasValue("1"),
					check.That(name).Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(name).Key("assignments.0.group_id").Exists(),
					check.That(name).Key("assignments.0.filter_type").HasValue("none"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_06_AssignmentsMaximal(t *testing.T) {
	name := resourceType + ".test_006"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyWithGroups(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "create bios configuration template with filtered and excluded assignments")
				},
				Config: loadAcceptanceTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).ExistsInGraph(testResource),
					check.That(name).Key("assignments.#").HasValue("5"),
					check.That(name).Key("assignments.*.filter_type").ContainsTypeSetElement("include"),
					check.That(name).Key("assignments.*.filter_type").ContainsTypeSetElement("exclude"),
					check.That(name).Key("assignments.*.type").ContainsTypeSetElement("exclusionGroupAssignmentTarget"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	name := resourceType + ".test_007"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyWithGroups(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "create bios configuration template with a single assignment")
				},
				Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "expand to five assignments") },
				Config:    loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).Key("assignments.#").HasValue("5"),
					check.That(name).Key("assignments.*.filter_type").ContainsTypeSetElement("include"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	name := resourceType + ".test_008"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             checkDestroyWithGroups(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "create bios configuration template with five assignments") },
				Config:    loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).Key("assignments.#").HasValue("5"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "reduce to a single assignment") },
				Config:    loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForAssignmentConsistency(),
					check.That(name).Key("assignments.#").HasValue("1"),
					check.That(name).Key("assignments.0.type").HasValue("groupAssignmentTarget"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

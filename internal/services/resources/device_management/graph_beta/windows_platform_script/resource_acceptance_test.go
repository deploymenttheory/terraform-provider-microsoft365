package graphBetaWindowsPlatformScript_test

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
	graphBetaWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaWindowsPlatformScript.ResourceName

var testResource = graphBetaWindowsPlatformScript.WindowsPlatformScriptTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceWindowsPlatformScript_01_Scenario_Minimal(t *testing.T) {
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
			graphBetaWindowsPlatformScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("display_name").Exists(),
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

func TestAccResourceWindowsPlatformScript_02_Scenario_Maximal(t *testing.T) {
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
			graphBetaWindowsPlatformScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("display_name").Exists(),
					check.That(resourceType+".test_002").Key("description").Exists(),
					check.That(resourceType+".test_002").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_002").Key("enforce_signature_check").HasValue("true"),
					check.That(resourceType+".test_002").Key("run_as_32_bit").HasValue("true"),
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

func TestAccResourceWindowsPlatformScript_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	var scriptID string
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
			graphBetaWindowsPlatformScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config:           loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_003", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".test_003", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".test_003").ExistsInGraph(testResource),
					check.That(resourceType+".test_003").Key("file_name").HasValue("test-script-003.ps1"),
					check.That(resourceType+".test_003").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("run_as_account").HasValue("system"),
				),
			},
			{
				ResourceName:      resourceType + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_003", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_003", "id", &scriptID),
					check.That(resourceType+".test_003").ExistsInGraph(testResource),
					check.That(resourceType+".test_003").Key("file_name").HasValue("test-script-003.ps1"),
					check.That(resourceType+".test_003").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("description").Exists(),
					check.That(resourceType+".test_003").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_003").Key("enforce_signature_check").HasValue("true"),
					check.That(resourceType+".test_003").Key("run_as_32_bit").HasValue("true"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_003", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_003", "id", &scriptID),
					check.That(resourceType+".test_003").ExistsInGraph(testResource),
					check.That(resourceType+".test_003").Key("file_name").HasValue("test-script-003.ps1"),
					check.That(resourceType+".test_003").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("run_as_account").HasValue("system"),
				),
			},
			{
				ResourceName:      resourceType + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceWindowsPlatformScript_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	var scriptID string
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
			graphBetaWindowsPlatformScript.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config:           loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_004", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".test_004", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".test_004").ExistsInGraph(testResource),
					check.That(resourceType+".test_004").Key("file_name").HasValue("test-script-004.ps1"),
					check.That(resourceType+".test_004").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_004").Key("enforce_signature_check").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_004", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_004", "id", &scriptID),
					check.That(resourceType+".test_004").ExistsInGraph(testResource),
					check.That(resourceType+".test_004").Key("file_name").HasValue("test-script-004.ps1"),
					check.That(resourceType+".test_004").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("run_as_account").HasValue("system"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_004", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_004", "id", &scriptID),
					check.That(resourceType+".test_004").ExistsInGraph(testResource),
					check.That(resourceType+".test_004").Key("file_name").HasValue("test-script-004.ps1"),
					check.That(resourceType+".test_004").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".test_004").Key("enforce_signature_check").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceWindowsPlatformScript_05_AssignmentsMinimal(t *testing.T) {
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
				ResourceType: graphBetaWindowsPlatformScript.ResourceName,
				TestResource: graphBetaWindowsPlatformScript.WindowsPlatformScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_005").Key("display_name").Exists(),
					check.That(resourceType+".test_005").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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

func TestAccResourceWindowsPlatformScript_06_AssignmentsMaximal(t *testing.T) {
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
				ResourceType: graphBetaWindowsPlatformScript.ResourceName,
				TestResource: graphBetaWindowsPlatformScript.WindowsPlatformScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_006").Key("display_name").Exists(),
					check.That(resourceType+".test_006").Key("assignments.#").HasValue("5"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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

func TestAccResourceWindowsPlatformScript_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	var scriptID string
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
				ResourceType: graphBetaWindowsPlatformScript.ResourceName,
				TestResource: graphBetaWindowsPlatformScript.WindowsPlatformScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config:           loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_007", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".test_007", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".test_007").ExistsInGraph(testResource),
					check.That(resourceType+".test_007").Key("file_name").HasValue("test-script-007.ps1"),
					check.That(resourceType+".test_007").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("display_name").Exists(),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_007", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_007", "id", &scriptID),
					check.That(resourceType+".test_007").ExistsInGraph(testResource),
					check.That(resourceType+".test_007").Key("file_name").HasValue("test-script-007.ps1"),
					check.That(resourceType+".test_007").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("display_name").Exists(),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("5"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_007", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_007", "id", &scriptID),
					check.That(resourceType+".test_007").ExistsInGraph(testResource),
					check.That(resourceType+".test_007").Key("file_name").HasValue("test-script-007.ps1"),
					check.That(resourceType+".test_007").Key("script_content").HasValue("Write-Host 'Hello World'"),
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("display_name").Exists(),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceWindowsPlatformScript_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	var scriptID string
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
				ResourceType: graphBetaWindowsPlatformScript.ResourceName,
				TestResource: graphBetaWindowsPlatformScript.WindowsPlatformScriptTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config:           loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_008", plancheck.ResourceActionCreate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceType+".test_008", "id", func(value string) error { scriptID = value; return nil }),
					check.That(resourceType+".test_008").ExistsInGraph(testResource),
					check.That(resourceType+".test_008").Key("file_name").HasValue("test-script-008.ps1"),
					check.That(resourceType+".test_008").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("5"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_008", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_008", "id", &scriptID),
					check.That(resourceType+".test_008").ExistsInGraph(testResource),
					check.That(resourceType+".test_008").Key("file_name").HasValue("test-script-008.ps1"),
					check.That(resourceType+".test_008").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_3.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_008", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_008", "id", &scriptID),
					check.That(resourceType+".test_008").ExistsInGraph(testResource),
					check.That(resourceType+".test_008").Key("file_name").HasValue("test-script-008.ps1"),
					check.That(resourceType+".test_008").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "assignments.#"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_3.tf"), PlanOnly: true},
			{
				Config:           loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceType+".test_008", plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceType+".test_008", "id", &scriptID),
					check.That(resourceType+".test_008").ExistsInGraph(testResource),
					check.That(resourceType+".test_008").Key("file_name").HasValue("test-script-008.ps1"),
					check.That(resourceType+".test_008").Key("script_content").HasValue("Write-Host 'Maximal script'"),
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("display_name").Exists(),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("5"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows platform script assignments", 20*time.Second)
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
			{Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"), PlanOnly: true},
		},
	})
}

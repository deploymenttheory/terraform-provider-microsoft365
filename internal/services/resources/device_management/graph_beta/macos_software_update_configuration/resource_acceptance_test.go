package graphBetaMacOSSoftwareUpdateConfiguration_test

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
	graphBetaMacOSSoftwareUpdateConfiguration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_software_update_configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaMacOSSoftwareUpdateConfiguration.ResourceName

	// testResource is the test resource implementation for macOS software update configurations
	testResource = graphBetaMacOSSoftwareUpdateConfiguration.MacOSSoftwareUpdateConfigurationTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Deploy Minimal Configuration
func TestAccResourceMacOSSoftwareUpdateConfiguration_01_Minimal(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal macOS software update configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_01_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".test_01_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_01_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-01-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".test_01_minimal").Key("update_schedule_type").HasValue("alwaysUpdate"),
					check.That(resourceType+".test_01_minimal").Key("critical_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("config_data_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("firmware_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("all_other_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".test_01_minimal").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal configuration")
				},
				ResourceName:            resourceType + ".test_01_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 02: Deploy Maximal Configuration
func TestAccResourceMacOSSoftwareUpdateConfiguration_02_Maximal(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal macOS software update configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_02_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".test_02_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_02_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-02-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".test_02_maximal").Key("description").HasValue("Maximal software update configuration for acceptance testing with all features"),
					check.That(resourceType+".test_02_maximal").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_02_maximal").Key("priority").HasValue("high"),
					check.That(resourceType+".test_02_maximal").Key("max_user_deferrals_count").HasValue("5"),
					check.That(resourceType+".test_02_maximal").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".test_02_maximal").Key("custom_update_time_windows.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal configuration")
				},
				ResourceName:            resourceType + ".test_02_maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 03: Minimal to Maximal in Steps
func TestAccResourceMacOSSoftwareUpdateConfiguration_03_MinimalToMaximal(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_03_minimal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_03_progression").ExistsInGraph(testResource),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("alwaysUpdate"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to intermediate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_03_intermediate_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_03_progression").ExistsInGraph(testResource),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_03_progression").Key("priority").HasValue("low"),
					check.That(resourceType+".test_03_progression").Key("custom_update_time_windows.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 3: Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_03_maximal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_03_progression").ExistsInGraph(testResource),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_03_progression").Key("priority").HasValue("high"),
					check.That(resourceType+".test_03_progression").Key("max_user_deferrals_count").HasValue("5"),
					check.That(resourceType+".test_03_progression").Key("custom_update_time_windows.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal configuration")
				},
				ResourceName:            resourceType + ".test_03_progression",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 04: Maximal to Minimal in Steps
func TestAccResourceMacOSSoftwareUpdateConfiguration_04_MaximalToMinimal(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_04_maximal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_04_regression").ExistsInGraph(testResource),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_04_regression").Key("priority").HasValue("high"),
					check.That(resourceType+".test_04_regression").Key("custom_update_time_windows.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to intermediate configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_04_intermediate_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_04_regression").ExistsInGraph(testResource),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_04_regression").Key("custom_update_time_windows.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 3: Updating to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_04_minimal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_04_regression").ExistsInGraph(testResource),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("alwaysUpdate"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal configuration")
				},
				ResourceName:            resourceType + ".test_04_regression",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 05: Minimal Resource with Minimal Assignments
func TestAccResourceMacOSSoftwareUpdateConfiguration_05_MinimalAssignments(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating configuration with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_05_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_05_min_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".test_05_min_assignments").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".test_05_min_assignments").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing configuration with minimal assignments")
				},
				ResourceName:            resourceType + ".test_05_min_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 06: Minimal Resource with Maximal Assignments
func TestAccResourceMacOSSoftwareUpdateConfiguration_06_MaximalAssignments(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating configuration with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_06_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_06_max_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".test_06_max_assignments").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing configuration with maximal assignments")
				},
				ResourceName:            resourceType + ".test_06_max_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 07: Minimal Assignments to Maximal Assignments
func TestAccResourceMacOSSoftwareUpdateConfiguration_07_MinimalToMaximalAssignments(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating configuration with minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_07_minimal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_07_assignments_progression").ExistsInGraph(testResource),
					check.That(resourceType+".test_07_assignments_progression").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_07_maximal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_07_assignments_progression").ExistsInGraph(testResource),
					check.That(resourceType+".test_07_assignments_progression").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing configuration with maximal assignments")
				},
				ResourceName:            resourceType + ".test_07_assignments_progression",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 08: Maximal Assignments to Minimal Assignments
func TestAccResourceMacOSSoftwareUpdateConfiguration_08_MaximalToMinimalAssignments(t *testing.T) {
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
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating configuration with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_08_maximal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_08_assignments_regression").ExistsInGraph(testResource),
					check.That(resourceType+".test_08_assignments_regression").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to minimal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_08_minimal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("macOS software update configuration", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_08_assignments_regression").ExistsInGraph(testResource),
					check.That(resourceType+".test_08_assignments_regression").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing configuration with minimal assignments")
				},
				ResourceName:            resourceType + ".test_08_assignments_regression",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

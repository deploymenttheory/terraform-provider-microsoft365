package graphBetaWindowsAutopilotDevicePreparationPolicy_test

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
	graphBetaWindowsAutopilotDevicePreparationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_device_preparation_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// resourceType is the full Terraform resource type name
	resourceType = graphBetaWindowsAutopilotDevicePreparationPolicy.ResourceName

	// testResource is the test resource implementation for Windows Autopilot Device Preparation Policies
	testResource = graphBetaWindowsAutopilotDevicePreparationPolicy.WindowsAutopilotDevicePreparationPolicyTestResource{}
)

// loadAcceptanceTestTerraform loads a Terraform config from the acceptance test directory.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_01_AutomaticMinimal tests an automatic mode minimal policy.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_01_AutomaticMinimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating automatic mode minimal Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("001_scenario_automatic_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".auto_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".auto_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auto_minimal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-auto-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".auto_minimal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_1"),
					check.That(resourceType+".auto_minimal").Key("allowed_apps.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing automatic mode minimal Windows Autopilot Device Preparation policy")
				},
				ResourceName:            resourceType + ".auto_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_02_AutomaticMaximal tests an automatic mode maximal policy.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_02_AutomaticMaximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating automatic mode maximal Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("002_scenario_automatic_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".auto_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".auto_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auto_maximal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-auto-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".auto_maximal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_1"),
					check.That(resourceType+".auto_maximal").Key("allowed_apps.#").HasValue("1"),
					check.That(resourceType+".auto_maximal").Key("allowed_scripts.#").HasValue("1"),
				),
			},
		},
	})
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_03_UserDrivenMinimal tests a user-driven mode minimal policy.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_03_UserDrivenMinimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating user-driven mode minimal Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("003_scenario_user_driven_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ud_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".ud_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_minimal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-ud-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".ud_minimal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_minimal").Key("device_security_group").IsSet(),
					check.That(resourceType+".ud_minimal").Key("deployment_settings.deployment_mode").HasValue("enrollment_autopilot_dpp_deploymentmode_0"),
					check.That(resourceType+".ud_minimal").Key("oobe_settings.timeout_in_minutes").HasValue("60"),
					check.That(resourceType+".ud_minimal").Key("allowed_apps.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing user-driven mode minimal Windows Autopilot Device Preparation policy")
				},
				ResourceName:            resourceType + ".ud_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_04_UserDrivenMaximal tests a user-driven mode maximal policy.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_04_UserDrivenMaximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating user-driven mode maximal Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("004_scenario_user_driven_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ud_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".ud_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_maximal").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-ud-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".ud_maximal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_maximal").Key("deployment_settings.deployment_mode").HasValue("enrollment_autopilot_dpp_deploymentmode_1"),
					check.That(resourceType+".ud_maximal").Key("oobe_settings.allow_skip").HasValue("true"),
					check.That(resourceType+".ud_maximal").Key("allowed_apps.#").HasValue("1"),
					check.That(resourceType+".ud_maximal").Key("allowed_scripts.#").HasValue("1"),
				),
			},
		},
	})
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_05_UserDrivenMinimalWithMinimalAssignments tests user-driven with minimal assignments.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_05_UserDrivenMinimalWithMinimalAssignments(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating user-driven mode with minimal assignments Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("005_scenario_user_driven_minimal_with_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ud_min_assign").ExistsInGraph(testResource),
					check.That(resourceType+".ud_min_assign").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_min_assign").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-ud-min-assign-[a-z0-9]{8}$`)),
					check.That(resourceType+".ud_min_assign").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_min_assign").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

// TestAccResourceWindowsAutopilotDevicePreparationPolicy_06_UserDrivenMinimalWithMaximalAssignments tests user-driven with maximal assignments.
func TestAccResourceWindowsAutopilotDevicePreparationPolicy_06_UserDrivenMinimalWithMaximalAssignments(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating user-driven mode with maximal assignments Windows Autopilot Device Preparation policy")
				},
				Config: loadAcceptanceTestTerraform("006_scenario_user_driven_minimal_with_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopilot device preparation policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".ud_max_assign").ExistsInGraph(testResource),
					check.That(resourceType+".ud_max_assign").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_max_assign").Key("name").MatchesRegex(regexp.MustCompile(`^acc-test-autopilot-dpp-ud-max-assign-[a-z0-9]{8}$`)),
					check.That(resourceType+".ud_max_assign").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_max_assign").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

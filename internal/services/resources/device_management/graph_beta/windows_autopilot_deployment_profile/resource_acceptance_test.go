package graphBetaWindowsAutopilotDeploymentProfile_test

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
	graphBetaWindowsAutopilotDeploymentProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_deployment_profile"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaWindowsAutopilotDeploymentProfile.ResourceName

	// testResource is the test resource implementation for Windows Autopilot deployment profiles
	testResource = graphBetaWindowsAutopilotDeploymentProfile.WindowsAutopilotDeploymentProfileTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsAutopilotDeploymentProfile_01_SelfDeployingOSDefaultLocale(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating Windows Autopilot deployment profile with OS default locale")
				},
				Config: loadAcceptanceTestTerraform("01_self_deploying_os_default_locale.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows Autopilot deployment profile", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".user_driven").ExistsInGraph(testResource),
					check.That(resourceType+".user_driven").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven").Key("display_name").HasValue("acc test user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("description").HasValue("user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("locale").HasValue("os-default"),
					check.That(resourceType+".user_driven").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven").Key("hardware_hash_extraction_enabled").HasValue("true"),
					check.That(resourceType+".user_driven").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".user_driven").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("false"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.device_usage_type").HasValue("singleUser"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),
					check.That(resourceType+".user_driven").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Windows Autopilot deployment profile")
				},
				ResourceName:            resourceType + ".user_driven",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_02_UserDrivenHybridDomainJoin(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating Windows Autopilot deployment profile with hybrid domain join")
				},
				Config: loadAcceptanceTestTerraform("02_user_driven_hybrid_domain_join.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows Autopilot deployment profile", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("display_name").HasValue("acc_test_user_driven_japanese_preprovisioned"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("description").HasValue("user driven autopilot profile with japanese locale and allow pre provisioned deployment"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("locale").HasValue("ja-JP"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("device_join_type").HasValue("microsoft_entra_hybrid_joined"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Windows Autopilot deployment profile")
				},
				ResourceName:            resourceType + ".user_driven_japanese_preprovisioned_with_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_02_UserDrivenWithGroupAssignments(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating Windows Autopilot deployment profile with group assignments")
				},
				Config: loadAcceptanceTestTerraform("03_user_driven_with_group_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows Autopilot deployment profile", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".user_driven").ExistsInGraph(testResource),
					check.That(resourceType+".user_driven").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven").Key("display_name").HasValue("acc test user driven autopilot with group assignments"),
					check.That(resourceType+".user_driven").Key("description").HasValue("user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("locale").HasValue("os-default"),
					check.That(resourceType+".user_driven").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".user_driven").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Windows Autopilot deployment profile")
				},
				ResourceName:            resourceType + ".user_driven",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_04_HololensWithAllDeviceAssignment(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating Windows Autopilot deployment profile for HoloLens")
				},
				Config: loadAcceptanceTestTerraform("04_hololens_with_all_device_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows Autopilot deployment profile", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".hololens_with_all_device_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("display_name").HasValue("acc_test_hololens_with_all_device_assignment"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("description").HasValue("hololens autopilot profile with hk locale and group assignment"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("locale").HasValue("zh-HK"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("preprovisioning_allowed").HasValue("false"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("device_type").HasValue("holoLens"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("hardware_hash_extraction_enabled").HasValue("false"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.device_usage_type").HasValue("shared"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Windows Autopilot deployment profile")
				},
				ResourceName:            resourceType + ".hololens_with_all_device_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

package graphBetaWindowsEnrollmentStatusPage_test

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
	graphBetaWindowsEnrollmentStatusPage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_enrollment_status_page"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testResource = graphBetaWindowsEnrollmentStatusPage.WindowsEnrollmentStatusPageTestResource{}

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// Test 001: Minimal Configuration
func TestAccWindowsEnrollmentStatusPageResource_001_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Creating Minimal Configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("show_installation_progress").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("install_progress_timeout_in_minutes").HasValue("120"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("custom_error_message").HasValue("Contact IT support for assistance"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Importing Minimal Configuration")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 002: Maximal Configuration
func TestAccWindowsEnrollmentStatusPageResource_002_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Creating Maximal Configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("description").HasValue("Test description for maximal enrollment status page"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("show_installation_progress").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("custom_error_message").HasValue("Contact IT support for assistance"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_quality_updates").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_progress_timeout_in_minutes").HasValue("120"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_log_collection_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_use_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_fail_selected_blocking_apps_in_technician_phase").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Importing Maximal Configuration")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 003: Configuration with Assignments
func TestAccWindowsEnrollmentStatusPageResource_003_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Creating Configuration with Assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_with_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".with_assignments").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".with_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".with_assignments").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".with_assignments").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Importing Configuration with Assignments")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".with_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 004: Update Lifecycle Test
func TestAccWindowsEnrollmentStatusPageResource_004_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Creating Initial Configuration (Step 1 - Minimal)")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("install_progress_timeout_in_minutes").HasValue("120"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Updating to Maximal Configuration (Step 2)")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_progress_timeout_in_minutes").HasValue("120"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Importing Updated Configuration")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 005: Lifecycle Minimal to Maximal
func TestAccWindowsEnrollmentStatusPageResource_005_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Lifecycle: Creating minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
					time.Sleep(15 * time.Second)
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Lifecycle: Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Lifecycle: Importing configuration")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".lifecycle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 006: Lifecycle Maximal to Minimal
func TestAccWindowsEnrollmentStatusPageResource_006_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsEnrollmentStatusPage.ResourceName,
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
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Downgrade: Creating maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
					time.Sleep(15 * time.Second)
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Downgrade: Updating to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("Windows enrollment status page", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").ExistsInGraph(testResource),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_Windows10EnrollmentCompletionPageConfiguration$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").IsNotEmpty(),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(graphBetaWindowsEnrollmentStatusPage.ResourceName, "Downgrade: Importing configuration")
				},
				ResourceName:            graphBetaWindowsEnrollmentStatusPage.ResourceName + ".lifecycle",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

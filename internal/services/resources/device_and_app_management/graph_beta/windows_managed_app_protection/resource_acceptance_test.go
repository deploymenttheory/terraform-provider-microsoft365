package graphBetaDeviceAndAppManagementWindowsManagedAppProtection_test

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
	graphBetaWindowsManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/windows_managed_app_protection"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaWindowsManagedAppProtection.ResourceName

var testResource = graphBetaWindowsManagedAppProtection.WindowsManagedAppProtectionTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceWindowsManagedAppProtection_01_Scenario_Minimal tests the creation
// of a Windows MAM policy with only the required fields set.
func TestAccResourceWindowsManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("display_name").Exists(),
					check.That(resourceType+".test_001").Key("print_blocked").HasValue("false"),
					check.That(resourceType+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(resourceType+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
					check.That(resourceType+".test_001").Key("maximum_allowed_device_threat_level").HasValue("notConfigured"),
					check.That(resourceType+".test_001").Key("mobile_threat_defense_remediation_action").HasValue("block"),
					check.That(resourceType+".test_001").Key("period_offline_before_wipe_is_enforced").HasValue("P90D"),
					check.That(resourceType+".test_001").Key("period_offline_before_access_check").HasValue("P30D"),
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

// TestAccResourceWindowsManagedAppProtection_02_Scenario_Maximal tests the creation
// of a Windows MAM policy with all available fields set.
func TestAccResourceWindowsManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("display_name").Exists(),
					check.That(resourceType+".test_002").Key("description").Exists(),
					check.That(resourceType+".test_002").Key("print_blocked").HasValue("true"),
					check.That(resourceType+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(resourceType+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(resourceType+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("none"),
					check.That(resourceType+".test_002").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
					check.That(resourceType+".test_002").Key("maximum_allowed_device_threat_level").HasValue("low"),
					check.That(resourceType+".test_002").Key("mobile_threat_defense_remediation_action").HasValue("wipe"),
					check.That(resourceType+".test_002").Key("minimum_required_os_version").HasValue("10.0.19041"),
					check.That(resourceType+".test_002").Key("minimum_warning_os_version").HasValue("10.0.18363"),
					check.That(resourceType+".test_002").Key("minimum_wipe_os_version").HasValue("10.0.17763"),
					check.That(resourceType+".test_002").Key("minimum_required_app_version").HasValue("1.0.0"),
					check.That(resourceType+".test_002").Key("minimum_warning_app_version").HasValue("1.1.0"),
					check.That(resourceType+".test_002").Key("minimum_wipe_app_version").HasValue("0.9.0"),
					check.That(resourceType+".test_002").Key("period_offline_before_wipe_is_enforced").HasValue("P30D"),
					check.That(resourceType+".test_002").Key("period_offline_before_access_check").HasValue("P7D"),
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

// TestAccResourceWindowsManagedAppProtection_03_Lifecycle_MinimalToMaximal tests
// updating a policy from a minimal config to a maximal config.
func TestAccResourceWindowsManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("print_blocked").HasValue("false"),
					check.That(resourceType+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("display_name").Exists(),
					check.That(resourceType+".test_003").Key("description").Exists(),
					check.That(resourceType+".test_003").Key("print_blocked").HasValue("true"),
					check.That(resourceType+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(resourceType+".test_003").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(resourceType+".test_003").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
					check.That(resourceType+".test_003").Key("minimum_required_os_version").HasValue("10.0.19041"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows managed app protection", 20*time.Second)
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
		},
	})
}

// TestAccResourceWindowsManagedAppProtection_04_Lifecycle_MaximalToMinimal tests
// updating a policy from a maximal config back down to a minimal config.
func TestAccResourceWindowsManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("print_blocked").HasValue("true"),
					check.That(resourceType+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(resourceType+".test_004").Key("app_action_if_unable_to_authenticate_user").HasValue("block"),
					check.That(resourceType+".test_004").Key("minimum_required_os_version").HasValue("10.0.19041"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("display_name").Exists(),
					check.That(resourceType+".test_004").Key("print_blocked").HasValue("false"),
					check.That(resourceType+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows managed app protection", 20*time.Second)
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
		},
	})
}

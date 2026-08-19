package graphBetaDeviceAndAppManagementAndroidManagedAppProtection_test

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
	graphBetaAndroidManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/android_managed_app_protection"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestResource = graphBetaAndroidManagedAppProtection.AndroidManagedAppProtectionTestResource{}

func loadAcceptanceTestTerraform(t *testing.T, filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		t.Skipf("skipping acceptance test: fixture file not found: %s", filename)
		return ""
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceAndroidManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaAndroidManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("pin_required").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("encrypt_app_data").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_wipe_is_enforced").HasValue("P90D"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_access_check").HasValue("P30D"),
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAndroidManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaAndroidManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("description").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("blocked"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("screen_capture_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_pin_length").HasValue("6"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("maximum_pin_retries").HasValue("10"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("pin_character_set").HasValue("alphanumericAndSymbol"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_required_os_version").HasValue("9.0"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("minimum_required_app_version").HasValue("2.0.0"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_002").Key("allowed_data_storage_locations.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAndroidManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaAndroidManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("false"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform(t, "003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("description").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_003").Key("minimum_pin_length").HasValue("6"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("android managed app protection", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAndroidManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaAndroidManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("minimum_pin_length").HasValue("6"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform(t, "004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("display_name").Exists(),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("false"),
					check.That(graphBetaAndroidManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("android managed app protection", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      graphBetaAndroidManagedAppProtection.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

package graphBetaDeviceAndAppManagementIosManagedAppProtection_test

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
	graphBetaIosManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_app_protection"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestResource = graphBetaIosManagedAppProtection.IosManagedAppProtectionTestResource{}

func loadAcceptanceTestTerraform(t *testing.T, filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		t.Skipf("skipping acceptance test: fixture file not found: %s", filename)
		return ""
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceIosManagedAppProtection_01_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
    		"random": {
        		Source: "hashicorp/random",
    		},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaIosManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("pin_required").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("face_id_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("allowed_outbound_data_transfer_destinations").HasValue("allApps"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_wipe_is_enforced").HasValue("P90D"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("period_offline_before_access_check").HasValue("P30D"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_001").Key("app_data_encryption_type").HasValue("whenDeviceLocked"),
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceIosManagedAppProtection_02_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
    		"random": {
        		Source: "hashicorp/random",
    		},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaIosManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("description").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_data_transfer_destinations").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_outbound_clipboard_sharing_level").HasValue("blocked"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("face_id_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("third_party_keyboards_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("filter_open_in_to_only_managed_apps").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("app_data_encryption_type").HasValue("afterDeviceRestart"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_pin_length").HasValue("6"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("maximum_pin_retries").HasValue("10"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("pin_character_set").HasValue("alphanumericAndSymbol"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_required_os_version").HasValue("15.0"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("minimum_required_app_version").HasValue("2.0.0"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_002").Key("allowed_data_storage_locations.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceIosManagedAppProtection_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
    		"random": {
        		Source: "hashicorp/random",
    		},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaIosManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform(t, "003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("description").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_003").Key("minimum_pin_length").HasValue("6"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios managed app protection", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceIosManagedAppProtection_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
    		"random": {
        		Source: "hashicorp/random",
    		},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			graphBetaIosManagedAppProtection.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform(t, "004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("none"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("data_backup_blocked").HasValue("true"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("minimum_pin_length").HasValue("6"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform(t, "004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^T_[0-9a-fA-F-]+$`)),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("display_name").Exists(),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("print_blocked").HasValue("false"),
					check.That(graphBetaIosManagedAppProtection.ResourceName+".test_004").Key("allowed_inbound_data_transfer_sources").HasValue("allApps"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios managed app protection", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:      graphBetaIosManagedAppProtection.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
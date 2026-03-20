package graphBetaWindowsUpdatesAutopatchDeviceRegistration_test

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
	graphBetaWindowsAutopatchDeviceRegistration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_device_registration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaWindowsAutopatchDeviceRegistration.ResourceName

var testResource = graphBetaWindowsAutopatchDeviceRegistration.WindowsAutopatchDeviceRegistrationTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceWindowsAutopatchDeviceRegistration_01_Scenario_Minimal(t *testing.T) {
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
			graphBetaWindowsAutopatchDeviceRegistration.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").HasValue("feature"),
					check.That(resourceType+".test_001").Key("update_category").HasValue("feature"),
					check.That(resourceType+".test_001").Key("entra_device_object_ids.#").Exists(),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entra_device_object_ids"},
			},
		},
	})
}

func TestAccResourceWindowsAutopatchDeviceRegistration_02_Scenario_Maximal(t *testing.T) {
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
			graphBetaWindowsAutopatchDeviceRegistration.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").HasValue("quality"),
					check.That(resourceType+".test_002").Key("update_category").HasValue("quality"),
					check.That(resourceType+".test_002").Key("entra_device_object_ids.#").MatchesRegex(regexp.MustCompile(`^[1-9]\d*$`)),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entra_device_object_ids"},
			},
		},
	})
}

func TestAccResourceWindowsAutopatchDeviceRegistration_03_Lifecycle_AddDevices(t *testing.T) {
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
			graphBetaWindowsAutopatchDeviceRegistration.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_add_devices_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").HasValue("feature"),
					check.That(resourceType+".test_003").Key("entra_device_object_ids.#").Exists(),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_add_devices_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").HasValue("feature"),
					check.That(resourceType+".test_003").Key("entra_device_object_ids.#").MatchesRegex(regexp.MustCompile(`^[2-9]\d*$`)),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entra_device_object_ids"},
			},
		},
	})
}

func TestAccResourceWindowsAutopatchDeviceRegistration_04_Lifecycle_RemoveDevices(t *testing.T) {
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
			graphBetaWindowsAutopatchDeviceRegistration.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_remove_devices_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").HasValue("feature"),
					check.That(resourceType+".test_004").Key("entra_device_object_ids.#").MatchesRegex(regexp.MustCompile(`^[2-9]\d*$`)),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_remove_devices_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").HasValue("feature"),
					check.That(resourceType+".test_004").Key("entra_device_object_ids.#").Exists(),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch device registration", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_004",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entra_device_object_ids"},
			},
		},
	})
}

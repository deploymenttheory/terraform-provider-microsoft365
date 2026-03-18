package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_test

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
	graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_updatable_asset_group_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const accResourceType = graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName

var accTestResource = graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentTestResource{}

func externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: constants.ExternalProviderTimeVersion,
		},
	}
}

func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_01_Minimal creates an
// updatable asset group, assigns one device to it, then verifies import.
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			accResourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Create group assignment with one device")
				},
				Config: loadAccTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(accResourceType+".test_001").Key("id").Exists(),
					check.That(accResourceType+".test_001").Key("updatable_asset_group_id").Exists(),
					check.That(accResourceType+".test_001").Key("entra_device_ids.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(accResourceType, 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            accResourceType + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_02_LifecycleAddMembers
// starts with one device member and updates to two, verifying the diff-based add.
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_02_LifecycleAddMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			accResourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Create group assignment with one device")
				},
				Config: loadAccTestTerraform("002_lifecycle_add_members_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(accResourceType+".test_002").Key("entra_device_ids.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(accResourceType, 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Add second device to group assignment")
				},
				Config: loadAccTestTerraform("002_lifecycle_add_members_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(accResourceType+".test_002").Key("entra_device_ids.#").MatchesRegex(regexp.MustCompile(`^[1-9]\d*$`)),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(accResourceType, 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            accResourceType + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_03_LifecycleRemoveMembers
// starts with two device members and updates to one, verifying the diff-based remove.
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_03_LifecycleRemoveMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			accTestResource,
			accResourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Create group assignment with two devices")
				},
				Config: loadAccTestTerraform("003_lifecycle_remove_members_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(accResourceType+".test_003").Key("entra_device_ids.#").MatchesRegex(regexp.MustCompile(`^[1-9]\d*$`)),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(accResourceType, 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Remove one device from group assignment")
				},
				Config: loadAccTestTerraform("003_lifecycle_remove_members_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(accResourceType+".test_003").Key("entra_device_ids.#").HasValue("1"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency(accResourceType, 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            accResourceType + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

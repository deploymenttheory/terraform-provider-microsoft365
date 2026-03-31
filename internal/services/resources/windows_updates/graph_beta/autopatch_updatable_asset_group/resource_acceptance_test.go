package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_updatable_asset_group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testResource = graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup.WindowsUpdatesAutopatchUpdatableAssetGroupTestResource{}

// loadAcceptanceTestTerraform loads an acceptance test HCL config from the acceptance test directory.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_01_CreateEmpty tests that
// an empty updatable asset group (no members) can be created, imported, and destroyed.
//
// API calls exercised:
//   - POST /admin/windows/updates/updatableAssets
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//   - DELETE /admin/windows/updates/updatableAssets/{id}
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_01_CreateEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating empty updatable asset group")
				},
				Config: loadAcceptanceTestTerraform("01_create_empty.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "entra_device_object_ids"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_02_WithMembers tests that
// an updatable asset group can be created with device members via entra_device_object_ids.
//
// API calls exercised:
//   - POST /admin/windows/updates/updatableAssets
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//   - DELETE /admin/windows/updates/updatableAssets/{id}
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_02_WithMembers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating updatable asset group with 1 device member")
				},
				Config: loadAcceptanceTestTerraform("02_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("asset group members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "entra_device_object_ids"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_03_LifecycleAddMember tests the
// diff-based update path for adding a member: starting from an empty group and adding one
// device via addMembersById through the Update function.
//
// API calls exercised (step 1):
//   - POST /admin/windows/updates/updatableAssets
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// API calls exercised (step 2):
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_03_LifecycleAddMember(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating empty asset group with no members")
				},
				Config: loadAcceptanceTestTerraform("01_create_empty.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Adding one device member via diff-based addMembersById")
				},
				Config: loadAcceptanceTestTerraform("03_lifecycle_add_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("asset group members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 3: Import state verification")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "entra_device_object_ids"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_04_LifecycleRemoveMember tests the
// diff-based update path for removing a member: starting with 1 device and removing it via
// removeMembersById through the Update function, leaving the group empty.
//
// API calls exercised (step 1):
//   - POST /admin/windows/updates/updatableAssets
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// API calls exercised (step 2):
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.removeMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
func TestAccResourceWindowsUpdatesAutopatchUpdatableAssetGroup_04_LifecycleRemoveMember(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating asset group with 1 device member")
				},
				Config: loadAcceptanceTestTerraform("03_lifecycle_add_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("asset group members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Removing device member via diff-based removeMembersById")
				},
				Config: loadAcceptanceTestTerraform("01_create_empty.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("asset group members", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 3: Import state verification")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "entra_device_object_ids"},
			},
		},
	})
}

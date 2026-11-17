package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assetGroupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_updatable_asset_group/mocks"
)

const resourceType = "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group"

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *assetGroupMocks.WindowsUpdateUpdatableAssetGroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assetGroupMock := &assetGroupMocks.WindowsUpdateUpdatableAssetGroupMock{}
	assetGroupMock.RegisterMocks()
	return mockClient, assetGroupMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *assetGroupMocks.WindowsUpdateUpdatableAssetGroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assetGroupMock := &assetGroupMocks.WindowsUpdateUpdatableAssetGroupMock{}
	assetGroupMock.RegisterErrorMocks()
	return mockClient, assetGroupMock
}

// TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_01_CreateEmpty tests that
// an empty updatable asset group (no members) can be created and read back.
//
// API calls exercised:
//   - POST /admin/windows/updates/updatableAssets
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//   - DELETE /admin/windows/updates/updatableAssets/{id}
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_01_CreateEmpty(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create_empty.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("0"),
				),
			},
		},
	})
}

// TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_02_CreateWithMembers tests that
// an updatable asset group can be created with an initial device member.
//
// API calls exercised:
//   - POST /admin/windows/updates/updatableAssets
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//   - DELETE /admin/windows/updates/updatableAssets/{id}
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_02_CreateWithMembers(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_create_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"),
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_03_UpdateAddMember tests the
// diff-based update path: starting with 1 member and adding a second via addMembersById.
//
// API calls exercised (step 1):
//   - POST /admin/windows/updates/updatableAssets
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// API calls exercised (step 2):
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_03_UpdateAddMember(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_create_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("entra_device_object_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("03_update_add_member.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("entra_device_object_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_04_Import tests that an
// updatable asset group can be imported by ID.
//
// API calls exercised:
//   - GET /admin/windows/updates/updatableAssets/{id}
//   - GET /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_04_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create_empty.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "entra_device_object_ids"},
			},
		},
	})
}

// TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_05_Error tests that API errors
// (e.g. 403 Forbidden) surface correctly as Terraform diagnostics.
//
// API calls exercised:
//   - POST /admin/windows/updates/updatableAssets (returns 403)
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_05_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_create_empty.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}

// The acceptance fixtures must exercise membership even when Intune has no devices.
func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_06_NoManagedDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices`,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{"value": []any{}}))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_with_members.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("entra_device_object_ids.#").HasValue("2"),
					check.That(resourceType+".test").Key("entra_device_object_ids.*").ContainsTypeSetElement("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					check.That(resourceType+".test").Key("entra_device_object_ids.*").ContainsTypeSetElement("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config: loadAcceptanceTestTerraform("01_create_empty.tf"),
				Check:  check.That(resourceType + ".test").Key("entra_device_object_ids.#").HasValue("0"),
			},
			{
				Config: loadAcceptanceTestTerraform("03_lifecycle_add_step1.tf"),
				Check:  check.That(resourceType + ".test").Key("entra_device_object_ids.#").HasValue("2"),
			},
			{
				Config: loadAcceptanceTestTerraform("04_lifecycle_remove_step2.tf"),
				Check:  check.That(resourceType + ".test").Key("entra_device_object_ids.#").HasValue("0"),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_07_NoWindowsDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices\?`,
		httpmock.NewJsonResponderOrPanic(200, map[string]any{"value": []any{}}))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadAcceptanceTestTerraform("02_with_members.tf"),
				ExpectError: regexp.MustCompile("This acceptance test requires at least one enabled Windows Entra"),
			},
		},
	})
}

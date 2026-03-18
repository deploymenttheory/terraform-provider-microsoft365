package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_updatable_asset_group_assignment"
	assignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_updatable_asset_group_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *assignmentMocks.WindowsUpdateUpdatableAssetGroupAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	assignmentMock := &assignmentMocks.WindowsUpdateUpdatableAssetGroupAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *assignmentMocks.WindowsUpdateUpdatableAssetGroupAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	assignmentMock := &assignmentMocks.WindowsUpdateUpdatableAssetGroupAssignmentMock{}
	assignmentMock.RegisterErrorMocks()
	return mockClient, assignmentMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_01_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName+".test").Key("id").HasValue("d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"),
					check.That(graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName+".test").Key("updatable_asset_group_id").HasValue("d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"),
					check.That(graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName+".test").Key("entra_device_ids.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_02_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName+".test").Key("entra_device_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("02_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName+".test").Key("entra_device_ids.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_03_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
			},
			{
				ResourceName:            graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment.ResourceName + ".test",
				ImportState:             true,
				ImportStateId:           "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroupAssignment_04_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_create.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}

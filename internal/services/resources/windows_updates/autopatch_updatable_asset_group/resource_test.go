package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assetGroupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_updatable_asset_group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
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

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_01_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_02_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdatesAutopatchUpdatableAssetGroup_03_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assetGroupMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assetGroupMock.CleanupMockState()

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

package graphBetaWindowsUpdatesAutopatchRing_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	ringMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_ring/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *ringMocks.WindowsUpdateRingMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	ringMock := &ringMocks.WindowsUpdateRingMock{}
	ringMock.RegisterMocks()
	return mockClient, ringMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *ringMocks.WindowsUpdateRingMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	ringMock := &ringMocks.WindowsUpdateRingMock{}
	ringMock.RegisterErrorMocks()
	return mockClient, ringMock
}

func TestUnitResourceWindowsUpdateRing_01_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ringMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ringMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("b2c3d4e5-2345-6789-bcde-b2c3d4e5f6a7"),
					check.That(resourceType+".test").Key("policy_id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
					check.That(resourceType+".test").Key("display_name").HasValue("Test Ring"),
					check.That(resourceType+".test").Key("description").HasValue("A test ring for unit tests"),
					check.That(resourceType+".test").Key("is_paused").HasValue("false"),
					check.That(resourceType+".test").Key("deferral_in_days").HasValue("7"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("last_modified_date_time").Exists(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateRing_02_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ringMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ringMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_create.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_paused").HasValue("false"),
					check.That(resourceType+".test").Key("deferral_in_days").HasValue("7"),
				),
			},
			{
				Config: loadUnitTestTerraform("02_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_paused").HasValue("true"),
					check.That(resourceType+".test").Key("deferral_in_days").HasValue("14"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateRing_03_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ringMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ringMock.CleanupMockState()

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
				ImportStateId:           "983f03cd-03cd-983f-cd03-3f98cd033f98/b2c3d4e5-2345-6789-bcde-b2c3d4e5f6a7",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdateRing_04_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ringMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ringMock.CleanupMockState()

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

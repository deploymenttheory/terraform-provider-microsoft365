package graphBetaDevice_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	deviceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/device/mocks"
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

func setupMockEnvironment() (*mocks.Mocks, *deviceMocks.DeviceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	dMock := &deviceMocks.DeviceMock{}
	dMock.RegisterMocks()
	return mockClient, dMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *deviceMocks.DeviceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	dMock := &deviceMocks.DeviceMock{}
	dMock.RegisterErrorMocks()
	return mockClient, dMock
}

func TestUnitDatasourceDevice_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-000481110457"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.1.display_name").HasValue("DT-TEST-DEVICE-002"),
					check.That(dataSourceType+".test").Key("items.2.display_name").HasValue("DT-TEST-DEVICE-003"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_02_ByObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-000481110457"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.0.is_compliant").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").HasValue("DT-000481110457"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-000481110457"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_04_ByDeviceId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_id").HasValue("06771871-1375-494e-97f9-ab87ba64edeb"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.device_id").HasValue("06771871-1375-494e-97f9-ab87ba64edeb"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-000481110457"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_05_ByODataQuery(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows' and isCompliant eq true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.0.is_compliant").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.1.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.1.is_compliant").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_06_WithMemberOf(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_with_member_of.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("list_member_of").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "member_of.#"),
					check.That(dataSourceType+".test").Key("member_of.0.id").HasValue("3df4b46e-776a-4c46-9aef-7350661f6529"),
					check.That(dataSourceType+".test").Key("member_of.0.odata_type").HasValue("#microsoft.graph.group"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_07_WithRegisteredOwners(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_with_registered_owners.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("list_registered_owners").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_owners.#"),
					check.That(dataSourceType+".test").Key("registered_owners.0.id").HasValue("7e9c394c-aa2d-4444-a782-a57d200d6c74"),
					check.That(dataSourceType+".test").Key("registered_owners.0.odata_type").HasValue("#microsoft.graph.user"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_08_WithRegisteredUsers(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_with_registered_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("list_registered_users").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_users.#"),
					check.That(dataSourceType+".test").Key("registered_users.0.id").HasValue("7e9c394c-aa2d-4444-a782-a57d200d6c74"),
					check.That(dataSourceType+".test").Key("registered_users.0.odata_type").HasValue("#microsoft.graph.user"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_09_Comprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("list_member_of").HasValue("true"),
					check.That(dataSourceType+".test").Key("list_registered_owners").HasValue("true"),
					check.That(dataSourceType+".test").Key("list_registered_users").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "member_of.#"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_owners.#"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_users.#"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("23ace577-ee29-416f-8566-11c948310bff"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-000481110457"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
				),
			},
		},
	})
}

func TestUnitDatasourceDevice_10_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, dMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer dMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("02_by_object_id.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|insufficient|privileges"),
			},
		},
	})
}

package graphBetaManagedDevice_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaManagedDevice "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/managed_device"
	managedDeviceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/managed_device/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaManagedDevice.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *managedDeviceMocks.ManagedDeviceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mdMock := &managedDeviceMocks.ManagedDeviceMock{}
	mdMock.RegisterMocks()
	return mockClient, mdMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *managedDeviceMocks.ManagedDeviceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mdMock := &managedDeviceMocks.ManagedDeviceMock{}
	mdMock.RegisterErrorMocks()
	return mockClient, mdMock
}

func TestUnitDatasourceManagedDevice_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.1.device_name").HasValue("DESKTOP-WIN-002"),
					check.That(dataSourceType+".test").Key("items.2.device_name").HasValue("LAPTOP-WIN-003"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_02_ByDeviceId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.0.compliance_state").HasValue("compliant"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_03_ByDeviceName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_device_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_04_ByOperatingSystem(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_operating_system.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.1.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.2.operating_system").HasValue("Windows"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_05_ByOsAndVersion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_by_os_and_version.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("os_version").HasValue("10.0.19045"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.0.os_version").HasValue("10.0.19045"),
					check.That(dataSourceType+".test").Key("items.0.compliance_state").HasValue("compliant"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_06_ByAzureADDeviceId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_by_azure_ad_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("azure_ad_device_id").HasValue("aaaaaaaa-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.azure_ad_device_id").HasValue("aaaaaaaa-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_07_BySerialNumber(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_by_serial_number.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("serial_number").HasValue("SN-WIN-001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.serial_number").HasValue("SN-WIN-001"),
					check.That(dataSourceType+".test").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_08_ByUserPrincipalName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_by_user_principal_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("user_principal_name").HasValue("user1@contoso.com"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.user_principal_name").HasValue("user1@contoso.com"),
					check.That(dataSourceType+".test").Key("items.1.user_principal_name").HasValue("user2@contoso.com"),
					check.That(dataSourceType+".test").Key("items.2.user_principal_name").HasValue("user3@contoso.com"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_09_ByODataQuery(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_by_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows' and complianceState eq 'compliant'"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.0.compliance_state").HasValue("compliant"),
					check.That(dataSourceType+".test").Key("items.1.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".test").Key("items.1.compliance_state").HasValue("compliant"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_10_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("02_by_device_id.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|insufficient|authorisation"),
			},
		},
	})
}

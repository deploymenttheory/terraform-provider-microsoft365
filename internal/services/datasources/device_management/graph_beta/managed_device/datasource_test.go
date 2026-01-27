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

// Helper function to load unit test Terraform configurations
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

func TestUnitDatasourceManagedDevice_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".all").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".all").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".all").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".all").Key("items.1.device_name").HasValue("DESKTOP-WIN-002"),
					check.That(dataSourceType+".all").Key("items.2.device_name").HasValue("LAPTOP-WIN-003"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".by_id").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".by_id").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".by_id").Key("items.0.operating_system").HasValue("Windows"),
					check.That(dataSourceType+".by_id").Key("items.0.compliance_state").HasValue("compliant"),
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
					check.That(dataSourceType+".by_device_name").Key("filter_type").HasValue("device_name"),
					check.That(dataSourceType+".by_device_name").Key("filter_value").HasValue("DESKTOP-WIN-001"),
					check.That(dataSourceType+".by_device_name").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".by_device_name").Key("items.0.device_name").HasValue("DESKTOP-WIN-001"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_04_BySerialNumber(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_serial_number.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_serial_number").Key("filter_type").HasValue("serial_number"),
					check.That(dataSourceType+".by_serial_number").Key("filter_value").HasValue("SN-WIN-001"),
					check.That(dataSourceType+".by_serial_number").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".by_serial_number").Key("items.0.serial_number").HasValue("SN-WIN-001"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_05_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("complianceState eq 'compliant'"),
					check.That(dataSourceType+".odata_filter").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".odata_filter").Key("items.0.compliance_state").HasValue("compliant"),
					check.That(dataSourceType+".odata_filter").Key("items.1.compliance_state").HasValue("compliant"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_06_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_advanced").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_advanced").Key("odata_filter").HasValue("operatingSystem eq 'Windows'"),
					check.That(dataSourceType+".odata_advanced").Key("odata_orderby").HasValue("deviceName"),
					check.That(dataSourceType+".odata_advanced").Key("odata_select").HasValue("id,deviceName,operatingSystem,complianceState"),
					check.That(dataSourceType+".odata_advanced").Key("items.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_07_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_comprehensive").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_filter").HasValue("operatingSystem eq 'Windows'"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_top").HasValue("50"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_orderby").HasValue("lastSyncDateTime desc"),
					check.That(dataSourceType+".odata_comprehensive").Key("items.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitDatasourceManagedDevice_08_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_all.tf"),
				ExpectError: regexp.MustCompile("Forbidden - 403"),
			},
		},
	})
}

package graphBetaManagedDevice_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	managedDeviceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/managed_device/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

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

func TestManagedDeviceDataSource_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "filter_type", "all"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.#", "3"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.device_name", "DESKTOP-WIN-001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.operating_system", "Windows"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.1.device_name", "DESKTOP-WIN-002"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "items.2.device_name", "LAPTOP-WIN-003"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "filter_value", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "items.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "items.0.device_name", "DESKTOP-WIN-001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "items.0.operating_system", "Windows"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_id", "items.0.compliance_state", "compliant"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ByDeviceName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDeviceName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "filter_type", "device_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "filter_value", "DESKTOP-WIN-001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "items.0.device_name", "DESKTOP-WIN-001"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_BySerialNumber(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBySerialNumber(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_serial_number", "filter_type", "serial_number"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_serial_number", "filter_value", "SN-WIN-001"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_serial_number", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_serial_number", "items.0.serial_number", "SN-WIN-001"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "odata_filter", "complianceState eq 'compliant'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "items.#", "2"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "items.0.compliance_state", "compliant"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "items.1.compliance_state", "compliant"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_filter", "operatingSystem eq 'Windows'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_orderby", "deviceName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_select", "id,deviceName,operatingSystem,complianceState"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "items.#", "2"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_filter", "operatingSystem eq 'Windows'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_top", "50"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_orderby", "lastSyncDateTime desc"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "items.#", "2"),
				),
			},
		},
	})
}

func TestManagedDeviceDataSource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mdMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigAll(),
				ExpectError: regexp.MustCompile("Forbidden - 403"),
			},
		},
	})
}

// Configuration functions
func testConfigAll() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_all.tf")
	if err != nil {
		panic("failed to load 01_all config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigById() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_by_id.tf")
	if err != nil {
		panic("failed to load 02_by_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByDeviceName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_device_name.tf")
	if err != nil {
		panic("failed to load 03_by_device_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBySerialNumber() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_by_serial_number.tf")
	if err != nil {
		panic("failed to load 04_by_serial_number config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataFilter() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_odata_filter.tf")
	if err != nil {
		panic("failed to load 05_odata_filter config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataAdvanced() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_odata_advanced.tf")
	if err != nil {
		panic("failed to load 06_odata_advanced config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataComprehensive() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/07_odata_comprehensive.tf")
	if err != nil {
		panic("failed to load 07_odata_comprehensive config: " + err.Error())
	}
	return unitTestConfig
}

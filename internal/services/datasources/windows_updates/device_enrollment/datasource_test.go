package graphBetaWindowsUpdatesDeviceEnrollment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesDeviceEnrollment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/device_enrollment"
	deviceEnrollmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/device_enrollment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaWindowsUpdatesDeviceEnrollment.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Lookup device by Entra device ID
func TestUnitDatasourceDeviceEnrollment_01_ByEntraDeviceId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterGetDeviceByIdSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_by_entra_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entra_device_id").HasValue("fb95f07d-9e73-411d-99ab-7eca3a5122b1"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("fb95f07d-9e73-411d-99ab-7eca3a5122b1"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.0.update_category").HasValue("feature"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.1.update_category").HasValue("quality"),
					check.That(dataSourceType+".test").Key("devices.0.errors.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 02: List all enrolled devices
func TestUnitDatasourceDeviceEnrollment_02_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterListAllDevicesSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("devices.1.id").HasValue("90b91efa-6d46-42cd-ad4d-381831773a85"),
					check.That(dataSourceType+".test").Key("devices.1.enrollments.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("devices.2.id").HasValue("0ee3eb63-caf3-44ce-9769-b83188cc683d"),
					check.That(dataSourceType+".test").Key("devices.2.enrollments.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 03: Filter devices by quality update category
func TestUnitDatasourceDeviceEnrollment_03_FilterByUpdateCategory(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterListAllDevicesSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_filter_by_update_category.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("update_category").HasValue("quality"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("3"),
				),
			},
		},
	})
}

// Test 04: List devices with custom OData filter
func TestUnitDatasourceDeviceEnrollment_04_WithODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterListDevicesWithFilterSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_with_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("odata_filter").HasValue("id eq 'fb95f07d-9e73-411d-99ab-7eca3a5122b1'"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("fb95f07d-9e73-411d-99ab-7eca3a5122b1"),
				),
			},
		},
	})
}

// Test 05: Device with registration errors
func TestUnitDatasourceDeviceEnrollment_05_DeviceWithErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterGetDeviceWithRegistrationErrorMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_device_with_errors.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entra_device_id").HasValue("0ee3eb63-caf3-44ce-9769-b83188cc683d"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("0ee3eb63-caf3-44ce-9769-b83188cc683d"),
					check.That(dataSourceType+".test").Key("devices.0.errors.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.errors.0.error_code").HasValue("AzureADDeviceRegistrationError"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 06: Filter by feature update category
func TestUnitDatasourceDeviceEnrollment_06_FilterFeatureCategory(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterListAllDevicesSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_filter_feature_category.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("update_category").HasValue("feature"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("983f03cd-03cd-983f-cd03-3f98cd033f98"),
				),
			},
		},
	})
}

// Test 07: Filter by driver update category
func TestUnitDatasourceDeviceEnrollment_07_FilterDriverCategory(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterListAllDevicesSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_filter_driver_category.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("update_category").HasValue("driver"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("90b91efa-6d46-42cd-ad4d-381831773a85"),
				),
			},
		},
	})
}

// Test 08: Lookup device by device name
func TestUnitDatasourceDeviceEnrollment_08_ByDeviceName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterGetDeviceByNameSuccessMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_by_device_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_name").HasValue("TEST-DEVICE-001"),
					check.That(dataSourceType+".test").Key("devices.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("devices.0.id").HasValue("fb95f07d-9e73-411d-99ab-7eca3a5122b1"),
					check.That(dataSourceType+".test").Key("devices.0.enrollments.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 09: Error scenario - device not found
func TestUnitDatasourceDeviceEnrollment_09_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceEnrollmentMocks.RegisterGetDeviceByIdErrorMock()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_by_entra_device_id.tf"),
				ExpectError: regexp.MustCompile("device not found|NotFound|404"),
			},
		},
	})
}

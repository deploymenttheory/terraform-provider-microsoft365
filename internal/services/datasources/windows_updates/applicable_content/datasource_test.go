package graphBetaWindowsUpdatesApplicableContent_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesApplicableContent "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/applicable_content"
	applicableContentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/applicable_content/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaWindowsUpdatesApplicableContent.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *applicableContentMocks.ApplicableContentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicableContentMock := &applicableContentMocks.ApplicableContentMock{}
	applicableContentMock.RegisterMocks()
	return mockClient, applicableContentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *applicableContentMocks.ApplicableContentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicableContentMock := &applicableContentMocks.ApplicableContentMock{}
	applicableContentMock.RegisterErrorMocks()
	return mockClient, applicableContentMock
}

func TestUnitDatasourceApplicableContent_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").HasValue("f660d844-30b7-46e4-a6cf-47e36164d3cb"),
					check.That(dataSourceType+".test").Key("applicable_content.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry_id").Exists(),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry.id").Exists(),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry.display_name").Exists(),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry.manufacturer").Exists(),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry.version").Exists(),
					check.That(dataSourceType+".test").Key("applicable_content.0.matched_devices.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplicableContent_02_DriverUpdates(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_driver_updates.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").HasValue("f660d844-30b7-46e4-a6cf-47e36164d3cb"),
					check.That(dataSourceType+".test").Key("catalog_entry_type").HasValue("driver"),
					check.That(dataSourceType+".test").Key("applicable_content.#").HasValue("2"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.0.catalog_entry.driver_class"),
				),
			},
		},
	})
}

func TestUnitDatasourceApplicableContent_03_DisplayDrivers(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_display_drivers.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").HasValue("f660d844-30b7-46e4-a6cf-47e36164d3cb"),
					check.That(dataSourceType+".test").Key("catalog_entry_type").HasValue("driver"),
					check.That(dataSourceType+".test").Key("driver_class").HasValue("Display"),
					check.That(dataSourceType+".test").Key("applicable_content.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("applicable_content.0.catalog_entry.driver_class").HasValue("Display"),
				),
			},
		},
	})
}

func TestUnitDatasourceApplicableContent_04_WithMatchesOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_with_matches_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").HasValue("f660d844-30b7-46e4-a6cf-47e36164d3cb"),
					check.That(dataSourceType+".test").Key("include_no_matches").HasValue("false"),
					check.That(dataSourceType+".test").Key("applicable_content.#").HasValue("2"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.0.matched_devices.#"),
				),
			},
		},
	})
}

func TestUnitDatasourceApplicableContent_05_DeviceSpecific(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_device_specific.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").HasValue("f660d844-30b7-46e4-a6cf-47e36164d3cb"),
					check.That(dataSourceType+".test").Key("device_id").HasValue("fb95f07d-9e73-411d-99ab-7eca3a5122b1"),
					check.That(dataSourceType+".test").Key("applicable_content.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitDatasourceApplicableContent_06_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicableContentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicableContentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_basic.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}

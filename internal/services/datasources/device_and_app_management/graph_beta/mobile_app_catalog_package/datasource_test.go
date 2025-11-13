package graphBetaMobileAppCatalogPackage_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	MobileAppCatalogPackageMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app_catalog_package/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	packagesMock := &MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock{}
	packagesMock.RegisterMocks()
	return mockClient, packagesMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	packagesMock := &MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock{}
	packagesMock.RegisterErrorMocks()
	return mockClient, packagesMock
}

func TestMobileAppCatalogPackageDataSource_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "filter_type", "all"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.#", "4"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.product_display_name", "7-Zip"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.publisher_display_name", "Igor Pavlov"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.version_display_name", "25.01"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.branch_display_name", "7-Zip (x64)"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.applicable_architectures", "x64"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.locales.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.locales.0", "mul"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.0.package_auto_update_capable", "false"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.1.product_display_name", "3CXPhone for Windows"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.1.publisher_display_name", "3CX"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.1.package_auto_update_capable", "true"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "filter_value", "3a6307ef-6991-faf1-01e1-35e1557287aa"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.0.id", "5af1ade9-6966-3608-7e04-848252e29681"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.0.product_id", "3a6307ef-6991-faf1-01e1-35e1557287aa"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.0.product_display_name", "7-Zip"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.0.publisher_display_name", "Igor Pavlov"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ByProductName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByProductName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "filter_type", "product_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "filter_value", "7-Zip"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "items.0.product_display_name", "7-Zip"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ByPublisherName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByPublisherName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "filter_type", "publisher_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "filter_value", "Igor Pavlov"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "items.0.publisher_display_name", "Igor Pavlov"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_filter", "productDisplayName eq '7-Zip'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_orderby", "productDisplayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_search", "\"productDisplayName:7-Zip\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "items.#", "2"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "items.0.product_display_name", "7-Zip"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "items.1.product_display_name", "3CXPhone for Windows"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced", "odata_select", "id,productId,productDisplayName,publisherDisplayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced", "odata_top", "10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced", "odata_skip", "0"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced", "items.#", "4"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_filter", "productDisplayName eq '7-Zip'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_orderby", "productDisplayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_search", "\"productDisplayName:7-Zip\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_select", "id,productId,productDisplayName,publisherDisplayName,versionDisplayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_top", "5"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "odata_skip", "0"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_comprehensive", "items.#", "2"),
				),
			},
		},
	})
}

func TestMobileAppCatalogPackageDataSource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

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

func testConfigByProductName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_product_name.tf")
	if err != nil {
		panic("failed to load 03_by_product_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByPublisherName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_by_publisher_name.tf")
	if err != nil {
		panic("failed to load 04_by_publisher_name config: " + err.Error())
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

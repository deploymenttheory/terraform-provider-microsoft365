package graphBetaMobileAppCatalogPackage_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMobileAppCatalogPackageDataSource_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcceptanceConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "filter_type", "all"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all", "items.#"),
				),
			},
		},
	})
}

func TestAccMobileAppCatalogPackageDataSource_ById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcceptanceConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "filter_value", "3a6307ef-6991-faf1-01e1-35e1557287aa"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id", "items.#", "1"),
				),
			},
		},
	})
}

func TestAccMobileAppCatalogPackageDataSource_ByProductName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcceptanceConfigByProductName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "filter_type", "product_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "filter_value", "7-Zip"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name", "items.#"),
				),
			},
		},
	})
}

func TestAccMobileAppCatalogPackageDataSource_ByPublisherName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcceptanceConfigByPublisherName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "filter_type", "publisher_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "filter_value", "Microsoft"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name", "items.#"),
				),
			},
		},
	})
}

func TestAccMobileAppCatalogPackageDataSource_ODataFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAcceptanceConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_filter", "startswith(publisherDisplayName, 'Microsoft')"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "odata_top", "10"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter", "items.#"),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAcceptanceConfigAll() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_all.tf")
	if err != nil {
		panic("failed to load 01_all acceptance config: " + err.Error())
	}
	return acceptanceTestConfig
}

func testAcceptanceConfigById() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_by_id.tf")
	if err != nil {
		panic("failed to load 02_by_id acceptance config: " + err.Error())
	}
	return acceptanceTestConfig
}

func testAcceptanceConfigByProductName() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_by_product_name.tf")
	if err != nil {
		panic("failed to load 03_by_product_name acceptance config: " + err.Error())
	}
	return acceptanceTestConfig
}

func testAcceptanceConfigByPublisherName() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_by_publisher_name.tf")
	if err != nil {
		panic("failed to load 04_by_publisher_name acceptance config: " + err.Error())
	}
	return acceptanceTestConfig
}

func testAcceptanceConfigODataFilter() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/05_odata_filter.tf")
	if err != nil {
		panic("failed to load 05_odata_filter acceptance config: " + err.Error())
	}
	return acceptanceTestConfig
}

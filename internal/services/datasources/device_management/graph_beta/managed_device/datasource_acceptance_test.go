package graphBetaManagedDevice_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManagedDeviceDataSource_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.all", "filter_type", "all"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.all", "items.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.device_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.all", "items.0.operating_system"),
				),
			},
		},
	})
}

// this test will fail if a real dvice in intune does not have the word "DESKTOP" in the device name.
// since this is a lab, this is highly likely to fail based upon what test devices are in lab at any
// given point in time.
//
// func TestAccManagedDeviceDataSource_ByDeviceName(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		ExternalProviders: map[string]resource.ExternalProvider{
// 			"azuread": {
// 				Source:            "hashicorp/azuread",
// 				VersionConstraint: ">= 2.47.0",
// 			},
// 			"random": {
// 				Source:            "hashicorp/random",
// 				VersionConstraint: ">= 3.7.2",
// 			},
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccConfigByDeviceName(),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "filter_type", "device_name"),
// 					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "filter_value", "DESKTOP"),
// 					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.by_device_name", "items.#"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccManagedDeviceDataSource_ODataFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "odata_filter", "operatingSystem eq 'Windows'"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.odata_filter", "items.#"),
				),
			},
		},
	})
}

func TestAccManagedDeviceDataSource_ODataAdvanced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_filter", "operatingSystem eq 'Windows'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_orderby", "deviceName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "odata_select", "id,deviceName,operatingSystem,complianceState"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.odata_advanced", "items.#"),
				),
			},
		},
	})
}

func TestAccManagedDeviceDataSource_ODataComprehensive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_filter", "operatingSystem eq 'Windows'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_top", "50"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "odata_orderby", "lastSyncDateTime desc"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_device_management_managed_device.odata_comprehensive", "items.#"),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigAll() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_all.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigByDeviceName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_by_device_name.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataFilter() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_odata_filter.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataAdvanced() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_odata_advanced.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataComprehensive() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/05_odata_comprehensive.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

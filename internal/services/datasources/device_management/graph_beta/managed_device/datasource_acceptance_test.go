package graphBetaManagedDevice_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccDatasourceManagedDevice_01_ListAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_02_ByDeviceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_device_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_name").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_03_ODataSimpleFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_odata_simple_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows'"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_04_ODataAndFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_odata_and_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows' and complianceState eq 'compliant'"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_05_ODataStartswithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_odata_startswith_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("startswith(deviceName, 'WIN-')"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_06_ODataOrFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("06_odata_or_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'iOS' or operatingSystem eq 'Android'"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_07_ByDeviceId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("07_by_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_id").Exists(),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").Exists(),
					check.That(dataSourceType+".test").Key("items.0.device_name").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_08_ByOperatingSystem(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("08_by_operating_system.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("operating_system").HasValue("Windows"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_09_ByOsAndVersion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("09_by_os_and_version.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("operating_system").Exists(),
					check.That(dataSourceType+".test").Key("os_version").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_10_ByAzureADDeviceId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("10_by_azure_ad_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("azure_ad_device_id").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_11_BySerialNumber(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("11_by_serial_number.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("serial_number").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceManagedDevice_12_ByUserPrincipalName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("12_by_user_principal_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("user_principal_name").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func testCheckItemsCountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		itemsCount, ok := rs.Primary.Attributes["items.#"]
		if !ok {
			return fmt.Errorf("items.# attribute not found in resource %s", resourceName)
		}

		count, err := strconv.Atoi(itemsCount)
		if err != nil {
			return fmt.Errorf("items.# is not a valid number: %s", itemsCount)
		}

		if count < 0 {
			return fmt.Errorf("items.# cannot be negative: %d", count)
		}

		return nil
	}
}

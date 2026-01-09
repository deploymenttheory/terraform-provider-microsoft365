package graphBetaManagedDevice_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to load acceptance test Terraform configurations
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccManagedDeviceDataSource_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					testCheckItemsCountExists(dataSourceType+".all"),
					// Note: Not checking specific items.0.* fields as test environment may have zero managed devices
				),
			},
		},
	})
}

func TestAccManagedDeviceDataSource_ByDeviceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_device_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_device_name").Key("filter_type").HasValue("device_name"),
					check.That(dataSourceType+".by_device_name").Key("filter_value").HasValue("DESKTOP"),
					testCheckItemsCountExists(dataSourceType+".by_device_name"),
					// Note: Not checking specific items.0.* fields as filtered results may return zero devices depending on lab state
				),
			},
		},
	})
}

func TestAccManagedDeviceDataSource_ODataFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("operatingSystem eq 'Windows'"),
					testCheckItemsCountExists(dataSourceType+".odata_filter"),
					// Note: Not checking specific items.0.* fields as filtered results may return zero devices
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
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_advanced").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_advanced").Key("odata_filter").HasValue("operatingSystem eq 'Windows'"),
					check.That(dataSourceType+".odata_advanced").Key("odata_orderby").HasValue("deviceName"),
					check.That(dataSourceType+".odata_advanced").Key("odata_select").HasValue("id,deviceName,operatingSystem,complianceState"),
					testCheckItemsCountExists(dataSourceType+".odata_advanced"),
					// Note: Not checking specific items.0.* fields as filtered results may return zero devices
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
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_comprehensive").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_filter").HasValue("operatingSystem eq 'Windows'"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_top").HasValue("50"),
					check.That(dataSourceType+".odata_comprehensive").Key("odata_orderby").HasValue("lastSyncDateTime desc"),
					testCheckItemsCountExists(dataSourceType+".odata_comprehensive"),
					// Note: Not checking specific items.0.* fields as filtered results may return zero devices
				),
			},
		},
	})
}

// Helper function to check for the condition that items.# exists and is >= 0
// this is used in scenarios where the number of items is not known in advance
// and we are asserting against a lab intune environment where the number of items may be zero.
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

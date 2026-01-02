package utilityLicensingServicePlanReference_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLicensingServicePlanReferenceDataSource_SearchByProductName(t *testing.T) {
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
				Config: testAccConfigSearchByProductName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_products.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.product_name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.string_id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.guid").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.service_plans_included.#").IsSet(),
				),
			},
		},
	})
}

func TestAccLicensingServicePlanReferenceDataSource_SearchByStringId(t *testing.T) {
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
				Config: testAccConfigSearchByStringId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_products.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.string_id").HasValue("SPE_E3_RPA1"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.product_name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.guid").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.service_plans_included.#").IsSet(),
				),
			},
		},
	})
}

func TestAccLicensingServicePlanReferenceDataSource_SearchByServicePlan(t *testing.T) {
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
				Config: testAccConfigSearchByServicePlan(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.guid").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.0.product_name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.0.string_id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.0.guid").Exists(),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigSearchByProductName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_search_by_product_name.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigSearchByStringId() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_search_by_string_id.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigSearchByServicePlan() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_search_by_service_plan.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

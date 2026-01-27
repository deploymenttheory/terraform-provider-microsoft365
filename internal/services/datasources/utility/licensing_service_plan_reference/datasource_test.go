package utilityLicensingServicePlanReference_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityMicrosoft365ServicePlan "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/licensing_service_plan_reference"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityMicrosoft365ServicePlan.DataSourceName
)

func TestUnitDatasourceLicensingServicePlanReference_01_SearchByProductName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByProductName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_products.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.product_name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.string_id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.guid").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_02_SearchByStringId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByStringId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").HasValue("string_id:ENTERPRISEPACK"),
					check.That("data."+dataSourceType+".test").Key("matching_products.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.string_id").HasValue("ENTERPRISEPACK"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.guid").HasValue("6fd2c87f-b296-42f0-b197-1e91e994b900"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.service_plans_included.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_03_SearchByGuid(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByGuid(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").HasValue("guid:6fd2c87f-b296-42f0-b197-1e91e994b900"),
					check.That("data."+dataSourceType+".test").Key("matching_products.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.guid").HasValue("6fd2c87f-b296-42f0-b197-1e91e994b900"),
					check.That("data."+dataSourceType+".test").Key("matching_products.0.string_id").HasValue("ENTERPRISEPACK"),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_04_SearchByServicePlanName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByServicePlanName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.name").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.guid").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_05_SearchByServicePlanId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByServicePlanId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.id").Exists(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_06_SearchByServicePlanGuid(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchByServicePlanGuid(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("id").HasValue("service_plan_guid:113feb6c-3fe4-4440-bddc-54d774bf0318"),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.guid").HasValue("113feb6c-3fe4-4440-bddc-54d774bf0318"),
					check.That("data."+dataSourceType+".test").Key("matching_service_plans.0.included_in_skus.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_07_MultipleParametersProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMultipleParamsProvided(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_08_NoParametersProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigNoParamsProvided(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestUnitDatasourceLicensingServicePlanReference_09_InvalidGuidFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidGuidFormat(),
				ExpectError: regexp.MustCompile("must be a valid GUID"),
			},
		},
	})
}

// Configuration functions
func testConfigSearchByProductName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_search_by_product_name.tf")
	if err != nil {
		panic("failed to load search_by_product_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchByStringId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_search_by_string_id.tf")
	if err != nil {
		panic("failed to load search_by_string_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchByGuid() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_search_by_guid.tf")
	if err != nil {
		panic("failed to load search_by_guid config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchByServicePlanName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_search_by_service_plan_name.tf")
	if err != nil {
		panic("failed to load search_by_service_plan_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchByServicePlanId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_search_by_service_plan_id.tf")
	if err != nil {
		panic("failed to load search_by_service_plan_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchByServicePlanGuid() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_search_by_service_plan_guid.tf")
	if err != nil {
		panic("failed to load search_by_service_plan_guid config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMultipleParamsProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/07_multiple_params_provided.tf")
	if err != nil {
		panic("failed to load multiple_params_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigNoParamsProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/08_no_params_provided.tf")
	if err != nil {
		panic("failed to load no_params_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidGuidFormat() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/09_invalid_guid_format.tf")
	if err != nil {
		panic("failed to load invalid_guid_format config: " + err.Error())
	}
	return unitTestConfig
}

package graphBetaMobileAppRelationship_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	relMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app_relationship/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const dataSourceType = "data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship"

// TestMobileAppRelationshipDataSource_All tests fetching all mobile app relationships
func TestUnitDatasourceMobileAppRelationship_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relMocks.RegisterMobileAppRelationshipMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").HasValue("3"),

					// Verify first relationship
					check.That(dataSourceType+".all").Key("items.0.id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".all").Key("items.0.target_id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.target_display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.source_id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.source_display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.target_type").Exists(),

					// Verify specific values from mock data
					check.That(dataSourceType+".all").Key("items.0.target_display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".all").Key("items.0.source_display_name").HasValue("Company Portal"),
					check.That(dataSourceType+".all").Key("items.0.target_type").HasValue("parent"),
				),
			},
		},
	})
}

// TestMobileAppRelationshipDataSource_ById tests fetching a specific mobile app relationship by ID
func TestUnitDatasourceMobileAppRelationship_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relMocks.RegisterMobileAppRelationshipMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),

					// Verify the specific relationship
					check.That(dataSourceType+".by_id").Key("items.0.id").HasValue("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					check.That(dataSourceType+".by_id").Key("items.0.target_id").HasValue("app-target-001"),
					check.That(dataSourceType+".by_id").Key("items.0.target_display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".by_id").Key("items.0.target_display_version").HasValue("120.0.2210.91"),
					check.That(dataSourceType+".by_id").Key("items.0.target_publisher").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".by_id").Key("items.0.target_publisher_display_name").HasValue("Microsoft"),
					check.That(dataSourceType+".by_id").Key("items.0.source_id").HasValue("app-source-001"),
					check.That(dataSourceType+".by_id").Key("items.0.source_display_name").HasValue("Company Portal"),
					check.That(dataSourceType+".by_id").Key("items.0.source_display_version").HasValue("5.0.5954.0"),
					check.That(dataSourceType+".by_id").Key("items.0.source_publisher_display_name").HasValue("Microsoft"),
					check.That(dataSourceType+".by_id").Key("items.0.target_type").HasValue("parent"),
				),
			},
		},
	})
}

// TestMobileAppRelationshipDataSource_BySourceId tests filtering relationships by source app ID
func TestUnitDatasourceMobileAppRelationship_03_BySourceId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relMocks.RegisterMobileAppRelationshipMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBySourceId(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_source_id").Key("filter_type").HasValue("source_id"),
					check.That(dataSourceType+".by_source_id").Key("filter_value").HasValue("app-source-001"),
					check.That(dataSourceType+".by_source_id").Key("items.#").HasValue("1"),

					// Verify the relationship has the correct source_id
					check.That(dataSourceType+".by_source_id").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_source_id").Key("items.0.source_id").HasValue("app-source-001"),
					check.That(dataSourceType+".by_source_id").Key("items.0.target_display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".by_source_id").Key("items.0.source_display_name").HasValue("Company Portal"),
				),
			},
		},
	})
}

// TestMobileAppRelationshipDataSource_ByTargetId tests filtering relationships by target app ID
func TestUnitDatasourceMobileAppRelationship_04_ByTargetId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relMocks.RegisterMobileAppRelationshipMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByTargetId(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_target_id").Key("filter_type").HasValue("target_id"),
					check.That(dataSourceType+".by_target_id").Key("filter_value").HasValue("app-target-001"),
					check.That(dataSourceType+".by_target_id").Key("items.#").HasValue("1"),

					// Verify the relationship has the correct target_id
					check.That(dataSourceType+".by_target_id").Key("items.0.target_id").HasValue("app-target-001"),
					check.That(dataSourceType+".by_target_id").Key("items.0.target_display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".by_target_id").Key("items.0.source_display_name").HasValue("Company Portal"),
				),
			},
		},
	})
}

// TestMobileAppRelationshipDataSource_ODataFilter tests using OData filter queries
func TestUnitDatasourceMobileAppRelationship_05_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relMocks.RegisterMobileAppRelationshipMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("sourceId eq 'app-source-001'"),
					check.That(dataSourceType+".odata_filter").Key("items.#").HasValue("2"),

					// Verify filtered results
					check.That(dataSourceType+".odata_filter").Key("items.0.id").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.source_id").HasValue("app-source-001"),
					check.That(dataSourceType+".odata_filter").Key("items.0.target_display_name").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.target_type").Exists(),
				),
			},
		},
	})
}

// Helper functions to load test configurations
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

func testConfigBySourceId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_source_id.tf")
	if err != nil {
		panic("failed to load 03_by_source_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByTargetId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_by_target_id.tf")
	if err != nil {
		panic("failed to load 04_by_target_id config: " + err.Error())
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

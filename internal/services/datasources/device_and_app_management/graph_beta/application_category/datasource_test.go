package graphBetaApplicationCategory_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	catmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/application_category/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const dataSourceType = "data.microsoft365_graph_beta_device_and_app_management_application_category"

// TestApplicationCategoryDataSource_All tests fetching all application categories
func TestUnitDatasourceApplicationCategory_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catmocks.RegisterApplicationCategoryMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").HasValue("4"),

					// Verify first category
					check.That(dataSourceType+".all").Key("items.0.id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".all").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.last_modified_date_time").Exists(),

					// Verify specific categories exist in the results
					check.That(dataSourceType+".all").Key("items.0.display_name").HasValue("Productivity"),
					check.That(dataSourceType+".all").Key("items.1.display_name").HasValue("Business"),
					check.That(dataSourceType+".all").Key("items.2.display_name").HasValue("Finance"),
					check.That(dataSourceType+".all").Key("items.3.display_name").HasValue("Education"),
				),
			},
		},
	})
}

// TestApplicationCategoryDataSource_ById tests fetching a specific application category by ID
func TestUnitDatasourceApplicationCategory_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catmocks.RegisterApplicationCategoryMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("5b0e1e8d-7a5c-4f3a-9c2d-1e4f5a6b7c8d"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),

					// Verify the specific category
					check.That(dataSourceType+".by_id").Key("items.0.id").HasValue("5b0e1e8d-7a5c-4f3a-9c2d-1e4f5a6b7c8d"),
					check.That(dataSourceType+".by_id").Key("items.0.display_name").HasValue("Productivity"),
					check.That(dataSourceType+".by_id").Key("items.0.last_modified_date_time").HasValue("2024-01-15T10:30:00Z"),
				),
			},
		},
	})
}

// TestApplicationCategoryDataSource_ByDisplayName tests filtering categories by display name
func TestUnitDatasourceApplicationCategory_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catmocks.RegisterApplicationCategoryMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_display_name").Key("filter_type").HasValue("display_name"),
					check.That(dataSourceType+".by_display_name").Key("filter_value").HasValue("Business"),
					check.That(dataSourceType+".by_display_name").Key("items.#").HasValue("1"),

					// Verify category contains "Business" in the name
					check.That(dataSourceType+".by_display_name").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)Business`)),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").HasValue("Business"),
					check.That(dataSourceType+".by_display_name").Key("items.0.last_modified_date_time").Exists(),
				),
			},
		},
	})
}

// TestApplicationCategoryDataSource_ODataFilter tests using OData filter queries
func TestUnitDatasourceApplicationCategory_04_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catmocks.RegisterApplicationCategoryMockResponders()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("startswith(displayName, 'Business')"),
					check.That(dataSourceType+".odata_filter").Key("items.#").HasValue("2"),

					// Verify filtered results
					check.That(dataSourceType+".odata_filter").Key("items.0.id").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.last_modified_date_time").Exists(),

					// Verify specific categories returned
					check.That(dataSourceType+".odata_filter").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`^Business`)),
					check.That(dataSourceType+".odata_filter").Key("items.1.display_name").MatchesRegex(regexp.MustCompile(`^Business`)),
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

func testConfigByDisplayName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_display_name.tf")
	if err != nil {
		panic("failed to load 03_by_display_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataFilter() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_odata_filter.tf")
	if err != nil {
		panic("failed to load 04_odata_filter config: " + err.Error())
	}
	return unitTestConfig
}

package graphBetaMobileApp_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	MobileAppMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Shared constant for both unit and acceptance tests
const dataSourceType = "data.microsoft365_graph_beta_device_and_app_management_mobile_app"

func setupMockEnvironment() (*mocks.Mocks, *MobileAppMocks.MobileAppsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appsMock := &MobileAppMocks.MobileAppsMock{}
	appsMock.RegisterMocks()
	return mockClient, appsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *MobileAppMocks.MobileAppsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appsMock := &MobileAppMocks.MobileAppsMock{}
	appsMock.RegisterErrorMocks()
	return mockClient, appsMock
}

// Test 01: Get all mobile apps - comprehensive field validation
func TestUnitDatasourceMobileApp_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").HasValue("4"),

					// ============================================
					// Item 0: Microsoft Edge (Win32 LOB App)
					// ============================================
					check.That(dataSourceType+".all").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".all").Key("items.0.display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".all").Key("items.0.description").HasValue("Microsoft Edge browser for Windows"),
					check.That(dataSourceType+".all").Key("items.0.publisher").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".all").Key("items.0.developer").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".all").Key("items.0.owner").HasValue("IT Department"),
					check.That(dataSourceType+".all").Key("items.0.notes").HasValue("Latest stable version"),
					check.That(dataSourceType+".all").Key("items.0.created_date_time").HasValue("2024-01-15T10:30:00Z"),
					check.That(dataSourceType+".all").Key("items.0.last_modified_date_time").HasValue("2024-01-15T10:30:00Z"),
					check.That(dataSourceType+".all").Key("items.0.privacy_information_url").HasValue("https://privacy.microsoft.com"),
					check.That(dataSourceType+".all").Key("items.0.information_url").HasValue("https://www.microsoft.com/edge"),
					check.That(dataSourceType+".all").Key("items.0.is_featured").HasValue("true"),
					check.That(dataSourceType+".all").Key("items.0.is_assigned").HasValue("true"),
					check.That(dataSourceType+".all").Key("items.0.publishing_state").HasValue("published"),
					check.That(dataSourceType+".all").Key("items.0.upload_state").HasValue("1"),
					check.That(dataSourceType+".all").Key("items.0.dependent_app_count").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.superseding_app_count").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.superseded_app_count").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.role_scope_tag_ids.#").HasValue("1"),
					check.That(dataSourceType+".all").Key("items.0.role_scope_tag_ids.0").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.categories.#").HasValue("1"),
					check.That(dataSourceType+".all").Key("items.0.categories.0").HasValue("Productivity"),

					// ============================================
					// Item 1: Adobe Acrobat Reader (macOS PKG App)
					// ============================================
					check.That(dataSourceType+".all").Key("items.1.id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(dataSourceType+".all").Key("items.1.display_name").HasValue("Adobe Acrobat Reader"),
					check.That(dataSourceType+".all").Key("items.1.publisher").HasValue("Adobe Inc."),
					check.That(dataSourceType+".all").Key("items.1.is_featured").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.1.categories.#").HasValue("0"),

					// ============================================
					// Item 2: Microsoft Teams (iOS Store App)
					// ============================================
					check.That(dataSourceType+".all").Key("items.2.id").HasValue("00000000-0000-0000-0000-000000000003"),
					check.That(dataSourceType+".all").Key("items.2.display_name").HasValue("Microsoft Teams"),
					check.That(dataSourceType+".all").Key("items.2.publisher").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".all").Key("items.2.categories.#").HasValue("1"),
					check.That(dataSourceType+".all").Key("items.2.categories.0").HasValue("Communication"),

					// ============================================
					// Item 3: Slack (Android Managed Store App)
					// ============================================
					check.That(dataSourceType+".all").Key("items.3.id").HasValue("00000000-0000-0000-0000-000000000004"),
					check.That(dataSourceType+".all").Key("items.3.display_name").HasValue("Slack"),
					check.That(dataSourceType+".all").Key("items.3.publisher").HasValue("Slack Technologies"),
					check.That(dataSourceType+".all").Key("items.3.is_assigned").HasValue("false"),
				),
			},
		},
	})
}

// Test 02: Get mobile app by ID
func TestUnitDatasourceMobileApp_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),

					// Complete field validation for single item
					check.That(dataSourceType+".by_id").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".by_id").Key("items.0.display_name").HasValue("Microsoft Edge"),
					check.That(dataSourceType+".by_id").Key("items.0.description").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.publisher").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".by_id").Key("items.0.developer").HasValue("Microsoft Corporation"),
					check.That(dataSourceType+".by_id").Key("items.0.owner").HasValue("IT Department"),
					check.That(dataSourceType+".by_id").Key("items.0.privacy_information_url").HasValue("https://privacy.microsoft.com"),
					check.That(dataSourceType+".by_id").Key("items.0.information_url").HasValue("https://www.microsoft.com/edge"),
					check.That(dataSourceType+".by_id").Key("items.0.is_featured").HasValue("true"),
					check.That(dataSourceType+".by_id").Key("items.0.publishing_state").HasValue("published"),
					check.That(dataSourceType+".by_id").Key("items.0.is_assigned").HasValue("true"),
					check.That(dataSourceType+".by_id").Key("items.0.categories.#").HasValue("1"),
					check.That(dataSourceType+".by_id").Key("items.0.categories.0").HasValue("Productivity"),
				),
			},
		},
	})
}

// Test 03: Get by display name
func TestUnitDatasourceMobileApp_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_display_name").Key("filter_type").HasValue("display_name"),
					check.That(dataSourceType+".by_display_name").Key("filter_value").HasValue("Microsoft"),
					check.That(dataSourceType+".by_display_name").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
					check.That(dataSourceType+".by_display_name").Key("items.1.display_name").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
				),
			},
		},
	})
}

// Test 04: Get by publisher name
func TestUnitDatasourceMobileApp_04_ByPublisherName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByPublisherName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_publisher").Key("filter_type").HasValue("publisher_name"),
					check.That(dataSourceType+".by_publisher").Key("filter_value").HasValue("Microsoft"),
					check.That(dataSourceType+".by_publisher").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".by_publisher").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
					check.That(dataSourceType+".by_publisher").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.1.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
				),
			},
		},
	})
}

// Test 05: OData filter
func TestUnitDatasourceMobileApp_05_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("startswith(publisher, 'Microsoft')"),
					check.That(dataSourceType+".odata_filter").Key("odata_top").HasValue("10"),
					check.That(dataSourceType+".odata_filter").Key("items.#").Exists(),
				),
			},
		},
	})
}

// Test 06: With app type filter
func TestUnitDatasourceMobileApp_06_WithAppTypeFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithAppTypeFilter(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".win32_apps").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".win32_apps").Key("app_type_filter").HasValue("win32LobApp"),
					check.That(dataSourceType+".win32_apps").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".win32_apps").Key("items.0.display_name").HasValue("Microsoft Edge"),
				),
			},
		},
	})
}

// Terraform config helpers
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

func testConfigWithAppTypeFilter() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_with_app_type_filter.tf")
	if err != nil {
		panic("failed to load 06_with_app_type_filter config: " + err.Error())
	}
	return unitTestConfig
}

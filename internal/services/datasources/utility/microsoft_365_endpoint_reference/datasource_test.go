package utilityMicrosoft365EndpointReference_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityMicrosoft365EndpointReference "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/microsoft_365_endpoint_reference"
	microsoft365EndpointReferenceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/microsoft_365_endpoint_reference/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityMicrosoft365EndpointReference.DataSourceName
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *microsoft365EndpointReferenceMocks.Microsoft365EndpointReferenceMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	microsoft365EndpointsMock := &microsoft365EndpointReferenceMocks.Microsoft365EndpointReferenceMock{}
	microsoft365EndpointsMock.RegisterMocks()

	return mockClient, microsoft365EndpointsMock
}

// Helper functions to load each test configuration
func testConfigWorldwide() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_worldwide.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigFilterExchange() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_filter_exchange.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigFilterOptimize() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_filter_optimize.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigRequiredOnly() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_required_only.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigExpressRoute() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_expressroute.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigUSGovDoD() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_usgov_dod.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigUSGovGCCHigh() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_usgov_gcchigh.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigChina() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_china.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMultipleFilters() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_multiple_filters.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Test cases
func TestMicrosoft365EndpointReferenceDataSource_Worldwide(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWorldwide(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// Should have at least 6 endpoints from our mock data
					check.That("data."+dataSourceType+".test").Key("endpoints.#").MatchesRegex(regexp.MustCompile(`^[6-9]$|^[1-9][0-9]+$`)),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_FilterByServiceArea(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigFilterExchange(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("service_areas.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("service_areas.*").ContainsTypeSetElement("Exchange"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_FilterByCategory(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigFilterOptimize(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("categories.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("categories.*").ContainsTypeSetElement("Optimize"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_RequiredOnly(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigRequiredOnly(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("required_only").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_ExpressRoute(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExpressRoute(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("express_route").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_USGovDoD(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUSGovDoD(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("usgov-dod"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_USGovGCCHigh(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUSGovGCCHigh(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("usgov-gcchigh"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_China(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigChina(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("china"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_MultipleFilters(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMultipleFilters(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("service_areas.#").HasValue("2"),
					check.That("data."+dataSourceType+".test").Key("service_areas.*").ContainsTypeSetElement("Exchange"),
					check.That("data."+dataSourceType+".test").Key("service_areas.*").ContainsTypeSetElement("Skype"),
					check.That("data."+dataSourceType+".test").Key("categories.#").HasValue("2"),
					check.That("data."+dataSourceType+".test").Key("categories.*").ContainsTypeSetElement("Optimize"),
					check.That("data."+dataSourceType+".test").Key("categories.*").ContainsTypeSetElement("Allow"),
					check.That("data."+dataSourceType+".test").Key("required_only").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_InvalidInstance(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoint_reference" "test" {
  instance = "invalid"
}
`,
				ExpectError: regexp.MustCompile(`Attribute instance value must be one of`),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_InvalidServiceArea(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoint_reference" "test" {
  instance      = "worldwide"
  service_areas = ["Invalid"]
}
`,
				ExpectError: regexp.MustCompile(`Attribute service_areas\[Value\("Invalid"\)\] value must be one of`),
			},
		},
	})
}

func TestMicrosoft365EndpointReferenceDataSource_InvalidCategory(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoint_reference" "test" {
  instance   = "worldwide"
  categories = ["Invalid"]
}
`,
				ExpectError: regexp.MustCompile(`Attribute categories\[Value\("Invalid"\)\] value must be one of`),
			},
		},
	})
}

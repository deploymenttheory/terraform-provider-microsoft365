package utilityMicrosoft365Endpoints_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	microsoft365EndpointsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/microsoft_365_endpoints/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *microsoft365EndpointsMocks.Microsoft365EndpointsMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	microsoft365EndpointsMock := &microsoft365EndpointsMocks.Microsoft365EndpointsMock{}
	microsoft365EndpointsMock.RegisterMocks()

	return mockClient, microsoft365EndpointsMock
}

// testCheckExists is a basic check to ensure the datasource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
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
func TestMicrosoft365EndpointsDataSource_Worldwide(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// Should have at least 6 endpoints from our mock data
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#", regexp.MustCompile(`^[6-9]$|^[1-9][0-9]+$`)),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_FilterByServiceArea(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.*", "Exchange"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_FilterByCategory(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.*", "Optimize"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_RequiredOnly(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "required_only", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_ExpressRoute(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "express_route", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_USGovDoD(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "usgov-dod"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_USGovGCCHigh(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "usgov-gcchigh"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_China(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "china"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_MultipleFilters(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.#", "2"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.*", "Exchange"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.*", "Skype"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.#", "2"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.*", "Optimize"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.*", "Allow"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "required_only", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
				),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_InvalidInstance(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoints" "test" {
  instance = "invalid"
}
`,
				ExpectError: regexp.MustCompile(`Attribute instance value must be one of`),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_InvalidServiceArea(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoints" "test" {
  instance      = "worldwide"
  service_areas = ["Invalid"]
}
`,
				ExpectError: regexp.MustCompile(`Attribute service_areas\[Value\("Invalid"\)\] value must be one of`),
			},
		},
	})
}

func TestMicrosoft365EndpointsDataSource_InvalidCategory(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, microsoft365EndpointsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer microsoft365EndpointsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_microsoft_365_endpoints" "test" {
  instance   = "worldwide"
  categories = ["Invalid"]
}
`,
				ExpectError: regexp.MustCompile(`Attribute categories\[Value\("Invalid"\)\] value must be one of`),
			},
		},
	})
}

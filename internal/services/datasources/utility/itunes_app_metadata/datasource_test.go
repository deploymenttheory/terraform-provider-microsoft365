package itunes_app_metadata_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	itunes_app_metadata "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/itunes_app_metadata"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/itunes_app_metadata/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = itunes_app_metadata.DataSourceName
)

// Helper functions to return the test configurations by reading from files
func testConfigFirefox() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_firefox.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigOffice() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_office.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigError() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_error.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigEmpty() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "datasource_empty.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.ItunesAppMetadataMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	itunesMock := &localMocks.ItunesAppMetadataMock{}
	itunesMock.RegisterMocks()

	return mockClient, itunesMock
}

// TestUnitItunesAppMetadataDataSource_Firefox tests fetching Firefox app metadata
func TestUnitItunesAppMetadataDataSource_Firefox(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigFirefox(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".firefox").Key("id").HasValue("us_firefox"),
					check.That("data."+dataSourceType+".firefox").Key("search_term").HasValue("firefox"),
					check.That("data."+dataSourceType+".firefox").Key("country_code").HasValue("us"),
					check.That("data."+dataSourceType+".firefox").Key("results.#").HasValue("1"),
					check.That("data."+dataSourceType+".firefox").Key("results.0.track_id").HasValue("989804926"),
					check.That("data."+dataSourceType+".firefox").Key("results.0.track_name").HasValue("Firefox Fast & Private Browser"),
					check.That("data."+dataSourceType+".firefox").Key("results.0.bundle_id").HasValue("org.mozilla.ios.Firefox"),
					check.That("data."+dataSourceType+".firefox").Key("results.0.seller_name").HasValue("Mozilla Corporation"),
					check.That("data."+dataSourceType+".firefox").Key("results.0.version").HasValue("140.2"),
				),
			},
		},
	})
}

// TestUnitItunesAppMetadataDataSource_Office tests fetching Office app metadata
func TestUnitItunesAppMetadataDataSource_Office(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigOffice(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".office").Key("id").HasValue("gb_office"),
					check.That("data."+dataSourceType+".office").Key("search_term").HasValue("office"),
					check.That("data."+dataSourceType+".office").Key("country_code").HasValue("gb"),
					check.That("data."+dataSourceType+".office").Key("results.#").HasValue("1"),
					check.That("data."+dataSourceType+".office").Key("results.0.track_id").HasValue("541164041"),
					check.That("data."+dataSourceType+".office").Key("results.0.track_name").HasValue("Microsoft 365 Copilot"),
					check.That("data."+dataSourceType+".office").Key("results.0.bundle_id").HasValue("com.microsoft.officemobile"),
					check.That("data."+dataSourceType+".office").Key("results.0.seller_name").HasValue("Microsoft Corporation"),
					check.That("data."+dataSourceType+".office").Key("results.0.version").HasValue("2.98.4"),
				),
			},
		},
	})
}

// TestUnitItunesAppMetadataDataSource_Empty tests fetching empty app metadata results
func TestUnitItunesAppMetadataDataSource_Empty(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigEmpty(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".empty").Key("id").HasValue("us_empty"),
					check.That("data."+dataSourceType+".empty").Key("search_term").HasValue("empty"),
					check.That("data."+dataSourceType+".empty").Key("country_code").HasValue("us"),
					check.That("data."+dataSourceType+".empty").Key("results.#").HasValue("0"),
				),
			},
		},
	})
}

// TestUnitItunesAppMetadataDataSource_Error tests error handling
func TestUnitItunesAppMetadataDataSource_Error(t *testing.T) {
	// Set up mock environment
	_, itunesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks instead of standard mocks
	itunesMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Error from iTunes API"),
			},
		},
	})
}

package itunes_app_metadata_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/itunes_app_metadata/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Helper functions to return the test configurations by reading from files
func testConfigFirefox() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_firefox.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigOffice() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_office.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigError() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_error.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigEmpty() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_empty.tf"))
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "id", "us_firefox"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "search_term", "firefox"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "country_code", "us"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.0.track_id", "989804926"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.0.track_name", "Firefox Fast & Private Browser"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.0.bundle_id", "org.mozilla.ios.Firefox"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.0.seller_name", "Mozilla Corporation"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.firefox", "results.0.version", "140.2"),
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "id", "gb_office"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "search_term", "office"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "country_code", "gb"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.0.track_id", "541164041"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.0.track_name", "Microsoft 365 Copilot"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.0.bundle_id", "com.microsoft.officemobile"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.0.seller_name", "Microsoft Corporation"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.office", "results.0.version", "2.98.4"),
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.empty", "id", "us_empty"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.empty", "search_term", "empty"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.empty", "country_code", "us"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_itunes_app_metadata.empty", "results.#", "0"),
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

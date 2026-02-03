package graphBetaApplication_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplication "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/applications/graph_beta/application"
	applicationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/applications/graph_beta/application/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = graphBetaApplication.DataSourceName
)

func setupMockEnvironment() (*mocks.Mocks, *applicationMocks.ApplicationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appMock := &applicationMocks.ApplicationMock{}
	appMock.RegisterMocks()
	return mockClient, appMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *applicationMocks.ApplicationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appMock := &applicationMocks.ApplicationMock{}
	appMock.RegisterErrorMocks()
	return mockClient, appMock
}

func TestUnitDatasourceApplication_01_ByObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_id").Key("object_id").HasValue("a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"),
					check.That("data."+dataSourceType+".by_id").Key("id").HasValue("a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"),
					check.That("data."+dataSourceType+".by_id").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".by_id").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".by_id").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_id").Key("publisher_domain").HasValue("example.com"),
					check.That("data."+dataSourceType+".by_id").Key("description").Exists(),
					check.That("data."+dataSourceType+".by_id").Key("notes").Exists(),
					check.That("data."+dataSourceType+".by_id").Key("identifier_uris.#").Exists(),
					check.That("data."+dataSourceType+".by_id").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".by_id").Key("info.logo_url").Exists(),
					check.That("data."+dataSourceType+".by_id").Key("web.home_page_url").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_02_ByAppId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_app_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_app_id").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".by_app_id").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".by_app_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_app_id").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("identifier_uris.#").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_display_name").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".by_display_name").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".by_display_name").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_display_name").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("identifier_uris.#").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_04_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_filter").Key("odata_query").HasValue("displayName eq 'Test Application' and signInAudience eq 'AzureADMyOrg'"),
					check.That("data."+dataSourceType+".odata_filter").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".odata_filter").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".odata_filter").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_filter").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("identifier_uris.#").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_05_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_advanced").Key("odata_query").HasValue("appId eq '12345678-1234-1234-1234-123456789012'"),
					check.That("data."+dataSourceType+".odata_advanced").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".odata_advanced").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_advanced").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("identifier_uris.#").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_06_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_comprehensive").Key("odata_query").HasValue("tags/any(t:t eq 'MyCustomTag') and signInAudience eq 'AzureADMyOrg'"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("display_name").HasValue("Test Application"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceApplication_07_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_by_object_id.tf"),
				ExpectError: regexp.MustCompile("Forbidden - 403"),
			},
		},
	})
}

// Helper function to load unit test Terraform configs
func loadUnitTestTerraform(filename string) string {
	unitTestConfig, err := helpers.ParseHCLFile(fmt.Sprintf("tests/terraform/unit/%s", filename))
	if err != nil {
		panic(fmt.Sprintf("failed to load unit test config: %s", err.Error()))
	}
	return unitTestConfig
}

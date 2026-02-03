package graphBetaServicePrincipal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	servicePrincipalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/applications/graph_beta/service_principal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *servicePrincipalMocks.ServicePrincipalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	spMock := &servicePrincipalMocks.ServicePrincipalMock{}
	spMock.RegisterMocks()
	return mockClient, spMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *servicePrincipalMocks.ServicePrincipalMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	spMock := &servicePrincipalMocks.ServicePrincipalMock{}
	spMock.RegisterErrorMocks()
	return mockClient, spMock
}

func TestUnitDatasourceServicePrincipal_01_ByObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_id").Key("object_id").HasValue("3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"),
					check.That("data."+dataSourceType+".by_id").Key("id").HasValue("3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"),
					check.That("data."+dataSourceType+".by_id").Key("display_name").HasValue("Microsoft Intune SCCM Connector"),
					check.That("data."+dataSourceType+".by_id").Key("app_id").HasValue("63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					check.That("data."+dataSourceType+".by_id").Key("app_display_name").HasValue("Microsoft Intune SCCM Connector"),
					check.That("data."+dataSourceType+".by_id").Key("publisher_name").HasValue("Microsoft Services"),
					check.That("data."+dataSourceType+".by_id").Key("account_enabled").HasValue("true"),
					check.That("data."+dataSourceType+".by_id").Key("service_principal_type").HasValue("Application"),
					check.That("data."+dataSourceType+".by_id").Key("app_role_assignment_required").HasValue("false"),
					check.That("data."+dataSourceType+".by_id").Key("sign_in_audience").HasValue("AzureADMultipleOrgs"),
					check.That("data."+dataSourceType+".by_id").Key("preferred_single_sign_on_mode").HasValue("oidc"),
					check.That("data."+dataSourceType+".by_id").Key("homepage").HasValue("https://intune.microsoft.com"),
					check.That("data."+dataSourceType+".by_id").Key("login_url").HasValue("https://login.microsoftonline.com"),
					check.That("data."+dataSourceType+".by_id").Key("logout_url").HasValue("https://login.microsoftonline.com/logout"),
					check.That("data."+dataSourceType+".by_id").Key("notes").HasValue("Production service principal"),
					check.That("data."+dataSourceType+".by_id").Key("reply_urls.#").HasValue("2"),
					check.That("data."+dataSourceType+".by_id").Key("service_principal_names.#").HasValue("2"),
					check.That("data."+dataSourceType+".by_id").Key("tags.#").HasValue("2"),
					check.That("data."+dataSourceType+".by_id").Key("notification_email_addresses.#").HasValue("2"),
					check.That("data."+dataSourceType+".by_id").Key("saml_single_sign_on_settings.relay_state").HasValue("https://example.com/relay"),
					check.That("data."+dataSourceType+".by_id").Key("verified_publisher.display_name").HasValue("Microsoft Corporation"),
					check.That("data."+dataSourceType+".by_id").Key("info.support_url").HasValue("https://support.microsoft.com"),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_02_ByAppId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_app_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_app_id").Key("app_id").HasValue("63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					check.That("data."+dataSourceType+".by_app_id").Key("display_name").HasValue("Microsoft Intune SCCM Connector"),
					check.That("data."+dataSourceType+".by_app_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("app_display_name").HasValue("Microsoft Intune SCCM Connector"),
					check.That("data."+dataSourceType+".by_app_id").Key("publisher_name").HasValue("Microsoft Services"),
					check.That("data."+dataSourceType+".by_app_id").Key("account_enabled").HasValue("true"),
					check.That("data."+dataSourceType+".by_app_id").Key("service_principal_type").HasValue("Application"),
					check.That("data."+dataSourceType+".by_app_id").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("service_principal_names.#").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("tags.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data." + dataSourceType + ".by_display_name").Key("display_name").HasValue("Microsoft Intune SCCM Connector"),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_04_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_filter").Key("odata_query").HasValue("preferredSingleSignOnMode ne 'notSupported' and displayName eq 'Microsoft Intune'"),
					check.That("data."+dataSourceType+".odata_filter").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("sign_in_audience").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_05_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_advanced").Key("odata_query").HasValue("servicePrincipalType eq 'Application' and accountEnabled eq true"),
					check.That("data."+dataSourceType+".odata_advanced").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("service_principal_type").HasValue("Application"),
					check.That("data."+dataSourceType+".odata_advanced").Key("account_enabled").HasValue("true"),
					check.That("data."+dataSourceType+".odata_advanced").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_06_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_comprehensive").Key("odata_query").HasValue("preferredSingleSignOnMode eq 'saml' and servicePrincipalType eq 'Application'"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("service_principal_type").HasValue("Application"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("preferred_single_sign_on_mode").HasValue("saml"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("saml_single_sign_on_settings.relay_state").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceServicePrincipal_07_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

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

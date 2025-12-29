package graphBetaServicePrincipal_test

import (
	"regexp"
	"testing"

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

func TestServicePrincipalDataSource_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "filter_type", "all"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.#", "4"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.0.display_name", "Microsoft Intune SCCM Connector"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.0.app_id", "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.1.display_name", "Microsoft Intune Service Discovery"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.1.app_id", "9cb77803-d937-493e-9a3b-4b49de3f5a74"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.2.display_name", "Microsoft Intune Web Company Portal"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "items.3.display_name", "MMD Intune Partner Sync"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "filter_value", "3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "items.0.id", "3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "items.0.display_name", "Microsoft Intune SCCM Connector"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "items.0.app_id", "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_id", "items.0.publisher_name", "Microsoft Services"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ByAppId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByAppId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "filter_type", "app_id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "filter_value", "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.0.app_id", "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.0.display_name", "Microsoft Intune SCCM Connector"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "filter_type", "display_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "filter_value", "Microsoft Intune SCCM Connector"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.0.display_name", "Microsoft Intune SCCM Connector"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_filter", "preferredSingleSignOnMode ne 'notSupported'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_search", "\"displayName:intune\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.#", "2"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.0.display_name", "Microsoft Intune SCCM Connector"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.1.display_name", "Microsoft Intune Service Discovery"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_select", "appId,displayName,publisherName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_top", "10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_skip", "0"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "items.#", "4"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_filter", "preferredSingleSignOnMode ne 'notSupported'"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_search", "\"displayName:intune\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_select", "id,appId,displayName,publisherName,servicePrincipalType"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_top", "5"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_skip", "0"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "items.#", "2"),
				),
			},
		},
	})
}

func TestServicePrincipalDataSource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, spMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer spMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigAll(),
				ExpectError: regexp.MustCompile("Forbidden - 403"),
			},
		},
	})
}

// Configuration functions
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

func testConfigByAppId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_app_id.tf")
	if err != nil {
		panic("failed to load 03_by_app_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByDisplayName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_by_display_name.tf")
	if err != nil {
		panic("failed to load 04_by_display_name config: " + err.Error())
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

func testConfigODataAdvanced() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_odata_advanced.tf")
	if err != nil {
		panic("failed to load 06_odata_advanced config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataComprehensive() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/07_odata_comprehensive.tf")
	if err != nil {
		panic("failed to load 07_odata_comprehensive config: " + err.Error())
	}
	return unitTestConfig
}

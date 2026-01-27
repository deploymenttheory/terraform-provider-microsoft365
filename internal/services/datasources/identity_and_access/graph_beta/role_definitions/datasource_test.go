package graphBetaRoleDefinitions_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleDefinitionsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/role_definitions/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *roleDefinitionsMocks.RoleDefinitionsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	rdMock := &roleDefinitionsMocks.RoleDefinitionsMock{}
	rdMock.RegisterMocks()
	return mockClient, rdMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *roleDefinitionsMocks.RoleDefinitionsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	rdMock := &roleDefinitionsMocks.RoleDefinitionsMock{}
	rdMock.RegisterErrorMocks()
	return mockClient, rdMock
}

func TestUnitDatasourceRoleDefinitions_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "filter_type", "all"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.#", "5"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.display_name", "Global Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.id", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.is_privileged", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.1.display_name", "Privileged Role Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.1.id", "e8611ab8-c189-46e8-94e1-60213ab1f814"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.2.display_name", "Security Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.3.display_name", "Application Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.4.display_name", "Conditional Access Administrator"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "filter_value", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.id", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.display_name", "Global Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.template_id", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.is_built_in", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.is_privileged", "true"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "filter_type", "display_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "filter_value", "Global Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.0.display_name", "Global Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.0.id", "62e90394-69f5-4237-9190-012177145e10"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_04_ODataFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "odata_filter", "isPrivileged eq true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "items.#", "2"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "items.0.is_privileged", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "items.1.is_privileged", "true"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_05_ODataAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_filter", "isBuiltIn eq true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_select", "id,displayName,isPrivileged"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "items.#", "2"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_06_ODataComprehensive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_filter", "isPrivileged eq true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_top", "5"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_skip", "0"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "items.#", "2"),
				),
			},
		},
	})
}

func TestUnitDatasourceRoleDefinitions_07_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, rdMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer rdMock.CleanupMockState()

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

func testConfigODataAdvanced() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/05_odata_advanced.tf")
	if err != nil {
		panic("failed to load 05_odata_advanced config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigODataComprehensive() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/06_odata_comprehensive.tf")
	if err != nil {
		panic("failed to load 06_odata_comprehensive config: " + err.Error())
	}
	return unitTestConfig
}

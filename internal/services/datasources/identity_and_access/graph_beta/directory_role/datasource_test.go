package graphBetaDirectoryRole_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDirectoryRole "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/directory_role"
	directoryRoleMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/directory_role/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaDirectoryRole.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *directoryRoleMocks.DirectoryRoleMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	drMock := &directoryRoleMocks.DirectoryRoleMock{}
	drMock.RegisterMocks()
	return mockClient, drMock
}


func TestUnitDatasourceDirectoryRole_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, drMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer drMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("aaaaaaaa-0001-0000-0000-000000000000"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("User Administrator"),
					check.That(dataSourceType+".test").Key("items.0.role_template_id").HasValue("fe930be7-5e62-47db-91af-98c3a49a38b1"),
					check.That(dataSourceType+".test").Key("items.1.display_name").HasValue("Helpdesk Administrator"),
					check.That(dataSourceType+".test").Key("items.2.display_name").HasValue("Global Administrator"),
				),
			},
		},
	})
}

func TestUnitDatasourceDirectoryRole_02_ByID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, drMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer drMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("aaaaaaaa-0001-0000-0000-000000000000"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("User Administrator"),
					check.That(dataSourceType+".test").Key("items.0.role_template_id").HasValue("fe930be7-5e62-47db-91af-98c3a49a38b1"),
				),
			},
		},
	})
}

func TestUnitDatasourceDirectoryRole_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, drMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer drMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("aaaaaaaa-0001-0000-0000-000000000000"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("User Administrator"),
					check.That(dataSourceType+".test").Key("items.0.role_template_id").HasValue("fe930be7-5e62-47db-91af-98c3a49a38b1"),
				),
			},
		},
	})
}


package graphBetaGroup_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/groups/graph_beta/group"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaGroup.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()
	return mockClient, groupMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *groupMocks.GroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterErrorMocks()
	return mockClient, groupMock
}

func TestUnitDatasourceGroup_01_ByObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("display_name").HasValue("IT Security Team"),
					check.That(dataSourceType+".test").Key("security_enabled").HasValue("true"),
					check.That(dataSourceType+".test").Key("mail_enabled").HasValue("false"),
				),
			},
		},
	})
}

func TestUnitDatasourceGroup_02_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").HasValue("IT Security Team"),
					check.That(dataSourceType+".test").Key("id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("security_enabled").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceGroup_03_ByDisplayNameWithSecurityFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name_with_security_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").HasValue("IT Security Team"),
					check.That(dataSourceType+".test").Key("security_enabled").HasValue("true"),
					check.That(dataSourceType+".test").Key("mail_enabled").HasValue("false"),
				),
			},
		},
	})
}

func TestUnitDatasourceGroup_04_ByMailNickname(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_mail_nickname.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("mail_nickname").HasValue("finance"),
					check.That(dataSourceType+".test").Key("id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(dataSourceType+".test").Key("display_name").HasValue("Finance Team"),
				),
			},
		},
	})
}

func TestUnitDatasourceGroup_05_ByODataQuery(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_by_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").HasValue("IT Security Team"),
					check.That(dataSourceType+".test").Key("id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".test").Key("security_enabled").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceGroup_06_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_by_object_id.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|insufficient|authorisation"),
			},
		},
	})
}

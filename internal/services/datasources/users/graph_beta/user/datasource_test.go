package graphBetaUser_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	userMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/users/graph_beta/user/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *userMocks.UserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	uMock := &userMocks.UserMock{}
	uMock.RegisterMocks()
	return mockClient, uMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *userMocks.UserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	uMock := &userMocks.UserMock{}
	uMock.RegisterErrorMocks()
	return mockClient, uMock
}

func TestUnitDatasourceUser_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-TEST-USER-001"),
					check.That(dataSourceType+".test").Key("items.1.display_name").HasValue("DT-TEST-USER-002"),
					check.That(dataSourceType+".test").Key("items.2.display_name").HasValue("DT-TEST-USER-003"),
					check.That(dataSourceType+".test").Key("items.2.user_type").HasValue("Guest"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_02_ByObjectId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-TEST-USER-001"),
					check.That(dataSourceType+".test").Key("items.0.user_principal_name").HasValue("dt.test.user001@contoso.com"),
					check.That(dataSourceType+".test").Key("items.0.employee_id").HasValue("EMP-0001"),
					check.That(dataSourceType+".test").Key("items.0.account_enabled").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").HasValue("DT-TEST-USER-001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-TEST-USER-001"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_04_ByEmployeeId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_employee_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("employee_id").HasValue("EMP-0001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.employee_id").HasValue("EMP-0001"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("DT-TEST-USER-001"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_05_ByGivenName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_by_given_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("given_name").HasValue("Test"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.given_name").HasValue("Test"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_06_ByUserPrincipalName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_by_user_principal_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("user_principal_name").HasValue("dt.test.user001@contoso.com"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.user_principal_name").HasValue("dt.test.user001@contoso.com"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_07_ByOnPremisesImmutableId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_by_on_premises_immutable_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("on_premises_immutable_id").HasValue("IMMUTABLE-0001"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.on_premises_immutable_id").HasValue("IMMUTABLE-0001"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_08_ByOnPremisesDistinguishedName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_by_on_premises_distinguished_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("on_premises_distinguished_name").HasValue("CN=Test UserOne,OU=Users,DC=contoso,DC=com"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.on_premises_distinguished_name").HasValue("CN=Test UserOne,OU=Users,DC=contoso,DC=com"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_09_ByODataQuery(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("accountEnabled eq true and userType eq 'Member'"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.user_type").HasValue("Member"),
					check.That(dataSourceType+".test").Key("items.1.user_type").HasValue("Member"),
				),
			},
		},
	})
}

func TestUnitDatasourceUser_10_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, uMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer uMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("02_by_object_id.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|insufficient|privileges"),
			},
		},
	})
}

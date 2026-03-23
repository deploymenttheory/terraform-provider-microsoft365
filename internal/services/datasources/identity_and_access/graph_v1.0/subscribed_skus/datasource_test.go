package graphSubscribedSkus_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphSubscribedSkus "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_v1.0/subscribed_skus"
	subscribedSkusMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_v1.0/subscribed_skus/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphSubscribedSkus.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *subscribedSkusMocks.SubscribedSkusMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	skuMock := &subscribedSkusMocks.SubscribedSkusMock{}
	skuMock.RegisterMocks()
	return mockClient, skuMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *subscribedSkusMocks.SubscribedSkusMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	skuMock := &subscribedSkusMocks.SubscribedSkusMock{}
	skuMock.RegisterErrorMocks()
	return mockClient, skuMock
}

func TestUnitDatasourceSubscribedSkus_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("48a80680-7326-48cd-9935-b556b81d3a4e_c7df2760-2c81-4ef7-b578-5b5392b571df"),
					check.That(dataSourceType+".test").Key("items.0.sku_part_number").HasValue("ENTERPRISEPREMIUM"),
					check.That(dataSourceType+".test").Key("items.0.account_name").HasValue("Contoso Corporation"),
					check.That(dataSourceType+".test").Key("items.1.sku_part_number").HasValue("CRMSTANDARD"),
					check.That(dataSourceType+".test").Key("items.2.sku_part_number").HasValue("AAD_PREMIUM"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_02_BySkuId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_sku_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("sku_id").HasValue("48a80680-7326-48cd-9935-b556b81d3a4e_c7df2760-2c81-4ef7-b578-5b5392b571df"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("48a80680-7326-48cd-9935-b556b81d3a4e_c7df2760-2c81-4ef7-b578-5b5392b571df"),
					check.That(dataSourceType+".test").Key("items.0.sku_part_number").HasValue("ENTERPRISEPREMIUM"),
					check.That(dataSourceType+".test").Key("items.0.account_name").HasValue("Contoso Corporation"),
					check.That(dataSourceType+".test").Key("items.0.applies_to").HasValue("User"),
					check.That(dataSourceType+".test").Key("items.0.consumed_units").HasValue("14"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_03_BySkuPartNumber(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_sku_part_number.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("sku_part_number").HasValue("ENTERPRISEPREMIUM"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.sku_part_number").HasValue("ENTERPRISEPREMIUM"),
					check.That(dataSourceType+".test").Key("items.0.account_name").HasValue("Contoso Corporation"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_04_ByAccountId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_by_account_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("account_id").HasValue("f97aeefc-af85-414d-8ae4-b457f90efc40"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.account_id").HasValue("f97aeefc-af85-414d-8ae4-b457f90efc40"),
					check.That(dataSourceType+".test").Key("items.0.account_name").HasValue("Contoso Corporation"),
					check.That(dataSourceType+".test").Key("items.1.account_id").HasValue("f97aeefc-af85-414d-8ae4-b457f90efc40"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_05_ByAccountName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_by_account_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("account_name").HasValue("Contoso"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.account_name").HasValue("Contoso Corporation"),
					check.That(dataSourceType+".test").Key("items.1.account_name").HasValue("Contoso Corporation"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_06_ByAppliesToUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_by_applies_to_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("applies_to").HasValue("User"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.applies_to").HasValue("User"),
					check.That(dataSourceType+".test").Key("items.1.applies_to").HasValue("User"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_07_ByAppliesToCompany(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_by_applies_to_company.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("applies_to").HasValue("Company"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.applies_to").HasValue("Company"),
					check.That(dataSourceType+".test").Key("items.0.sku_part_number").HasValue("AAD_PREMIUM"),
				),
			},
		},
	})
}

func TestUnitDatasourceSubscribedSkus_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, skuMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer skuMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_list_all.tf"),
				ExpectError: regexp.MustCompile("Internal Server Error|500"),
			},
		},
	})
}

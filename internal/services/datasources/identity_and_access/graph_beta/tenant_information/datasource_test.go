package graphBetaTenantInformation_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTenantInformation "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/tenant_information"
	tenantInformationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/tenant_information/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = graphBetaTenantInformation.DataSourceName
)

func setupMockEnvironment() (*mocks.Mocks, *tenantInformationMocks.TenantInformationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tiMock := &tenantInformationMocks.TenantInformationMock{}
	tiMock.RegisterMocks()
	return mockClient, tiMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *tenantInformationMocks.TenantInformationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tiMock := &tenantInformationMocks.TenantInformationMock{}
	tiMock.RegisterErrorMocks()
	return mockClient, tiMock
}

func TestTenantInformationDataSource_ByTenantId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tiMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tiMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByTenantId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_tenant_id").Key("filter_type").HasValue("tenant_id"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("filter_value").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("tenant_id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("display_name").HasValue("Deployment Theory"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("default_domain_name").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestTenantInformationDataSource_ByDomainName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tiMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tiMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDomainName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_domain_name").Key("filter_type").HasValue("domain_name"),
					check.That("data."+dataSourceType+".by_domain_name").Key("filter_value").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_domain_name").Key("tenant_id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_domain_name").Key("display_name").HasValue("Deployment Theory"),
					check.That("data."+dataSourceType+".by_domain_name").Key("default_domain_name").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_domain_name").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestTenantInformationDataSource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tiMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tiMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigByTenantId(),
				ExpectError: regexp.MustCompile("Forbidden - 403"),
			},
		},
	})
}

// Configuration functions
func testConfigByTenantId() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_by_tenant_id.tf")
	if err != nil {
		panic("failed to load 01_by_tenant_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByDomainName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_by_domain_name.tf")
	if err != nil {
		panic("failed to load 02_by_domain_name config: " + err.Error())
	}
	return unitTestConfig
}

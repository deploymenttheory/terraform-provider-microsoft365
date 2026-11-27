package graphBetaTenantInformation_test

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTenantInformation "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/tenants/tenant_information"
	tenantInformationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/tenants/tenant_information/mocks"
)

var dataSourceType = graphBetaTenantInformation.DataSourceName

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

func TestUnitDatasourceTenantInformation_01_ByTenantId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByTenantId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_tenant_id").Key("items.0.tenant_id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("items.0.display_name").HasValue("Deployment Theory"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("items.0.default_domain_name").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
				),
			},
		},
	})
}

func TestUnitDatasourceTenantInformation_02_ByDomainName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByDomainName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_domain_name").Key("items.0.tenant_id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_domain_name").Key("items.0.display_name").HasValue("Deployment Theory"),
					check.That("data."+dataSourceType+".by_domain_name").Key("items.0.default_domain_name").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_domain_name").Key("id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
				),
			},
		},
	})
}

func TestUnitDatasourceTenantInformation_03_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
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

func TestUnitDatasourceTenantInformation_InvalidLookup(t *testing.T) {
	for _, tc := range []struct{ name, attributes string }{
		{"missing", ""},
		{"both", `tenant_id = "6babcaad-604b-40ac-a9d7-9fd97c0b779f" domain_name = "contoso.com"`},
		{"invalid_id", `tenant_id = "invalid"`},
		{"invalid_domain", `domain_name = "contoso..com"`},
		{"empty_domain", `domain_name = ""`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{{
				Config:      "data \"" + dataSourceType + "\" \"invalid\" {\n" + strings.ReplaceAll(tc.attributes, " domain_name", "\ndomain_name") + "\n}",
				ExpectError: regexp.MustCompile("Invalid Attribute|Invalid Attribute Combination"),
			}}})
		})
	}
}

func TestUnitDatasourceTenantInformation_MissingResult(t *testing.T) {
	for _, tc := range []struct {
		name           string
		status         int
		body, expected string
	}{
		{"not_found", 404, "{\"error\":{\"code\":\"Request_ResourceNotFound\",\"message\":\"Tenant not found\"}}", "404|not found"},
		{"empty_object", 200, "{}", "Invalid tenant information response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				"https://graph.microsoft.com/beta/tenantRelationships/findTenantInformationByTenantId(tenantId='6babcaad-604b-40ac-a9d7-9fd97c0b779f')",
				httpmock.NewStringResponder(tc.status, tc.body).HeaderSet(http.Header{"Content-Type": []string{"application/json"}}),
			)
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{{
				Config: testConfigByTenantId(), ExpectError: regexp.MustCompile(tc.expected),
			}}})
		})
	}
}

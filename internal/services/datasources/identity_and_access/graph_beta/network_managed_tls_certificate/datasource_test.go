package graphBetaNetworkManagedTLSCertificate_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkManagedTLSCertificate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/network_managed_tls_certificate"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitDatasourceNetworkManagedTLSCertificate_01_ByID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	response, err := helpers.ParseJSONFile("tests/responses/get_managed_tls_certificate.json")
	if err != nil {
		t.Fatalf("failed to load response fixture: %v", err)
	}
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates/00000000-0000-0000-0000-000000000301",
		httpmock.NewStringResponder(200, response).HeaderSet(http.Header{"Content-Type": []string{"application/json"}}),
	)

	dataSourceType := graphBetaNetworkManagedTLSCertificate.DataSourceName
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: dataSourceTestConfig(t),
			Check: resource.ComposeTestCheckFunc(
				check.That("data."+dataSourceType+".test").Key("id").HasValue("00000000-0000-0000-0000-000000000301"),
				check.That("data."+dataSourceType+".test").Key("certificate_authority_id").HasValue("00000000-0000-0000-0000-000000000301"),
				check.That("data."+dataSourceType+".test").Key("status").HasValue("active"),
				check.That("data."+dataSourceType+".test").Key("certificate").HasValue("-----BEGIN CERTIFICATE-----\nUNIT-TEST\n-----END CERTIFICATE-----"),
				check.That("data."+dataSourceType+".test").Key("validity_end_date_time").HasValue("2036-09-03T12:44:24Z"),
			),
		}},
	})
}

func TestUnitDatasourceNetworkManagedTLSCertificate_02_NotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates/00000000-0000-0000-0000-000000000301",
		httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Microsoft-managed TLS certificate authority was not found"}}`),
	)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      dataSourceTestConfig(t),
			ExpectError: regexp.MustCompile(`Not Found - 404`),
		}},
	})
}

func dataSourceTestConfig(t *testing.T) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/datasource.tf")
	if err != nil {
		t.Fatalf("failed to load data source config: %v", err)
	}
	return config
}

package graphBetaNetworkManagedTLSCertificate_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaNetworkManagedTLSCertificate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_managed_tls_certificate"
	managedTLSCertificateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_managed_tls_certificate/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const unitTestResourceType = graphBetaNetworkManagedTLSCertificate.ResourceName

func setupMockEnvironment() (*mocks.Mocks, *managedTLSCertificateMocks.ManagedTLSCertificateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	certificateMock := &managedTLSCertificateMocks.ManagedTLSCertificateMock{}
	certificateMock.RegisterMocks()
	return mockClient, certificateMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *managedTLSCertificateMocks.ManagedTLSCertificateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	certificateMock := &managedTLSCertificateMocks.ManagedTLSCertificateMock{}
	certificateMock.RegisterErrorMocks()
	return mockClient, certificateMock
}

func TestUnitResourceNetworkManagedTLSCertificate_01_DefaultsAndImport(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceType+".test").Key("name").MatchesRegex(regexp.MustCompile(`^M-TLSi-[0-9a-z]{5}$`)),
					check.That(unitTestResourceType+".test").Key("common_name").HasValue("Microsoft Entra TLS Inspection Root CA"),
					check.That(unitTestResourceType+".test").Key("organization_name").HasValue("Microsoft"),
					check.That(unitTestResourceType+".test").Key("validity_months").HasValue("120"),
					check.That(unitTestResourceType+".test").Key("enabled").HasValue("false"),
					check.That(unitTestResourceType+".test").Key("status").HasValue("disabled"),
					check.That(unitTestResourceType+".test").Key("certificate").Exists(),
					check.That(unitTestResourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:      unitTestResourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceNetworkManagedTLSCertificate_02_CreateEnabledThenDisable(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceType+".test").Key("name").HasValue("M-TLSi-unit-enabled"),
					check.That(unitTestResourceType+".test").Key("common_name").HasValue("Terraform TLS Inspection Root CA"),
					check.That(unitTestResourceType+".test").Key("organization_name").HasValue("Contoso"),
					check.That(unitTestResourceType+".test").Key("validity_months").HasValue("60"),
					check.That(unitTestResourceType+".test").Key("enabled").HasValue("true"),
					check.That(unitTestResourceType+".test").Key("status").HasValue("active"),
				),
			},
			{
				Config: testConfigHelper("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceType+".test").Key("enabled").HasValue("false"),
					check.That(unitTestResourceType+".test").Key("status").HasValue("disabled"),
				),
			},
		},
	})
}

func TestUnitResourceNetworkManagedTLSCertificate_03_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      testConfigHelper("resource_minimal.tf"),
			ExpectError: regexp.MustCompile(`Bad Request - 400`),
		}},
	})
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

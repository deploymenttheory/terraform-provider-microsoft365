package graphBetaApplicationCertificateCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationCertificateCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_certificate_credential"
	certificateCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_certificate_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationCertificateCredential.ResourceName

	// testResource is the test resource implementation
	testResource = graphBetaApplicationCertificateCredential.ApplicationCertificateCredentialTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *certificateCredentialMocks.ApplicationCertificateCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	certificateMock := &certificateCredentialMocks.ApplicationCertificateCredentialMock{}
	certificateMock.RegisterMocks()

	return mockClient, certificateMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *certificateCredentialMocks.ApplicationCertificateCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	certificateMock := &certificateCredentialMocks.ApplicationCertificateCredentialMock{}
	certificateMock.RegisterErrorMocks()

	return mockClient, certificateMock
}

func TestUnitResourceApplicationCertificateCredential_01_Base64(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_base64.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_base64").Key("application_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_base64").Key("display_name").HasValue("unit-test-certificate-base64"),
					check.That(resourceType+".test_base64").Key("encoding").HasValue("base64"),
					check.That(resourceType+".test_base64").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_base64").Key("usage").HasValue("Verify"),
					check.That(resourceType+".test_base64").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

func TestUnitResourceApplicationCertificateCredential_02_DER(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_der.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_der").Key("application_id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test_der").Key("display_name").HasValue("unit-test-certificate-der"),
					check.That(resourceType+".test_der").Key("encoding").HasValue("base64"),
					check.That(resourceType+".test_der").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_der").Key("usage").HasValue("Verify"),
					check.That(resourceType+".test_der").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

func TestUnitResourceApplicationCertificateCredential_03_HEX(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_03_hex.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_hex").Key("application_id").HasValue("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".test_hex").Key("display_name").HasValue("unit-test-certificate-hex"),
					check.That(resourceType+".test_hex").Key("encoding").HasValue("hex"),
					check.That(resourceType+".test_hex").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_hex").Key("usage").HasValue("Verify"),
					check.That(resourceType+".test_hex").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

func TestUnitResourceApplicationCertificateCredential_04_PEM(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_04_pem.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_pem").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_pem").Key("display_name").HasValue("unit-test-certificate-pem"),
					check.That(resourceType+".test_pem").Key("encoding").HasValue("pem"),
					check.That(resourceType+".test_pem").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_pem").Key("usage").HasValue("Verify"),
					check.That(resourceType+".test_pem").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			// Note: Import is not supported for certificate credentials
		},
	})
}

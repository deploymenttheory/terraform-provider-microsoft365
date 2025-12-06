package graphBetaAgentIdentityBlueprintCertificateCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentIdentityBlueprintCertificateCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_certificate_credential"
	certificateCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_certificate_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentityBlueprintCertificateCredential.ResourceName

	// testResource is the test resource implementation
	testResource = graphBetaAgentIdentityBlueprintCertificateCredential.AgentIdentityBlueprintCertificateCredentialTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *certificateCredentialMocks.AgentIdentityBlueprintCertificateCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	certificateMock := &certificateCredentialMocks.AgentIdentityBlueprintCertificateCredentialMock{}
	certificateMock.RegisterMocks()

	return mockClient, certificateMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *certificateCredentialMocks.AgentIdentityBlueprintCertificateCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	certificateMock := &certificateCredentialMocks.AgentIdentityBlueprintCertificateCredentialMock{}
	certificateMock.RegisterErrorMocks()

	return mockClient, certificateMock
}

func TestAgentIdentityBlueprintCertificateCredentialResource_PEM(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigPEM(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_pem").Key("blueprint_id").HasValue("11111111-1111-1111-1111-111111111111"),
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

func TestAgentIdentityBlueprintCertificateCredentialResource_DER(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigDER(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_der").Key("blueprint_id").HasValue("22222222-2222-2222-2222-222222222222"),
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

func TestAgentIdentityBlueprintCertificateCredentialResource_HEX(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHEX(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_hex").Key("blueprint_id").HasValue("33333333-3333-3333-3333-333333333333"),
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

func testConfigPEM() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_pem.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testConfigDER() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_der.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testConfigHEX() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_hex.tf")
	if err != nil {
		panic(err)
	}
	return content
}

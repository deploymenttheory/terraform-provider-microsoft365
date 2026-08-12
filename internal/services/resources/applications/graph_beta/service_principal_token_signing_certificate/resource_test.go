package graphBetaServicePrincipalTokenSigningCertificate_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalTokenSigningCertificate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_token_signing_certificate"
	certificateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_token_signing_certificate/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalTokenSigningCertificate.ResourceName

	// testResource is the test resource implementation for token signing certificates
	testResource = graphBetaServicePrincipalTokenSigningCertificate.ServicePrincipalTokenSigningCertificateTestResource{}
)

const testServicePrincipalID = "11111111-1111-1111-1111-111111111111"

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *certificateMocks.ServicePrincipalTokenSigningCertificateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	certificateMock := &certificateMocks.ServicePrincipalTokenSigningCertificateMock{}
	certificateMock.RegisterMocks()
	return mockClient, certificateMock
}

func TestUnitResourceServicePrincipalTokenSigningCertificate_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		// Destroy must remove the certificate's Sign/Verify key credentials and password
		// credential via the read-modify-write PATCH while retaining the unrelated
		// credential pair the mock seeds (1 key + 1 password).
		CheckDestroy: func(s *terraform.State) error {
			keyCredentials, passwordCredentials, ok := certificateMocks.GetCredentialCounts(testServicePrincipalID)
			if !ok {
				return fmt.Errorf("mock service principal not found after destroy")
			}
			if keyCredentials != 1 || passwordCredentials != 1 {
				return fmt.Errorf("expected exactly the unrelated credentials to remain after destroy (1 key, 1 password), got: %d key, %d password", keyCredentials, passwordCredentials)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("service_principal_id").HasValue(testServicePrincipalID),
					check.That(resourceType+".test").Key("display_name").HasValue("CN=Test SAML Signing"),
					check.That(resourceType+".test").Key("end_date_time").HasValue("2029-01-01T00:00:00Z"),
					check.That(resourceType+".test").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]{36}$`)),
					check.That(resourceType+".test").Key("thumbprint").MatchesRegex(regexp.MustCompile(`^[0-9a-f]{40}$`)),
					check.That(resourceType+".test").Key("start_date_time").Exists(),
					check.That(resourceType+".test").Key("value").Exists(),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^`+testServicePrincipalID+`/[0-9a-fA-F-]{36}$`)),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
				ImportStateIdFunc:       testAccImportStateIdFunc(resourceType + ".test"),
			},
		},
	})
}

func TestUnitResourceServicePrincipalTokenSigningCertificate_02_Defaults(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, certificateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer certificateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_defaults.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("CN=Microsoft Azure Federated SSO Certificate"),
					check.That(resourceType+".test").Key("end_date_time").Exists(),
					check.That(resourceType+".test").Key("thumbprint").Exists(),
				),
			},
			{
				// API defaults must not produce a diff on re-plan
				Config:   loadUnitTestTerraform("resource_02_defaults.tf"),
				PlanOnly: true,
			},
		},
	})
}

// testAccImportStateIdFunc returns a function that constructs the import ID
func testAccImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_principal_id"], rs.Primary.Attributes["key_id"]), nil
	}
}

package graphBetaApplicationFederatedIdentityCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationFederatedIdentityCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_federated_identity_credential"
	applicationFederatedIdentityCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_federated_identity_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationFederatedIdentityCredential.ResourceName

	// testResource is the test resource implementation for federated identity credentials
	testResource = graphBetaApplicationFederatedIdentityCredential.ApplicationFederatedIdentityCredentialTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *applicationFederatedIdentityCredentialMocks.ApplicationFederatedIdentityCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationFederatedIdentityCredentialMock := &applicationFederatedIdentityCredentialMocks.ApplicationFederatedIdentityCredentialMock{}
	applicationFederatedIdentityCredentialMock.RegisterMocks()
	return mockClient, applicationFederatedIdentityCredentialMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *applicationFederatedIdentityCredentialMocks.ApplicationFederatedIdentityCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationFederatedIdentityCredentialMock := &applicationFederatedIdentityCredentialMocks.ApplicationFederatedIdentityCredentialMock{}
	applicationFederatedIdentityCredentialMock.RegisterErrorMocks()
	return mockClient, applicationFederatedIdentityCredentialMock
}

func TestUnitResourceApplicationFederatedIdentityCredential_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("name").HasValue("unit-test-fic-minimal"),
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_minimal").Key("subject").HasValue("repo:octo-org/octo-repo:environment:Production"),
					check.That(resourceType+".test_minimal").Key("audiences.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_minimal"),
			},
		},
	})
}

func TestUnitResourceApplicationFederatedIdentityCredential_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_maximal").Key("name").HasValue("unit-test-fic-maximal"),
					check.That(resourceType+".test_maximal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_maximal").Key("subject").HasValue("repo:octo-org/octo-repo:environment:Production"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("This is a test federated identity credential with all optional fields configured"),
					check.That(resourceType+".test_maximal").Key("audiences.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_maximal"),
			},
		},
	})
}

func TestUnitResourceApplicationFederatedIdentityCredential_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationFederatedIdentityCredentialMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationFederatedIdentityCredentialMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
				),
			},
			{
				Config: testConfigMinimalUpdated(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("issuer").HasValue("https://token.actions.githubusercontent.com"),
					check.That(resourceType+".test_minimal").Key("description").HasValue("Updated description for unit test"),
				),
			},
		},
	})
}

func testConfigMinimalUpdated() string {
	return `
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "test_minimal" {
  application_id = "11111111-1111-1111-1111-111111111111"
  name           = "unit-test-fic-minimal"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:octo-org/octo-repo:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]
  description    = "Updated description for unit test"
}
`
}

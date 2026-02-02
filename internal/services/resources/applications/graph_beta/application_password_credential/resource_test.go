package graphBetaApplicationPasswordCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationPasswordCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_password_credential"
	passwordCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_password_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationPasswordCredential.ResourceName

	// testResource is the test resource implementation for password credentials
	testResource = graphBetaApplicationPasswordCredential.ApplicationPasswordCredentialTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *passwordCredentialMocks.ApplicationPasswordCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	passwordCredMock := &passwordCredentialMocks.ApplicationPasswordCredentialMock{}
	passwordCredMock.RegisterMocks()

	return mockClient, passwordCredMock
}

func TestUnitResourceApplicationPasswordCredential_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, passwordCredMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer passwordCredMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("secret_text").MatchesRegex(regexp.MustCompile(`^generatedSecretText~`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("unit-test-password-credential"),
					check.That(resourceType+".test_minimal").Key("hint").MatchesRegex(regexp.MustCompile(`^gen`)),
				),
			},
		},
	})
}

func TestUnitResourceApplicationPasswordCredential_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, passwordCredMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer passwordCredMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("application_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_maximal").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("secret_text").MatchesRegex(regexp.MustCompile(`^generatedSecretText~`)),
					check.That(resourceType+".test_maximal").Key("display_name").HasValue("unit-test-password-credential-maximal"),
					check.That(resourceType+".test_maximal").Key("hint").MatchesRegex(regexp.MustCompile(`^gen`)),
					check.That(resourceType+".test_maximal").Key("start_date_time").HasValue("2027-01-01T00:00:00Z"),
					check.That(resourceType+".test_maximal").Key("end_date_time").HasValue("2029-01-01T00:00:00Z"),
				),
			},
		},
	})
}

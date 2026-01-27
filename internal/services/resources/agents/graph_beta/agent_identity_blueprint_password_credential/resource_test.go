package graphBetaAgentIdentityBlueprintPasswordCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	agentIdentityBlueprintMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint/mocks"
	graphBetaAgentIdentityBlueprintPasswordCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_password_credential"
	passwordCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_password_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentityBlueprintPasswordCredential.ResourceName

	// testResource is the test resource implementation for password credentials
	testResource = graphBetaAgentIdentityBlueprintPasswordCredential.AgentIdentityBlueprintPasswordCredentialTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintMocks.AgentIdentityBlueprintMock, *passwordCredentialMocks.AgentIdentityBlueprintPasswordCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register both the blueprint and password credential mocks
	blueprintMock := &agentIdentityBlueprintMocks.AgentIdentityBlueprintMock{}
	blueprintMock.RegisterMocks()

	passwordCredMock := &passwordCredentialMocks.AgentIdentityBlueprintPasswordCredentialMock{}
	passwordCredMock.RegisterMocks()

	return mockClient, blueprintMock, passwordCredMock
}

func TestUnitResourceAgentIdentityBlueprintPasswordCredential_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, blueprintMock, passwordCredMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer blueprintMock.CleanupMockState()
	defer passwordCredMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("secret_text").MatchesRegex(regexp.MustCompile(`^generatedSecretText~`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("unit-test-password-credential"),
					check.That(resourceType+".test_minimal").Key("hint").MatchesRegex(regexp.MustCompile(`^gen`)),
				),
			},
		},
	})
}

func testConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

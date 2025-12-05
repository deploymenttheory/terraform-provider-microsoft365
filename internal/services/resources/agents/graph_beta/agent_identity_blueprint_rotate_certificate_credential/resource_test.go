package graphBetaAgentIdentityBlueprintKeyCredential_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	agentIdentityBlueprintMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint/mocks"
	graphBetaAgentIdentityBlueprintKeyCredential "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_key_credential"
	keyCredentialMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_identity_blueprint_key_credential/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAgentIdentityBlueprintKeyCredential.ResourceName

	// testResource is the test resource implementation for key credentials
	testResource = graphBetaAgentIdentityBlueprintKeyCredential.AgentIdentityBlueprintKeyCredentialTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentIdentityBlueprintMocks.AgentIdentityBlueprintMock, *keyCredentialMocks.AgentIdentityBlueprintKeyCredentialMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register both the blueprint and key credential mocks
	blueprintMock := &agentIdentityBlueprintMocks.AgentIdentityBlueprintMock{}
	blueprintMock.RegisterMocks()

	keyCredMock := &keyCredentialMocks.AgentIdentityBlueprintKeyCredentialMock{}
	keyCredMock.RegisterMocks()

	return mockClient, blueprintMock, keyCredMock
}

func TestAgentIdentityBlueprintKeyCredentialResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, blueprintMock, keyCredMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer blueprintMock.CleanupMockState()
	defer keyCredMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("key_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("unit-test-key-credential"),
					check.That(resourceType+".test_minimal").Key("type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_minimal").Key("usage").HasValue("Verify"),
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

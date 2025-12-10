package graphBetaAgentUser_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAgentUser "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_user"
	agentUserMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/agents/graph_beta/agent_user/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaAgentUser.ResourceName
	testResource = graphBetaAgentUser.AgentUserTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *agentUserMocks.AgentUserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentUserMock := &agentUserMocks.AgentUserMock{}
	agentUserMock.RegisterMocks()
	return mockClient, agentUserMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *agentUserMocks.AgentUserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	agentUserMock := &agentUserMocks.AgentUserMock{}
	agentUserMock.RegisterErrorMocks()
	return mockClient, agentUserMock
}

func TestUnitAgentUserResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentUserMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentUserMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("Unit Test Agent User"),
					check.That(resourceType+".test_minimal").Key("agent_identity_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("user_principal_name").HasValue("testagentuser@contoso.onmicrosoft.com"),
					check.That(resourceType+".test_minimal").Key("mail_nickname").HasValue("testagentuser"),
					check.That(resourceType+".test_minimal").Key("sponsor_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitAgentUserResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, agentUserMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer agentUserMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").HasValue("Unit Test Agent User Maximal"),
					check.That(resourceType+".test_maximal").Key("agent_identity_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_maximal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("user_principal_name").HasValue("testagentusermaximal@contoso.onmicrosoft.com"),
					check.That(resourceType+".test_maximal").Key("mail_nickname").HasValue("testagentusermaximal"),
					check.That(resourceType+".test_maximal").Key("sponsor_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("given_name").HasValue("Test"),
					check.That(resourceType+".test_maximal").Key("surname").HasValue("AgentUser"),
					check.That(resourceType+".test_maximal").Key("job_title").HasValue("AI Agent"),
					check.That(resourceType+".test_maximal").Key("department").HasValue("Engineering"),
					check.That(resourceType+".test_maximal").Key("company_name").HasValue("Contoso"),
					check.That(resourceType+".test_maximal").Key("office_location").HasValue("Building A"),
					check.That(resourceType+".test_maximal").Key("city").HasValue("Seattle"),
					check.That(resourceType+".test_maximal").Key("state").HasValue("WA"),
					check.That(resourceType+".test_maximal").Key("country").HasValue("US"),
					check.That(resourceType+".test_maximal").Key("postal_code").HasValue("98101"),
					check.That(resourceType+".test_maximal").Key("street_address").HasValue("123 Main Street"),
					check.That(resourceType+".test_maximal").Key("usage_location").HasValue("US"),
					check.That(resourceType+".test_maximal").Key("preferred_language").HasValue("en-US"),
				),
			},
			{
				ResourceName: resourceType + ".test_maximal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_maximal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_maximal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
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

func testConfigMaximal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

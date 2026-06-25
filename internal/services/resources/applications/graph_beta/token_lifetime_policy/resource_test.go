package graphBetaApplicationsTokenLifetimePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	tokenLifetimePolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/token_lifetime_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = "microsoft365_graph_beta_applications_token_lifetime_policy"

func setupMockEnvironment() (*mocks.Mocks, *tokenLifetimePolicyMocks.TokenLifetimePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tlpMock := &tokenLifetimePolicyMocks.TokenLifetimePolicyMock{}
	tlpMock.RegisterMocks()
	return mockClient, tlpMock
}

func TestUnitResourceTokenLifetimePolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tlpMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tlpMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".minimal", "display_name", "unit-test-token-lifetime-policy-min"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "description", "Unit test minimal token lifetime policy"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "is_organization_default", "false"),
					resource.TestCheckResourceAttr(resourceType+".minimal", "definition.#", "1"),
					resource.TestMatchResourceAttr(resourceType+".minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceTokenLifetimePolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tlpMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tlpMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".maximal", "display_name", "unit-test-token-lifetime-policy-max"),
					resource.TestCheckResourceAttr(resourceType+".maximal", "description", "Unit test maximal token lifetime policy"),
					resource.TestCheckResourceAttr(resourceType+".maximal", "is_organization_default", "true"),
					resource.TestCheckResourceAttr(resourceType+".maximal", "definition.#", "1"),
					resource.TestMatchResourceAttr(resourceType+".maximal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load token lifetime policy minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load token lifetime policy maximal config: " + err.Error())
	}
	return unitTestConfig
}

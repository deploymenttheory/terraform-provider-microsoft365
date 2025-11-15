package graphBetaFilteringPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	filteringPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/filtering_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *filteringPolicyMocks.FilteringPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringPolicyMock := &filteringPolicyMocks.FilteringPolicyMock{}
	filteringPolicyMock.RegisterMocks()
	return mockClient, filteringPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *filteringPolicyMocks.FilteringPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringPolicyMock := &filteringPolicyMocks.FilteringPolicyMock{}
	filteringPolicyMock.RegisterErrorMocks()
	return mockClient, filteringPolicyMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestFilteringPolicyResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "name", "Test Filtering Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "description", "Test filtering policy for unit testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "action", "block"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_filtering_policy.test"),
				),
			},
		},
	})
}

func TestFilteringPolicyResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "name", "Test Filtering Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "action", "block"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "name", "Updated Filtering Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_filtering_policy.test", "action", "allow"),
				),
			},
		},
	})
}

func TestFilteringPolicyResource_InvalidAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidAction(),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestFilteringPolicyResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Invalid filtering policy data`),
			},
		},
	})
}

func testConfigBasic() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_filtering_policy" "test" {
  name        = "Test Filtering Policy"
  description = "Test filtering policy for unit testing"
  action      = "block"
}
`
}

func testConfigUpdate() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_filtering_policy" "test" {
  name        = "Updated Filtering Policy"
  description = "Updated description"
  action      = "allow"
}
`
}

func testConfigInvalidAction() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_filtering_policy" "test" {
  name        = "Test Filtering Policy"
  description = "Test filtering policy with invalid action"
  action      = "invalid_action"
}
`
}

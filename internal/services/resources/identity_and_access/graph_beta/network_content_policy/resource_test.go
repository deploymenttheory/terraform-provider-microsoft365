package graphBetaNetworkContentPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	contentPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_content_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_network_content_policy"

func setupMockEnvironment() (*mocks.Mocks, *contentPolicyMocks.ContentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	contentPolicyMock := &contentPolicyMocks.ContentPolicyMock{}
	contentPolicyMock.RegisterMocks()
	return mockClient, contentPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *contentPolicyMocks.ContentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	contentPolicyMock := &contentPolicyMocks.ContentPolicyMock{}
	contentPolicyMock.RegisterErrorMocks()
	return mockClient, contentPolicyMock
}

func TestUnitResourceNetworkContentPolicy_01_BasicAndImport(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-content-policy-minimal"),
					check.That(resourceType+".test").Key("description").HasValue(""),
					check.That(resourceType+".test").Key("default_action").HasValue("allow"),
					check.That(resourceType+".test").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceNetworkContentPolicy_02_UpdateAndClearDescription(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_with_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("description").HasValue("managed by Terraform"),
				),
			},
			{
				Config: testConfigHelper("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-content-policy-updated"),
					check.That(resourceType+".test").Key("description").HasValue(""),
				),
			},
		},
	})
}

func TestUnitResourceNetworkContentPolicy_03_InvalidDefaultAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      testConfigHelper("resource_invalid_default_action.tf"),
			ExpectError: regexp.MustCompile(`default_action.*allow|value must be one of`),
		}},
	})
}

func TestUnitResourceNetworkContentPolicy_04_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, contentPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer contentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      testConfigHelper("resource_minimal.tf"),
			ExpectError: regexp.MustCompile(`Bad Request - 400`),
		}},
	})
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

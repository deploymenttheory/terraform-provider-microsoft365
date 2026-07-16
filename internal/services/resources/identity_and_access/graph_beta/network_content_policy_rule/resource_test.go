package graphBetaNetworkContentPolicyRule_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	ruleMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_content_policy_rule/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_network_content_policy_rule"

func setupMockEnvironment(errorMocks bool) *ruleMocks.ContentPolicyRuleMock {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	ruleMock := &ruleMocks.ContentPolicyRuleMock{}
	if errorMocks {
		ruleMock.RegisterErrorMocks()
	} else {
		ruleMock.RegisterMocks()
	}
	return ruleMock
}

func TestUnitResourceNetworkContentPolicyRule_01_CRUDAndImport(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	ruleMock := setupMockEnvironment(false)
	defer httpmock.DeactivateAndReset()
	defer ruleMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("name").HasValue("unit-test-content-policy-rule"),
					check.That(resourceType+".test").Key("action").HasValue("scanPurview"),
					check.That(resourceType+".test").Key("status").HasValue("enabled"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "00000000-0000-0000-0000-000000000301/00000000-0000-0000-0000-000000000302",
				ImportStateVerify: true,
			},
			{
				Config: testConfig("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-content-policy-rule-updated"),
					check.That(resourceType+".test").Key("action").HasValue("block"),
					check.That(resourceType+".test").Key("priority").HasValue("102"),
					check.That(resourceType+".test").Key("status").HasValue("disabled"),
				),
			},
		},
	})
}

func TestUnitResourceNetworkContentPolicyRule_02_InvalidAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	ruleMock := setupMockEnvironment(false)
	defer httpmock.DeactivateAndReset()
	defer ruleMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: testConfig("resource_invalid_action.tf"), ExpectError: regexp.MustCompile(`action.*allow.*block.*scanPurview|value must be one of`)}},
	})
}

func TestUnitResourceNetworkContentPolicyRule_03_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	ruleMock := setupMockEnvironment(true)
	defer httpmock.DeactivateAndReset()
	defer ruleMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: testConfig("resource_minimal.tf"), ExpectError: regexp.MustCompile(`Bad Request - 400`)}},
	})
}

func testConfig(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

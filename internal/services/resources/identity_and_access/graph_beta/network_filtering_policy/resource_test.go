package graphBetaNetworkFilteringPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	networkFilteringPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// resourceType is declared in resource_acceptance_test.go and shared across the package

func setupMockEnvironment() (*mocks.Mocks, *networkFilteringPolicyMocks.FilteringPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringPolicyMock := &networkFilteringPolicyMocks.FilteringPolicyMock{}
	filteringPolicyMock.RegisterMocks()
	return mockClient, filteringPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *networkFilteringPolicyMocks.FilteringPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringPolicyMock := &networkFilteringPolicyMocks.FilteringPolicyMock{}
	filteringPolicyMock.RegisterErrorMocks()
	return mockClient, filteringPolicyMock
}

func TestNetworkFilteringPolicyResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-policy-minimal"),
					check.That(resourceType+".test").Key("description").HasValue("Test filtering policy for unit testing"),
					check.That(resourceType+".test").Key("action").HasValue("block"),
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

func TestNetworkFilteringPolicyResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-policy-minimal"),
					check.That(resourceType+".test").Key("action").HasValue("block"),
				),
			},
			{
				Config: testConfigHelper("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-policy-updated"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description"),
					check.That(resourceType+".test").Key("action").HasValue("allow"),
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

func TestNetworkFilteringPolicyResource_InvalidAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_invalid.tf"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestNetworkFilteringPolicyResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Invalid filtering policy data`),
			},
		},
	})
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

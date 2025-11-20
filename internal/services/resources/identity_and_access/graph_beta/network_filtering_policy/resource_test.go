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
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("Test Filtering Policy"),
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
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("Test Filtering Policy"),
					check.That(resourceType+".test").Key("action").HasValue("block"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("Updated Filtering Policy"),
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
				Config:      testConfigInvalidAction(),
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
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Invalid filtering policy data`),
			},
		},
	})
}

func testConfigBasic() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load resource_minimal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUpdate() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_updated.tf")
	if err != nil {
		panic("failed to load resource_updated.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidAction() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_invalid.tf")
	if err != nil {
		panic("failed to load resource_invalid.tf: " + err.Error())
	}
	return unitTestConfig
}

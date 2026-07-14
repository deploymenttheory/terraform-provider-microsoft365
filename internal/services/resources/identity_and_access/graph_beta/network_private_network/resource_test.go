package graphBetaNetworkPrivateNetwork_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	privateNetworkMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_private_network/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_network_private_network"

func setupMockEnvironment() (*mocks.Mocks, *privateNetworkMocks.PrivateNetworkMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	privateNetworkMock := &privateNetworkMocks.PrivateNetworkMock{}
	privateNetworkMock.RegisterMocks()
	return mockClient, privateNetworkMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *privateNetworkMocks.PrivateNetworkMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	privateNetworkMock := &privateNetworkMocks.PrivateNetworkMock{}
	privateNetworkMock.RegisterErrorMocks()
	return mockClient, privateNetworkMock
}

func TestUnitResourceNetworkPrivateNetwork_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, privateNetworkMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer privateNetworkMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-private-network-minimal"),
					check.That(resourceType+".test").Key("dns_resolution_identification.fqdn_to_resolve").HasValue("example.com"),
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

func TestUnitResourceNetworkPrivateNetwork_02_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, privateNetworkMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer privateNetworkMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("name").HasValue("unit-test-private-network-minimal"),
				),
			},
			{
				Config: testConfigHelper("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-private-network-updated"),
					check.That(resourceType+".test").Key("app_ids.#").HasValue("1"),
					check.That(resourceType+".test").Key("dns_resolution_identification.expected_ip_resolutions.#").HasValue("3"),
				),
			},
		},
	})
}

func TestUnitResourceNetworkPrivateNetwork_03_InvalidExpectedIPResolution(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, privateNetworkMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer privateNetworkMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_invalid.tf"),
				ExpectError: regexp.MustCompile(`type "ip_address" requires "value"`),
			},
			{
				Config:      testConfigHelper("resource_invalid_range.tf"),
				ExpectError: regexp.MustCompile(`type "ip_range" requires "begin_address"`),
			},
		},
	})
}

func TestUnitResourceNetworkPrivateNetwork_04_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, privateNetworkMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer privateNetworkMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Bad Request - 400`),
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

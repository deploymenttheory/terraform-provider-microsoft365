package graphBetaNetworkFilteringProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	networkFilteringProfileMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_filtering_profile/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// resourceType is declared in resource_acceptance_test.go and shared across the package

func setupMockEnvironment() (*mocks.Mocks, *networkFilteringProfileMocks.FilteringProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringProfileMock := &networkFilteringProfileMocks.FilteringProfileMock{}
	filteringProfileMock.RegisterMocks()
	return mockClient, filteringProfileMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *networkFilteringProfileMocks.FilteringProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	filteringProfileMock := &networkFilteringProfileMocks.FilteringProfileMock{}
	filteringProfileMock.RegisterErrorMocks()
	return mockClient, filteringProfileMock
}

func TestUnitResourceNetworkFilteringProfile_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringProfileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringProfileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-profile-minimal"),
					check.That(resourceType+".test").Key("description").HasValue("Test filtering profile for unit testing"),
					check.That(resourceType+".test").Key("priority").HasValue("100"),
					check.That(resourceType+".test").Key("state").HasValue("enabled"),
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

func TestUnitResourceNetworkFilteringProfile_02_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringProfileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringProfileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-profile-minimal"),
					check.That(resourceType+".test").Key("state").HasValue("enabled"),
				),
			},
			{
				Config: testConfigHelper("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("name").HasValue("unit-test-filtering-profile-updated"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description"),
					check.That(resourceType+".test").Key("priority").HasValue("200"),
					check.That(resourceType+".test").Key("state").HasValue("disabled"),
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

func TestUnitResourceNetworkFilteringProfile_03_InvalidState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringProfileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringProfileMock.CleanupMockState()

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

func TestUnitResourceNetworkFilteringProfile_04_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, filteringProfileMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer filteringProfileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Invalid filtering profile data`),
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

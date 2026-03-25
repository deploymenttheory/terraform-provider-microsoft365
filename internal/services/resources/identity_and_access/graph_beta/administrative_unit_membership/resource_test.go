package graphBetaAdministrativeUnitMembership_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	membershipMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit_membership/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *membershipMocks.AdministrativeUnitMembershipMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	membershipMock := &membershipMocks.AdministrativeUnitMembershipMock{}
	membershipMock.RegisterMocks()
	return mockClient, membershipMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *membershipMocks.AdministrativeUnitMembershipMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	membershipMock := &membershipMocks.AdministrativeUnitMembershipMock{}
	membershipMock.RegisterErrorMocks()
	return mockClient, membershipMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Basic membership with two users
func TestUnitResourceAdministrativeUnitMembership_01_AUM001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_aum001_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".aum001_basic").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aum001_basic").Key("administrative_unit_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".aum001_basic").Key("members.#").HasValue("2"),
				),
			},
			{
				ResourceName:      resourceType + ".aum001_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Single member
func TestUnitResourceAdministrativeUnitMembership_02_AUM002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_aum002_single.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".aum002_single").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aum002_single").Key("administrative_unit_id").HasValue("22222222-1111-1111-1111-111111111111"),
					check.That(resourceType+".aum002_single").Key("members.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".aum002_single",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

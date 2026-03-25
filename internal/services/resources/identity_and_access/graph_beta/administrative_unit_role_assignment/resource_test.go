package graphBetaAdministrativeUnitRoleAssignment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAdministrativeUnitRoleAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit_role_assignment"
	roleAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit_role_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaAdministrativeUnitRoleAssignment.ResourceName
	testResource = graphBetaAdministrativeUnitRoleAssignment.AdministrativeUnitRoleAssignmentTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *roleAssignmentMocks.AdministrativeUnitRoleAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	roleAssignmentMock := &roleAssignmentMocks.AdministrativeUnitRoleAssignmentMock{}
	roleAssignmentMock.RegisterMocks()
	return mockClient, roleAssignmentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *roleAssignmentMocks.AdministrativeUnitRoleAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	roleAssignmentMock := &roleAssignmentMocks.AdministrativeUnitRoleAssignmentMock{}
	roleAssignmentMock.RegisterErrorMocks()
	return mockClient, roleAssignmentMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Basic scoped role assignment with a single user
func TestUnitResourceAdministrativeUnitRoleAssignment_01_AURA001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_aura001_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".aura001_basic").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aura001_basic").Key("administrative_unit_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".aura001_basic").Key("role_id").HasValue("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					check.That(resourceType+".aura001_basic").Key("role_member_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".aura001_basic").Key("role_member_display_name").HasValue("Test User 1"),
				),
			},
			{
				ResourceName:      resourceType + ".aura001_basic",
				ImportState:       true,
				ImportStateIdFunc: importStateIDFunc(resourceType + ".aura001_basic"),
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Scoped role assignment with a different user
func TestUnitResourceAdministrativeUnitRoleAssignment_02_AURA002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_aura002_different_member.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".aura002_different_member").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aura002_different_member").Key("administrative_unit_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".aura002_different_member").Key("role_id").HasValue("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
					check.That(resourceType+".aura002_different_member").Key("role_member_id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".aura002_different_member").Key("role_member_display_name").HasValue("Test User 2"),
				),
			},
		},
	})
}

// importStateIDFunc returns a function that constructs the composite import ID from Terraform state
func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", nil
		}
		auID := rs.Primary.Attributes["administrative_unit_id"]
		id := rs.Primary.Attributes["id"]
		return auID + "/" + id, nil
	}
}

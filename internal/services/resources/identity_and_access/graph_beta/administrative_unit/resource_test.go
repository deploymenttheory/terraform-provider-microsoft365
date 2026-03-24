package graphBetaAdministrativeUnit_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	administrativeUnitMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *administrativeUnitMocks.AdministrativeUnitMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	administrativeUnitMock := &administrativeUnitMocks.AdministrativeUnitMock{}
	administrativeUnitMock.RegisterMocks()
	return mockClient, administrativeUnitMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *administrativeUnitMocks.AdministrativeUnitMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	administrativeUnitMock := &administrativeUnitMocks.AdministrativeUnitMock{}
	administrativeUnitMock.RegisterErrorMocks()
	return mockClient, administrativeUnitMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// AU001: User-Based Administrative Unit
func TestUnitResourceAdministrativeUnit_01_AU001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, administrativeUnitMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer administrativeUnitMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_au001_user_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".au001_user_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au001_user_based").Key("display_name").HasValue("AU001: User-Based Administrative Unit"),
					check.That(resourceType+".au001_user_based").Key("description").HasValue("Administrative unit for user-based testing"),
					check.That(resourceType+".au001_user_based").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				ResourceName: resourceType + ".au001_user_based",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au001_user_based"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au001_user_based")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU002: Group-Based Administrative Unit
func TestUnitResourceAdministrativeUnit_02_AU002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, administrativeUnitMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer administrativeUnitMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_au002_group_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".au002_group_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au002_group_based").Key("display_name").HasValue("AU002: Group-Based Administrative Unit"),
					check.That(resourceType+".au002_group_based").Key("description").HasValue("Administrative unit for group-based testing"),
					check.That(resourceType+".au002_group_based").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				ResourceName: resourceType + ".au002_group_based",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au002_group_based"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au002_group_based")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU003: Mixed User and Group Administrative Unit
func TestUnitResourceAdministrativeUnit_03_AU003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, administrativeUnitMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer administrativeUnitMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_au003_mixed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".au003_mixed").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au003_mixed").Key("display_name").HasValue("AU003: Mixed User and Group Administrative Unit"),
					check.That(resourceType+".au003_mixed").Key("description").HasValue("Administrative unit for mixed user and group testing"),
					check.That(resourceType+".au003_mixed").Key("visibility").HasValue("HiddenMembership"),
					check.That(resourceType+".au003_mixed").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				ResourceName: resourceType + ".au003_mixed",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au003_mixed"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au003_mixed")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU004: Dynamic Administrative Unit
func TestUnitResourceAdministrativeUnit_04_AU004(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, administrativeUnitMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer administrativeUnitMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_au004_dynamic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".au004_dynamic").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au004_dynamic").Key("display_name").HasValue("AU004: Dynamic Administrative Unit"),
					check.That(resourceType+".au004_dynamic").Key("description").HasValue("Administrative unit with dynamic membership"),
					check.That(resourceType+".au004_dynamic").Key("membership_type").HasValue("Dynamic"),
					check.That(resourceType+".au004_dynamic").Key("membership_rule").HasValue("(user.country -eq \"United States\")"),
					check.That(resourceType+".au004_dynamic").Key("membership_rule_processing_state").HasValue("On"),
					check.That(resourceType+".au004_dynamic").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				ResourceName: resourceType + ".au004_dynamic",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au004_dynamic"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au004_dynamic")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

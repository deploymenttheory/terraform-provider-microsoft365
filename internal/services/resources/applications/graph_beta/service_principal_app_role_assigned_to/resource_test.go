package graphBetaServicePrincipalAppRoleAssignedTo_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalAppRoleAssignedTo "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_app_role_assigned_to"
	appRoleAssignedToMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_app_role_assigned_to/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalAppRoleAssignedTo.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *appRoleAssignedToMocks.ServicePrincipalAppRoleAssignedToMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appRoleAssignedToMock := &appRoleAssignedToMocks.ServicePrincipalAppRoleAssignedToMock{}
	appRoleAssignedToMock.RegisterMocks()
	return mockClient, appRoleAssignedToMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *appRoleAssignedToMocks.ServicePrincipalAppRoleAssignedToMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appRoleAssignedToMock := &appRoleAssignedToMocks.ServicePrincipalAppRoleAssignedToMock{}
	appRoleAssignedToMock.RegisterErrorMocks()
	return mockClient, appRoleAssignedToMock
}

// TestUnitServicePrincipalAppRoleAssignedToResource_Minimal tests creating an app role assignment with minimal configuration
func TestUnitServicePrincipalAppRoleAssignedToResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appRoleAssignedToMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appRoleAssignedToMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("resource_object_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_minimal").Key("app_role_id").HasValue("df021288-bdef-4463-88db-98f22de89214"),
					check.That(resourceType+".test_minimal").Key("target_service_principal_object_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("principal_type").Exists(),
					check.That(resourceType+".test_minimal").Key("principal_display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("resource_display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("created_date_time").Exists(),
				),
			},
			{
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					resourceObjectID := rs.Primary.Attributes["resource_object_id"]
					id := rs.Primary.Attributes["id"]
					return fmt.Sprintf("%s/%s", resourceObjectID, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitServicePrincipalAppRoleAssignedToResource_DefaultRole tests creating an app role assignment with default role ID
func TestUnitServicePrincipalAppRoleAssignedToResource_DefaultRole(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appRoleAssignedToMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appRoleAssignedToMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigDefaultRole(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_default_role").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_default_role").Key("app_role_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".test_default_role").Key("resource_object_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_default_role").Key("target_service_principal_object_id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
			{
				ResourceName: resourceType + ".test_default_role",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_default_role"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_default_role")
					}
					resourceObjectID := rs.Primary.Attributes["resource_object_id"]
					id := rs.Primary.Attributes["id"]
					return fmt.Sprintf("%s/%s", resourceObjectID, id), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration helper functions using helpers.ParseHCLFile
func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load service principal app role assigned to minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigDefaultRole() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_default_role.tf")
	if err != nil {
		panic("failed to load service principal app role assigned to default role config: " + err.Error())
	}
	return unitTestConfig
}

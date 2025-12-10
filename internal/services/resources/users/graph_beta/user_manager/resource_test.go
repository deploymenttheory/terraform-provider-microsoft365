package graphBetaUsersUserManager_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUserManager "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user_manager"
	userManagerMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user_manager/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaUserManager.ResourceName
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *userManagerMocks.UserManagerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	userManagerMock := &userManagerMocks.UserManagerMock{}
	userManagerMock.RegisterMocks()
	return mockClient, userManagerMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *userManagerMocks.UserManagerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	userManagerMock := &userManagerMocks.UserManagerMock{}
	userManagerMock.RegisterErrorMocks()
	return mockClient, userManagerMock
}

// TestUserManagerResource_Lifecycle tests the full lifecycle of user manager relationship
func TestUserManagerResource_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userManagerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userManagerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceType+".test", "user_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr(resourceType+".test", "manager_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
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

// TestUserManagerResource_RequiredFields tests required field validation
func TestUserManagerResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userManagerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userManagerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_users_user_manager" "test" {
  # Missing user_id
  manager_id = "00000000-0000-0000-0000-000000000002"
}
`,
				ExpectError: regexp.MustCompile(`The argument "user_id" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_users_user_manager" "test" {
  user_id = "00000000-0000-0000-0000-000000000001"
  # Missing manager_id
}
`,
				ExpectError: regexp.MustCompile(`The argument "manager_id" is required`),
			},
		},
	})
}

// TestUserManagerResource_InvalidGUID tests GUID validation
func TestUserManagerResource_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userManagerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userManagerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_users_user_manager" "test" {
  user_id    = "invalid-guid"
  manager_id = "00000000-0000-0000-0000-000000000002"
}
`,
				ExpectError: regexp.MustCompile(`must be a valid GUID`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_users_user_manager" "test" {
  user_id    = "00000000-0000-0000-0000-000000000001"
  manager_id = "invalid-guid"
}
`,
				ExpectError: regexp.MustCompile(`must be a valid GUID`),
			},
		},
	})
}

// TestUserManagerResource_ErrorHandling tests error scenarios
func TestUserManagerResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userManagerMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userManagerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile(`Bad Request|BadRequest`),
			},
		},
	})
}

// Config loader functions
func testConfigMinimal() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return config
}

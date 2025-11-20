package graphBetaUsersUser_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	userMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *userMocks.UserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	userMock := &userMocks.UserMock{}
	userMock.RegisterMocks()
	return mockClient, userMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *userMocks.UserMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	userMock := &userMocks.UserMock{}
	userMock.RegisterErrorMocks()
	return mockClient, userMock
}

func TestUserResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("Minimal User"),
					check.That(resourceType+".minimal").Key("user_principal_name").HasValue("minimal.user@deploymenttheory.com"),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
				},
			},
		},
	})
}

func TestUserResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("Minimal User"),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("display_name").HasValue("Maximal User"),
					check.That(resourceType+".maximal").Key("given_name").HasValue("Maximal"),
					check.That(resourceType+".maximal").Key("surname").HasValue("User"),
					check.That(resourceType+".maximal").Key("job_title").HasValue("Senior Developer"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
				},
			},
		},
	})
}

func TestUserResource_RequiredAttributes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMissingRequired(),
				ExpectError: regexp.MustCompile(`Missing required argument|The argument .* is required`),
			},
		},
	})
}

func TestUserResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Bad Request - 400|The request was invalid or malformed|ApiError`),
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
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load resource_maximal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMissingRequired() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_invalid.tf")
	if err != nil {
		panic("failed to load resource_invalid.tf: " + err.Error())
	}
	return unitTestConfig
}

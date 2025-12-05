package graphBetaUsersAgentUser_test

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
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-agent-user-minimal"),
					check.That(resourceType+".minimal").Key("user_principal_name").HasValue("unit-test-agent-user-minimal@deploymenttheory.com"),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("identity_parent_id").HasValue("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					check.That(resourceType+".minimal").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
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
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-agent-user-minimal"),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-agent-user-maximal"),
					check.That(resourceType+".maximal").Key("user_principal_name").HasValue("unit-test-agent-user-maximal@deploymenttheory.com"),
					check.That(resourceType+".maximal").Key("given_name").HasValue("Maximal"),
					check.That(resourceType+".maximal").Key("surname").HasValue("User"),
					check.That(resourceType+".maximal").Key("job_title").HasValue("Marketing Agent"),
					check.That(resourceType+".maximal").Key("department").HasValue("Marketing"),
					check.That(resourceType+".maximal").Key("company_name").HasValue("Deployment Theory"),
					check.That(resourceType+".maximal").Key("employee_id").HasValue("1234567890"),
					check.That(resourceType+".maximal").Key("employee_type").HasValue("full time"),
					check.That(resourceType+".maximal").Key("age_group").HasValue("NotAdult"),
					check.That(resourceType+".maximal").Key("consent_provided_for_minor").HasValue("Granted"),
					check.That(resourceType+".maximal").Key("identity_parent_id").HasValue("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUserResource_CustomSecurityAttributes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, userMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer userMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCustomSecAtt(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_custom_security_attributes").Key("display_name").HasValue("unit-test-agent-user-custom-sec-att"),
					check.That(resourceType+".with_custom_security_attributes").Key("user_principal_name").HasValue("unit-test-agent-user-custom-sec-att@deploymenttheory.com"),
					check.That(resourceType+".with_custom_security_attributes").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".with_custom_security_attributes").Key("identity_parent_id").HasValue("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.#").HasValue("2"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.0.attribute_set").HasValue("Engineering"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.1.attribute_set").HasValue("Marketing"),
				),
			},
			{
				ResourceName:      resourceType + ".with_custom_security_attributes",
				ImportState:       true,
				ImportStateVerify: true,
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

func testConfigCustomSecAtt() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_custom_sec_att.tf")
	if err != nil {
		panic("failed to load resource_custom_sec_att.tf: " + err.Error())
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

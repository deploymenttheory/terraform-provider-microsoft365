package graphBetaCustomSecurityAttributeAllowedValue_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	allowedValueMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_allowed_value/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *allowedValueMocks.CustomSecurityAttributeAllowedValueMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	allowedValueMock := &allowedValueMocks.CustomSecurityAttributeAllowedValueMock{}
	allowedValueMock.RegisterMocks()
	return mockClient, allowedValueMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *allowedValueMocks.CustomSecurityAttributeAllowedValueMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	allowedValueMock := &allowedValueMocks.CustomSecurityAttributeAllowedValueMock{}
	allowedValueMock.RegisterErrorMocks()
	return mockClient, allowedValueMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAllowedValueResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "custom_security_attribute_definition_id", "Engineering_Project"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Alpine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "true"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test"),
				),
			},
		},
	})
}

func TestAllowedValueResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "custom_security_attribute_definition_id", "Engineering_Project"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "true"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test"),
				),
			},
		},
	})
}

func TestAllowedValueResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Alpine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "true"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Alpine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "false"),
				),
			},
		},
	})
}

func TestAllowedValueResource_Inactive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigInactive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "custom_security_attribute_definition_id", "Engineering_Project"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Legacy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "false"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test"),
				),
			},
		},
	})
}

func TestAllowedValueResource_WithSpaces(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithSpaces(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "custom_security_attribute_definition_id", "Engineering_Department"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "id", "Human Resources"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test", "is_active", "true"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test"),
				),
			},
		},
	})
}

func TestAllowedValueResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`(?i)(bad.*request|malformed|incorrect)`),
			},
		},
	})
}

func TestAllowedValueResource_NotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`(?i)(bad.*request|malformed|incorrect)`),
			},
		},
	})
}

func TestAllowedValueResource_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, allowedValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer allowedValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
			},
			{
				ResourceName:      "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.test",
				ImportState:       true,
				ImportStateId:     "Engineering_Project/Alpine",
				ImportStateVerify: true,
			},
		},
	})
}

// Config helper functions
func testConfigBasic() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Project"
  id                                       = "Alpine"
  is_active                                = true
}
`
}

func testConfigMinimal() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Project"
  id                                       = "Minimal"
  is_active                                = true
}
`
}

func testConfigUpdate() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Project"
  id                                       = "Alpine"
  is_active                                = false
}
`
}

func testConfigInactive() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Project"
  id                                       = "Legacy"
  is_active                                = false
}
`
}

func testConfigWithSpaces() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Department"
  id                                       = "Human Resources"
  is_active                                = true
}
`
}

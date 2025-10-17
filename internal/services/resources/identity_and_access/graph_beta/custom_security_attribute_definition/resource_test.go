package graphBetaCustomSecurityAttributeDefinition_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	definitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_definition/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *definitionMocks.CustomSecurityAttributeDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	definitionMock := &definitionMocks.CustomSecurityAttributeDefinitionMock{}
	definitionMock.RegisterMocks()
	return mockClient, definitionMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *definitionMocks.CustomSecurityAttributeDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	definitionMock := &definitionMocks.CustomSecurityAttributeDefinitionMock{}
	definitionMock.RegisterErrorMocks()
	return mockClient, definitionMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestCustomSecurityAttributeDefinitionResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "attribute_set", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "name", "ProjectName"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "description", "Name of the project the user is assigned to"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "type", "String"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "status", "Available"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "is_collection", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "is_searchable", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "use_pre_defined_values_only", "false"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test"),
				),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "attribute_set", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "name", "MinimalAttribute"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "type", "String"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "description"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test"),
				),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "description", "Name of the project the user is assigned to"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "status", "Available"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "description", "Updated project name description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "status", "Available"),
				),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_BooleanType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBoolean(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "type", "Boolean"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "is_collection", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "use_pre_defined_values_only", "false"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test"),
				),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_Collection(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCollection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "is_collection", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test", "type", "String"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.test"),
				),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_InvalidAttributeSetName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidAttributeSet(),
				ExpectError: regexp.MustCompile(`(?i)attribute.set.*cannot contain spaces`),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_InvalidName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidName(),
				ExpectError: regexp.MustCompile(`Name cannot contain spaces or special characters`),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_InvalidBooleanWithCollection(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidBooleanWithCollection(),
				ExpectError: regexp.MustCompile(`When type is set to Boolean, isCollection must be set to false`),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_InvalidBooleanWithPredefinedValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidBooleanWithPredefinedValues(),
				ExpectError: regexp.MustCompile(`When type is set to Boolean, usePreDefinedValuesOnly must be set to false`),
			},
		},
	})
}

func TestCustomSecurityAttributeDefinitionResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, definitionMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer definitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Invalid custom security attribute definition data`),
			},
		},
	})
}

func testConfigBasic() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_basic.tf")
	if err != nil {
		panic("failed to load test_config_basic config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_minimal.tf")
	if err != nil {
		panic("failed to load test_config_minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUpdate() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_update.tf")
	if err != nil {
		panic("failed to load test_config_update config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBoolean() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_boolean.tf")
	if err != nil {
		panic("failed to load test_config_boolean config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCollection() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_collection.tf")
	if err != nil {
		panic("failed to load test_config_collection config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInvalidAttributeSet() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Invalid Set With Spaces"
  name                        = "TestAttribute"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false
}
`
}

func testConfigInvalidName() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Engineering"
  name                        = "Invalid Name With Spaces"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false
}
`
}

func testConfigInvalidBooleanWithCollection() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Engineering"
  name                        = "TestBool"
  type                        = "Boolean"
  status                      = "Available"
  is_collection               = true
  is_searchable               = true
  use_pre_defined_values_only = false
}
`
}

func testConfigInvalidBooleanWithPredefinedValues() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Engineering"
  name                        = "TestBool"
  type                        = "Boolean"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = true
}
`
}

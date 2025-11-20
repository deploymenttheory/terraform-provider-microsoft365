package graphBetaCustomSecurityAttributeDefinition_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaCustomSecurityAttributeDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_definition"
	definitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_definition/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = graphBetaCustomSecurityAttributeDefinition.ResourceName

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
					check.That(resourceType+".test").Key("attribute_set").HasValue("Engineering"),
					check.That(resourceType+".test").Key("name").HasValue("ProjectName"),
					check.That(resourceType+".test").Key("description").HasValue("Name of the project the user is assigned to"),
					check.That(resourceType+".test").Key("type").HasValue("String"),
					check.That(resourceType+".test").Key("status").HasValue("Available"),
					check.That(resourceType+".test").Key("is_collection").HasValue("false"),
					check.That(resourceType+".test").Key("is_searchable").HasValue("true"),
					check.That(resourceType+".test").Key("use_pre_defined_values_only").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_ProjectName",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("attribute_set").HasValue("Engineering"),
					check.That(resourceType+".test").Key("name").HasValue("MinimalAttribute"),
					check.That(resourceType+".test").Key("type").HasValue("String"),
					check.That(resourceType+".test").Key("description").DoesNotExist(),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_MinimalAttribute",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("description").HasValue("Name of the project the user is assigned to"),
					check.That(resourceType+".test").Key("status").HasValue("Available"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("description").HasValue("Updated project name description"),
					check.That(resourceType+".test").Key("status").HasValue("Available"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_ProjectName",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("type").HasValue("Boolean"),
					check.That(resourceType+".test").Key("is_collection").HasValue("false"),
					check.That(resourceType+".test").Key("use_pre_defined_values_only").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Security_HasClearance",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("is_collection").HasValue("true"),
					check.That(resourceType+".test").Key("type").HasValue("String"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "HumanResources_Skills",
				ImportStateVerify: true,
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

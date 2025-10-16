package graphBetaAttributeSet_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	attributeSetMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/attribute_set/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *attributeSetMocks.AttributeSetMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	attributeSetMock := &attributeSetMocks.AttributeSetMock{}
	attributeSetMock.RegisterMocks()
	return mockClient, attributeSetMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *attributeSetMocks.AttributeSetMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	attributeSetMock := &attributeSetMocks.AttributeSetMock{}
	attributeSetMock.RegisterErrorMocks()
	return mockClient, attributeSetMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAttributeSetResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "id", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "description", "Attributes for engineering team"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "max_attributes_per_set", "25"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_attribute_set.test"),
				),
			},
		},
	})
}

func TestAttributeSetResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "id", "Marketing"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "description"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "max_attributes_per_set"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_attribute_set.test"),
				),
			},
		},
	})
}

func TestAttributeSetResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "id", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "description", "Attributes for engineering team"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "max_attributes_per_set", "25"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "id", "Engineering"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "description", "Updated attributes for engineering team"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_attribute_set.test", "max_attributes_per_set", "50"),
				),
			},
		},
	})
}

func TestAttributeSetResource_InvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidID(),
				ExpectError: regexp.MustCompile(`ID cannot contain spaces or special characters`),
			},
		},
	})
}

func TestAttributeSetResource_InvalidMaxAttributes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidMaxAttributes(),
				ExpectError: regexp.MustCompile(`value must be between 1 and 500`),
			},
		},
	})
}

func TestAttributeSetResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, attributeSetMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer attributeSetMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`Invalid attribute set data`),
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

func testConfigInvalidID() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_attribute_set" "test" {
  id = "invalid id with spaces"
  description = "This should fail due to invalid ID"
}
`
}

func testConfigInvalidMaxAttributes() string {
	return `
resource "microsoft365_graph_beta_identity_and_access_attribute_set" "test" {
  id = "TestSet"
  description = "This should fail due to invalid max attributes"
  max_attributes_per_set = 600
}
`
}

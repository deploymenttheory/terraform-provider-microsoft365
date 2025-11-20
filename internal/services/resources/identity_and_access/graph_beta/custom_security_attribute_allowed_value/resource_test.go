package graphBetaCustomSecurityAttributeAllowedValue_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaCustomSecurityAttributeAllowedValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_allowed_value"
	allowedValueMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/custom_security_attribute_allowed_value/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var resourceType = graphBetaCustomSecurityAttributeAllowedValue.ResourceName

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
					check.That(resourceType+".test").Key("custom_security_attribute_definition_id").HasValue("Engineering_Project"),
					check.That(resourceType+".test").Key("id").HasValue("Alpine"),
					check.That(resourceType+".test").Key("is_active").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_Project/Alpine",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("custom_security_attribute_definition_id").HasValue("Engineering_Project"),
					check.That(resourceType+".test").Key("id").HasValue("Minimal"),
					check.That(resourceType+".test").Key("is_active").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_Project/Minimal",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("id").HasValue("Alpine"),
					check.That(resourceType+".test").Key("is_active").HasValue("true"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("Alpine"),
					check.That(resourceType+".test").Key("is_active").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_Project/Alpine",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("custom_security_attribute_definition_id").HasValue("Engineering_Project"),
					check.That(resourceType+".test").Key("id").HasValue("Legacy"),
					check.That(resourceType+".test").Key("is_active").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_Project/Legacy",
				ImportStateVerify: true,
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
					check.That(resourceType+".test").Key("custom_security_attribute_definition_id").HasValue("Engineering_Department"),
					check.That(resourceType+".test").Key("id").HasValue("Human Resources"),
					check.That(resourceType+".test").Key("is_active").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "Engineering_Department/Human Resources",
				ImportStateVerify: true,
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

// Config helper functions
func testConfigBasic() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_basic.tf")
	if err != nil {
		panic("failed to load test_config_basic.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_minimal.tf")
	if err != nil {
		panic("failed to load test_config_minimal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUpdate() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_update.tf")
	if err != nil {
		panic("failed to load test_config_update.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigInactive() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_inactive.tf")
	if err != nil {
		panic("failed to load test_config_inactive.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWithSpaces() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/test_config_with_spaces.tf")
	if err != nil {
		panic("failed to load test_config_with_spaces.tf: " + err.Error())
	}
	return unitTestConfig
}

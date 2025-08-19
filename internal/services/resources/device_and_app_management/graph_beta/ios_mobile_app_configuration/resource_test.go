package graphBetaIOSMobileAppConfiguration_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	iosMobileAppConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_mobile_app_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *iosMobileAppConfigurationMocks.IOSMobileAppConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	iosMobileAppConfigurationMock := &iosMobileAppConfigurationMocks.IOSMobileAppConfigurationMock{}
	iosMobileAppConfigurationMock.RegisterMocks()

	return mockClient, iosMobileAppConfigurationMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *iosMobileAppConfigurationMocks.IOSMobileAppConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	iosMobileAppConfigurationMock := &iosMobileAppConfigurationMocks.IOSMobileAppConfigurationMock{}
	iosMobileAppConfigurationMock.RegisterErrorMocks()

	return mockClient, iosMobileAppConfigurationMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestIOSMobileAppConfigurationResource_Schema validates the resource schema
func TestIOSMobileAppConfigurationResource_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "display_name", "Test Minimal iOS Mobile App Configuration - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "version"),
				),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_Minimal tests basic CRUD operations
func TestIOSMobileAppConfigurationResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "display_name", "Test Minimal iOS Mobile App Configuration - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "display_name", "Test Maximal iOS Mobile App Configuration - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "description", "Maximal iOS mobile app configuration for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "targeted_mobile_apps.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "settings.#", "2"),
				),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_UpdateInPlace tests in-place updates
func TestIOSMobileAppConfigurationResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.minimal", "display_name", "Test Minimal iOS Mobile App Configuration - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "display_name", "Test Maximal iOS Mobile App Configuration - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "description", "Maximal iOS mobile app configuration for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_RequiredFields tests required field validation
func TestIOSMobileAppConfigurationResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  # Missing display_name
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_ErrorHandling tests error scenarios
func TestIOSMobileAppConfigurationResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Mobile App Configuration"
}
`,
				ExpectError: regexp.MustCompile(`Invalid iOS mobile app configuration data|BadRequest`),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_Settings tests settings handling
func TestIOSMobileAppConfigurationResource_Settings(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Mobile App Configuration"
  settings = [
    {
      app_config_key       = "testKey1"
      app_config_key_type  = "stringType"
      app_config_key_value = "testValue1"
    },
    {
      app_config_key       = "testKey2"
      app_config_key_type  = "integerType"
      app_config_key_value = "123"
    }
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "settings.#", "2"),
				),
			},
		},
	})
}

// TestIOSMobileAppConfigurationResource_TargetedMobileAppsValidation tests GUID validation for targeted mobile apps
func TestIOSMobileAppConfigurationResource_TargetedMobileAppsValidation(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, iosMobileAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer iosMobileAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid GUID format
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Mobile App Configuration"
  targeted_mobile_apps = ["invalid-guid", "another-invalid-guid"]
}
`,
				ExpectError: regexp.MustCompile(`Must be a valid GUID format`),
			},
			// Test valid GUID format
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Mobile App Configuration"
  targeted_mobile_apps = ["12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "targeted_mobile_apps.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "targeted_mobile_apps.*", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "targeted_mobile_apps.*", "87654321-4321-4321-4321-210987654321"),
				),
			},
		},
	})
}

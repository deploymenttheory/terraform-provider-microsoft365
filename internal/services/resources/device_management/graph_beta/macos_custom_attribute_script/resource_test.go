package graphBetaMacOSCustomAttributeScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	macosCustomAttributeScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_custom_attribute_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *macosCustomAttributeScriptMocks.MacOSCustomAttributeScriptMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	macosCustomAttributeScriptMock := &macosCustomAttributeScriptMocks.MacOSCustomAttributeScriptMock{}
	macosCustomAttributeScriptMock.RegisterMocks()

	return mockClient, macosCustomAttributeScriptMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *macosCustomAttributeScriptMocks.MacOSCustomAttributeScriptMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	macosCustomAttributeScriptMock := &macosCustomAttributeScriptMocks.MacOSCustomAttributeScriptMock{}
	macosCustomAttributeScriptMock.RegisterErrorMocks()

	return mockClient, macosCustomAttributeScriptMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// TestMacOSCustomAttributeScriptResource_Schema validates the resource schema
func TestUnitResourceMacOSCustomAttributeScript_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosCustomAttributeScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosCustomAttributeScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal", "display_name", "Test Minimal macOS Custom Attribute Script - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestMacOSCustomAttributeScriptResource_Minimal tests basic CRUD operations
func TestUnitResourceMacOSCustomAttributeScript_02_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosCustomAttributeScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosCustomAttributeScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal", "display_name", "Test Minimal macOS Custom Attribute Script - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestMacOSCustomAttributeScriptResource_Maximal tests maximal configuration
func TestUnitResourceMacOSCustomAttributeScript_03_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosCustomAttributeScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosCustomAttributeScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_custom_attribute_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.maximal", "display_name", "Test Maximal macOS Custom Attribute Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.maximal", "description", "Maximal custom attribute script for testing with all features"),
				),
			},
		},
	})
}

// TestMacOSCustomAttributeScriptResource_Update tests update operations
func TestUnitResourceMacOSCustomAttributeScript_04_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosCustomAttributeScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosCustomAttributeScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_custom_attribute_script.minimal"),
				),
			},
			// Update to maximal
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_custom_attribute_script.maximal"),
				),
			},
		},
	})
}

// TestMacOSCustomAttributeScriptResource_ErrorHandling tests error scenarios
func TestUnitResourceMacOSCustomAttributeScript_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosCustomAttributeScriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosCustomAttributeScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Validation error: Invalid display name"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

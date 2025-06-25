package graphBetaSettingsCatalog_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Helper functions to return the test configurations by reading from files
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMinimalToMaximal() string {
	// For minimal to maximal test, we need to use the maximal config
	// but with the minimal resource name to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	return updatedMaximal
}

func testConfigError() string {
	// Read the minimal config and modify for error scenario
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace resource name to create an error scenario
	updated := strings.Replace(string(content), "minimal", "error", 1)
	updated = strings.Replace(updated, "Minimal Settings Catalog", "Error Settings Catalog", 1)

	return updated
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.SettingsCatalogMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	settingsCatalogMock := &localMocks.SettingsCatalogMock{}
	settingsCatalogMock.RegisterMocks()

	return mockClient, settingsCatalogMock
}

// Helper function to check if a resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// TestUnitSettingsCatalogResource_Create_Minimal tests the creation of a settings catalog with minimal configuration
func TestUnitSettingsCatalogResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "name", "Minimal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "description", "Minimal settings catalog policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "platform", "windows10"),
				),
			},
		},
	})
}

// TestUnitSettingsCatalogResource_Create_Maximal tests the creation of a settings catalog with maximal configuration
func TestUnitSettingsCatalogResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "name", "Maximal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "description", "Maximal settings catalog policy with all options"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "platform", "windows10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "technologies", "mdm"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "settings.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.maximal", "assignments.#", "1"),
				),
			},
		},
	})
}

// TestUnitSettingsCatalogResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitSettingsCatalogResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "name", "Minimal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "description", "Minimal settings catalog policy"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "settings.#", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "assignments.#", "0"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "name", "Maximal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "description", "Maximal settings catalog policy with all options"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "platform", "windows10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "technologies", "mdm"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "settings.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.minimal", "assignments.#", "1"),
				),
			},
		},
	})
}

// TestUnitSettingsCatalogResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitSettingsCatalogResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "name", "Maximal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "settings.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "assignments.#", "1"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				// We expect a non-empty plan because computed fields will show as changes
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.test"),
					// Verify it now has only minimal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "name", "Minimal Settings Catalog"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "settings.#", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog.test", "assignments.#", "0"),
				),
			},
		},
	})
}

// Helper function to get maximal config with a custom resource name
func testConfigMaximalWithResourceName(resourceName string) string {
	// Read the maximal config
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name
	updated := strings.Replace(string(content), "maximal", resourceName, 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_settings_catalog" "%s" {
  name        = "Minimal Settings Catalog"
  description = "Minimal settings catalog policy"
  platform    = "windows10"
}`, resourceName)
}

// TestUnitSettingsCatalogResource_Delete_Minimal tests deleting a settings catalog with minimal configuration
func TestUnitSettingsCatalogResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_management_settings_catalog.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitSettingsCatalogResource_Delete_Maximal tests deleting a settings catalog with maximal configuration
func TestUnitSettingsCatalogResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_management_settings_catalog.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitSettingsCatalogResource_Import tests importing a resource
func TestUnitSettingsCatalogResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_settings_catalog.minimal"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitSettingsCatalogResource_Error(t *testing.T) {
	// Set up mock environment
	_, settingsCatalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	settingsCatalogMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Settings catalog policy creation failed"),
			},
		},
	})
}

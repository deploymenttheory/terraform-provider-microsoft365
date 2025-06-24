package graphM365AppsInstallationOptions_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/m365_admin/graph_beta/m365_apps_installation_options/mocks"
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
	// Create an error configuration with invalid update_channel
	return `
resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "error" {
  update_channel = "invalid"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
}
`
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.M365AppsInstallationOptionsMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	m365AppsMock := &localMocks.M365AppsInstallationOptionsMock{}
	m365AppsMock.RegisterMocks()

	return mockClient, m365AppsMock
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
	return fmt.Sprintf(`resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "%s" {
  update_channel = "current"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
}`, resourceName)
}

// TestUnitM365AppsInstallationOptionsResource_Create_Minimal tests the creation of M365 Apps Installation Options with minimal configuration
func TestUnitM365AppsInstallationOptionsResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "update_channel", "current"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_windows.is_skype_for_business_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_mac.is_microsoft_365_apps_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_mac.is_skype_for_business_enabled", "true"),
				),
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Create_Maximal tests the creation of M365 Apps Installation Options with maximal configuration
func TestUnitM365AppsInstallationOptionsResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal", "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal", "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal", "apps_for_windows.is_skype_for_business_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal", "apps_for_mac.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal", "apps_for_mac.is_skype_for_business_enabled", "false"),
				),
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitM365AppsInstallationOptionsResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "update_channel", "current"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal", "apps_for_mac.is_microsoft_365_apps_enabled", "false"),
				),
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitM365AppsInstallationOptionsResource_Update_MaximalToMinimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.test", "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.test", "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.test", "update_channel", "current"),
					resource.TestCheckResourceAttr("microsoft365_graph_m365_admin_m365_apps_installation_options.test", "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
				),
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Delete_Minimal tests deleting M365 Apps Installation Options with minimal configuration
func TestUnitM365AppsInstallationOptionsResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Delete_Maximal tests deleting M365 Apps Installation Options with maximal configuration
func TestUnitM365AppsInstallationOptionsResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_m365_admin_m365_apps_installation_options.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Import tests importing a resource
func TestUnitM365AppsInstallationOptionsResource_Import(t *testing.T) {
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
					testCheckExists("microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"),
				),
			},
			// Import
			{
				ResourceName:            "microsoft365_graph_m365_admin_m365_apps_installation_options.minimal",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// TestUnitM365AppsInstallationOptionsResource_Error tests error handling
func TestUnitM365AppsInstallationOptionsResource_Error(t *testing.T) {
	// Set up mock environment
	_, m365AppsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	m365AppsMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Attribute update_channel value must be one of"),
			},
		},
	})
}

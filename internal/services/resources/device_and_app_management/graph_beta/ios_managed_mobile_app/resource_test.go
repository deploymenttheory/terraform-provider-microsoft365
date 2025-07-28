package graphBetaDeviceAndAppManagementIOSManagedMobileApp_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_mobile_app/mocks"
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
	// but with the minimal resource name and managed_app_protection_id to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	// Replace the managed_app_protection_id to match the minimal one
	updatedMaximal = strings.Replace(updatedMaximal, "00000000-0000-0000-0000-000000000003", "00000000-0000-0000-0000-000000000002", 1)

	return updatedMaximal
}

func testConfigError() string {
	// Create an error configuration with invalid managed_app_protection_id
	return `
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "error" {
  managed_app_protection_id = "invalid-id"
  mobile_app_identifier = {
    bundle_id = "com.example.testapp"
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
func setupMockEnvironment() (*mocks.Mocks, *localMocks.IOSManagedMobileAppMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	managedAppMock := &localMocks.IOSManagedMobileAppMock{}
	managedAppMock.RegisterMocks()

	return mockClient, managedAppMock
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
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "%s" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000003"
  mobile_app_identifier = {
    bundle_id = "com.example.testapp"
  }
}`, resourceName)
}

// TestUnitIOSManagedMobileAppResource_Create_Minimal tests the creation of an iOS managed mobile app with minimal configuration
func TestUnitIOSManagedMobileAppResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "managed_app_protection_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "mobile_app_identifier.bundle_id", "com.example.testapp"),
				),
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Create_Maximal tests the creation of an iOS managed mobile app with maximal configuration
func TestUnitIOSManagedMobileAppResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal", "managed_app_protection_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal", "mobile_app_identifier.bundle_id", "com.example.complexapp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal", "version", "1.5"),
				),
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitIOSManagedMobileAppResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "managed_app_protection_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "mobile_app_identifier.bundle_id", "com.example.testapp"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "managed_app_protection_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "mobile_app_identifier.bundle_id", "com.example.complexapp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal", "version", "1.5"),
				),
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitIOSManagedMobileAppResource_Update_MaximalToMinimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.test", "version", "1.5"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.test", "mobile_app_identifier.bundle_id", "com.example.testapp"),
				),
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Delete_Minimal tests deleting an iOS managed mobile app with minimal configuration
func TestUnitIOSManagedMobileAppResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Delete_Maximal tests deleting an iOS managed mobile app with maximal configuration
func TestUnitIOSManagedMobileAppResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Import tests importing a resource
func TestUnitIOSManagedMobileAppResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Use a pre-existing app ID that exists in the mock
	importConfig := `
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "import_test" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000002"
  mobile_app_identifier = {
    bundle_id = "com.example.testapp"
  }
}
`

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import using the composite ID format: managedAppProtectionId/appId
			{
				Config:        importConfig,
				ResourceName:  "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.import_test",
				ImportState:   true,
				ImportStateId: "00000000-0000-0000-0000-000000000002/00000000-0000-0000-0000-000000000001",
				ImportStateCheck: func(state []*terraform.InstanceState) error {
					if len(state) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(state))
					}
					if state[0].ID != "00000000-0000-0000-0000-000000000001" {
						return fmt.Errorf("expected ID '00000000-0000-0000-0000-000000000001', got '%s'", state[0].ID)
					}
					if state[0].Attributes["managed_app_protection_id"] != "00000000-0000-0000-0000-000000000002" {
						return fmt.Errorf("expected managed_app_protection_id '00000000-0000-0000-0000-000000000002', got '%s'", state[0].Attributes["managed_app_protection_id"])
					}
					return nil
				},
			},
		},
	})
}

// TestUnitIOSManagedMobileAppResource_Error tests error handling
func TestUnitIOSManagedMobileAppResource_Error(t *testing.T) {
	// Set up mock environment
	_, managedAppMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	managedAppMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Error creating iOS managed mobile app"),
			},
		},
	})
}
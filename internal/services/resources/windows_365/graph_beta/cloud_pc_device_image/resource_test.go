package graphBetaCloudPcDeviceImage_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloudPcDeviceImageMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_device_image/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*cloudPcDeviceImageMocks.CloudPcDeviceImageMock, *cloudPcDeviceImageMocks.CloudPcDeviceImageMock) {
	httpmock.Activate()
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	errorMock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	return mock, errorMock
}

func setupTestEnvironment(t *testing.T) {
	// Set up any test-specific environment variables or configurations here if needed
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
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

	// Fix the display name to match test expectations
	updated = strings.Replace(updated, "Test Maximal Cloud PC Device Image - Unique", "Test Maximal Cloud PC Device Image", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "%s" {
  display_name              = "Test Minimal Cloud PC Device Image"
  version                   = "1.0.0"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-minimal-image"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitCloudPcDeviceImageResource_Create_Minimal tests the creation of a Cloud PC device image with minimal configuration
func TestUnitCloudPcDeviceImageResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "display_name", "Test Minimal Cloud PC Device Image - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "version", "1.0.0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "source_image_resource_id", "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-minimal-image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "status", "ready"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "os_status", "supported"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "operating_system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "os_build_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "os_version_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal", "last_modified_date_time"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Create_Maximal tests the creation of a Cloud PC device image with maximal configuration
func TestUnitCloudPcDeviceImageResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "display_name", "Test Maximal Cloud PC Device Image - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "version", "2.1.5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "source_image_resource_id", "/subscriptions/87654321-4321-4321-4321-cba987654321/resourceGroups/test-maximal-rg/providers/Microsoft.Compute/images/test-maximal-image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "status", "ready"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "os_status", "supported"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "operating_system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "os_build_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "os_version_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal", "last_modified_date_time"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitCloudPcDeviceImageResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Minimal Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "1.0.0"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Maximal Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "2.1.5"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitCloudPcDeviceImageResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Maximal Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "2.1.5"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Minimal Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "1.0.0"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Delete_Minimal tests deleting a Cloud PC device image with minimal configuration
func TestUnitCloudPcDeviceImageResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Delete_Maximal tests deleting a Cloud PC device image with maximal configuration
func TestUnitCloudPcDeviceImageResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.maximal"),
				),
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Import tests importing a Cloud PC device image
func TestUnitCloudPcDeviceImageResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcDeviceImageMocks.CloudPcDeviceImageMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_device_image.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitCloudPcDeviceImageResource_Error tests error handling
func TestUnitCloudPcDeviceImageResource_Error(t *testing.T) {
	// Set up mock environment
	_, errorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the error mocks
	errorMock.RegisterErrorMocks()

	// Run the test
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

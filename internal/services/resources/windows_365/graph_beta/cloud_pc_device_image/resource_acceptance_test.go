package graphBetaCloudPcDeviceImage_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudPcDeviceImageResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccCloudPcDeviceImageConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Acceptance Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "1.0.0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "source_image_resource_id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "status"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "operating_system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "os_build_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "os_version_number"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "last_modified_date_time"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_device_image.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccCloudPcDeviceImageConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Acceptance Cloud PC Device Image - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "2.0.0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "source_image_resource_id"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccCloudPcDeviceImageConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "display_name", "Test Acceptance Cloud PC Device Image"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_device_image.test", "version", "1.0.0"),
				),
			},
		},
	})
}

func TestAccCloudPcDeviceImageResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcDeviceImageConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcDeviceImageConfig_missingVersion(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcDeviceImageConfig_missingSourceImageResourceId(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccCloudPcDeviceImageResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcDeviceImageConfig_invalidSourceImageResourceId(),
				ExpectError: regexp.MustCompile("Must be a valid Azure image resource ID"),
			},
			{
				Config:      testAccCloudPcDeviceImageConfig_invalidVersionTooLong(),
				ExpectError: regexp.MustCompile("string length must be between 1 and 32"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("M365_TENANT_ID") == "" {
		t.Skip("M365_TENANT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_ID") == "" {
		t.Skip("M365_CLIENT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_SECRET") == "" {
		t.Skip("M365_CLIENT_SECRET must be set for acceptance tests")
	}
}

func testAccCloudPcDeviceImageConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name              = "Test Acceptance Cloud PC Device Image"
  version                   = "1.0.0"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-image"
}
`
}

func testAccCloudPcDeviceImageConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name              = "Test Acceptance Cloud PC Device Image - Updated"
  version                   = "2.0.0"
  source_image_resource_id  = "/subscriptions/87654321-4321-4321-4321-cba987654321/resourceGroups/test-maximal-rg/providers/Microsoft.Compute/images/test-maximal-image"
}
`
}

func testAccCloudPcDeviceImageConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  version                   = "1.0.0"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-image"
}
`
}

func testAccCloudPcDeviceImageConfig_missingVersion() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name              = "Test Image"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-image"
}
`
}

func testAccCloudPcDeviceImageConfig_missingSourceImageResourceId() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name = "Test Image"
  version      = "1.0.0"
}
`
}

func testAccCloudPcDeviceImageConfig_invalidSourceImageResourceId() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name              = "Test Image"
  version                   = "1.0.0"
  source_image_resource_id  = "invalid-resource-id"
}
`
}

func testAccCloudPcDeviceImageConfig_invalidVersionTooLong() string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "test" {
  display_name              = "Test Image"
  version                   = "%s"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-image"
}
`, "this-version-string-is-way-too-long-and-exceeds-the-32-character-limit")
}
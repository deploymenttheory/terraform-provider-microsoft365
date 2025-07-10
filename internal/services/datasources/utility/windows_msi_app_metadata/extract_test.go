package utilityWindowsMSIAppMetadata_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityWindowsMSIAppMetadata "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/windows_msi_app_metadata"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// Firefox MSI download URL
	firefoxMSIURL = "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=en-US"
)

// Helper functions to return the test configurations by reading from files
func testConfigFirefoxMSI() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_firefox_msi.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigLocalMSI() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "datasource_local_msi.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to download Firefox MSI and get its path
func downloadFirefoxMSI(t *testing.T) string {
	t.Helper()

	filePath, err := common.DownloadFile(firefoxMSIURL)
	if err != nil {
		t.Fatalf("Failed to download Firefox MSI: %v", err)
	}

	// Verify the file was downloaded
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %s", filePath)
	}

	t.Logf("Downloaded Firefox MSI to: %s", filePath)
	return filePath
}

// Helper function to clean up downloaded file
func cleanupDownloadedFile(t *testing.T, filePath string) {
	t.Helper()

	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			t.Logf("Warning: Failed to remove downloaded file %s: %v", filePath, err)
		} else {
			t.Logf("Cleaned up downloaded file: %s", filePath)
		}
	}
}

// Helper function to create temporary terraform config with file path
func createTerraformConfigWithPath(t *testing.T, filePath string) string {
	t.Helper()

	// Create a temporary terraform config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.tf")

	config := fmt.Sprintf(`
data "microsoft365_utility_windows_msi_app_metadata" "firefox" {
  installer_file_path_source = "%s"
}
`, filePath)

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to create temporary terraform config: %v", err)
	}

	return config
}

// TestUnitWindowsMSIAppMetadataDataSource_FirefoxMSI tests downloading and extracting Firefox MSI metadata
func TestUnitWindowsMSIAppMetadataDataSource_FirefoxMSI(t *testing.T) {

	setupTestEnvironment(t)

	msiPath := downloadFirefoxMSI(t)
	defer cleanupDownloadedFile(t, msiPath)

	config := createTerraformConfigWithPath(t, msiPath)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Check that the data source has an ID
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "id"),

					// Check that metadata was extracted
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_version"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_code"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.publisher"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.architecture"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.sha256_checksum"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.md5_checksum"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.size_mb"),

					// Check that Firefox-specific values are present
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.publisher", "Mozilla"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.architecture", "Unknown"),

					// Check that commands were generated
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.install_command"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.uninstall_command"),

					// Check that properties map is populated
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.properties.%"),
				),
			},
		},
	})
}

// TestUnitWindowsMSIAppMetadataDataSource_LocalMSI tests using a local MSI file (if available)
func TestUnitWindowsMSIAppMetadataDataSource_LocalMSI(t *testing.T) {
	// Skip this test if no test MSI file is available
	testMSIPath := "testdata/sample.msi"
	if _, err := os.Stat(testMSIPath); os.IsNotExist(err) {
		t.Skip("Test MSI file not found, skipping test")
	}

	// Set up the test environment
	setupTestEnvironment(t)

	// Create terraform config with the local file path
	config := createTerraformConfigWithPath(t, testMSIPath)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Check that the data source has an ID
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "id"),

					// Check that metadata was extracted
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_version"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_code"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.publisher"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.architecture"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.sha256_checksum"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.md5_checksum"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.size_mb"),
				),
			},
		},
	})
}

// Utility function to print all extracted metadata (useful for debugging)
func PrintMetadata(metadata *utilityWindowsMSIAppMetadata.MetadataDataSourceModel) {
	fmt.Println("=== MSI Metadata ===")
	fmt.Printf("Product Name: %s\n", getStringValue(metadata.ProductName))
	fmt.Printf("Product Version: %s\n", getStringValue(metadata.ProductVersion))
	fmt.Printf("Product Code: %s\n", getStringValue(metadata.ProductCode))
	fmt.Printf("Publisher: %s\n", getStringValue(metadata.Publisher))
	fmt.Printf("Upgrade Code: %s\n", getStringValue(metadata.UpgradeCode))
	fmt.Printf("Language: %s\n", getStringValue(metadata.Language))
	fmt.Printf("Package Type: %s\n", getStringValue(metadata.PackageType))
	fmt.Printf("Install Location: %s\n", getStringValue(metadata.InstallLocation))
	fmt.Printf("Architecture: %s\n", getStringValue(metadata.Architecture))
	fmt.Printf("Min OS Version: %s\n", getStringValue(metadata.MinOSVersion))

	if !metadata.SizeMB.IsNull() {
		fmt.Printf("Size (MB): %.2f\n", metadata.SizeMB.ValueFloat64())
	}

	fmt.Printf("SHA256: %s\n", getStringValue(metadata.SHA256Checksum))
	fmt.Printf("MD5: %s\n", getStringValue(metadata.MD5Checksum))
	fmt.Printf("Install Command: %s\n", getStringValue(metadata.InstallCommand))
	fmt.Printf("Uninstall Command: %s\n", getStringValue(metadata.UninstallCommand))

	if !metadata.Properties.IsNull() {
		fmt.Printf("Total Properties: %d\n", len(metadata.Properties.Elements()))
	}

	if !metadata.Files.IsNull() {
		fmt.Printf("Total Files: %d\n", len(metadata.Files.Elements()))
	}

	if !metadata.RequiredFeatures.IsNull() {
		fmt.Printf("Total Features: %d\n", len(metadata.RequiredFeatures.Elements()))
	}
}

// Helper to safely get string values
func getStringValue(attr types.String) string {
	if attr.IsNull() {
		return "<null>"
	}
	return attr.ValueString()
}

// Example Terraform configuration for using this data source
const ExampleTerraformConfig = `
# Extract metadata from a local MSI file
data "microsoft365_utility_windows_msi_app_metadata" "local_msi" {
  installer_file_path_source = "C:/path/to/your/installer.msi"
}

# Extract metadata from a remote MSI file
data "microsoft365_utility_windows_msi_app_metadata" "remote_msi" {
  installer_url_source = "https://example.com/path/to/installer.msi"
}

# Use the extracted metadata
output "msi_metadata" {
  value = {
    product_name      = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.product_name
    product_version   = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.product_version
    product_code      = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.product_code
    publisher         = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.publisher
    upgrade_code      = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.upgrade_code
    architecture      = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.architecture
    install_command   = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.install_command
    uninstall_command = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.uninstall_command
    size_mb           = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.size_mb
    sha256_checksum   = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.sha256_checksum
    files_count       = length(data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.files)
    features_count    = length(data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.required_features)
  }
}

# Example of using the metadata to create a Microsoft Graph application
resource "microsoft365_graph_application" "msi_app" {
  display_name = data.microsoft365_utility_windows_msi_app_metadata.local_msi.metadata.product_name
  
  # Use other metadata fields as needed...
}
`

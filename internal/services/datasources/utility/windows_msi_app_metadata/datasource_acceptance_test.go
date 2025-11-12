package utilityWindowsMSIAppMetadata_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsMSIAppMetadataDataSource_FirefoxMSI(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigFirefoxMSI(),
				Check: resource.ComposeTestCheckFunc(
					// Core attributes
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "installer_url_source", "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=en-US"),

					// Core MSI properties
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_code"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_code", regexp.MustCompile(`^{[A-F0-9-]+}$`)),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_version"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.product_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.publisher"),

					// Commands
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.install_command"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.install_command", regexp.MustCompile(`msiexec /i`)),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.uninstall_command"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.uninstall_command", regexp.MustCompile(`msiexec /x`)),

					// File information
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.size_mb"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.sha256_checksum"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.sha256_checksum", regexp.MustCompile(`^[a-f0-9]{64}$`)),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.md5_checksum"),
					resource.TestMatchResourceAttr("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.md5_checksum", regexp.MustCompile(`^[a-f0-9]{32}$`)),

					// Additional metadata
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.upgrade_code"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.language"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_windows_msi_app_metadata.firefox", "metadata.architecture"),
				),
			},
		},
	})
}

// Acceptance test configuration function
func testAccConfigFirefoxMSI() string {
	accTestConfig, err := helpers.ParseHCLFile("mocks/terraform/datasource_firefox_msi.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

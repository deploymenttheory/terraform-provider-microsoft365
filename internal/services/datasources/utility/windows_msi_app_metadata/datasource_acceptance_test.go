package utilityWindowsMSIAppMetadata_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
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
					check.That("data."+dataSourceType+".firefox").Key("id").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("installer_url_source").HasValue("https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=en-US"),

					// Core MSI properties
					check.That("data."+dataSourceType+".firefox").Key("metadata.product_code").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.product_code").MatchesRegex(regexp.MustCompile(`^{[A-F0-9-]+}$`)),
					check.That("data."+dataSourceType+".firefox").Key("metadata.product_version").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.product_name").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.publisher").IsSet(),

					// Commands
					check.That("data."+dataSourceType+".firefox").Key("metadata.install_command").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.install_command").MatchesRegex(regexp.MustCompile(`msiexec /i`)),
					check.That("data."+dataSourceType+".firefox").Key("metadata.uninstall_command").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.uninstall_command").MatchesRegex(regexp.MustCompile(`msiexec /x`)),

					// File information
					check.That("data."+dataSourceType+".firefox").Key("metadata.size_mb").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.sha256_checksum").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.sha256_checksum").MatchesRegex(regexp.MustCompile(`^[a-f0-9]{64}$`)),
					check.That("data."+dataSourceType+".firefox").Key("metadata.md5_checksum").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.md5_checksum").MatchesRegex(regexp.MustCompile(`^[a-f0-9]{32}$`)),

					// Additional metadata
					check.That("data."+dataSourceType+".firefox").Key("metadata.upgrade_code").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.language").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("metadata.architecture").IsSet(),
				),
			},
		},
	})
}

// Acceptance test configuration function
func testAccConfigFirefoxMSI() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/datasource_firefox_msi.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

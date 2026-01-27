package graphBetaMobileAppCatalogPackage_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// Test 01: Get all mobile app catalog packages
func TestAccDatasourceMobileAppCatalogPackage_01_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").Exists(),

					// Verify we have valid data structure for at least first item
					check.That(dataSourceType+".all").Key("items.0.id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.publisher").Exists(),
					check.That(dataSourceType+".all").Key("items.0.file_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.mobile_app_catalog_package_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),

					// Verify core fields are populated
					check.That(dataSourceType+".all").Key("items.0.install_command_line").Exists(),
					check.That(dataSourceType+".all").Key("items.0.uninstall_command_line").Exists(),
					check.That(dataSourceType+".all").Key("items.0.display_version").Exists(),

					// Verify nested structures exist
					check.That(dataSourceType+".all").Key("items.0.install_experience.run_as_account").Exists(),
					check.That(dataSourceType+".all").Key("items.0.install_experience.max_run_time_in_minutes").Exists(),
					check.That(dataSourceType+".all").Key("items.0.install_experience.device_restart_behavior").Exists(),
					check.That(dataSourceType+".all").Key("items.0.return_codes.#").Exists(),
					check.That(dataSourceType+".all").Key("items.0.rules.#").Exists(),
				),
				// Use ExpectNonEmptyPlan to allow for dynamic catalog changes between API calls
				// The Microsoft catalog is constantly updated with new packages
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test 02: Get by product ID (using 7-Zip as a known stable package)
func TestAccDatasourceMobileAppCatalogPackage_02_ById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("3a6307ef-6991-faf1-01e1-35e1557287aa"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".by_id").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`7-Zip`)),
					check.That(dataSourceType+".by_id").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`Igor Pavlov`)),
					check.That(dataSourceType+".by_id").Key("items.0.mobile_app_catalog_package_id").Exists(),

					// Verify complete win32CatalogApp structure
					check.That(dataSourceType+".by_id").Key("items.0.file_name").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.size").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.install_command_line").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.uninstall_command_line").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.setup_file_path").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.minimum_supported_windows_release").Exists(),

					// Verify rules exist
					check.That(dataSourceType+".by_id").Key("items.0.rules.#").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.rules.0.rule_type").Exists(),

					// Verify install experience
					check.That(dataSourceType+".by_id").Key("items.0.install_experience.run_as_account").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.install_experience.device_restart_behavior").Exists(),

					// Verify return codes
					check.That(dataSourceType+".by_id").Key("items.0.return_codes.#").Exists(),
				),
			},
		},
	})
}

// Test 03: Get by product name
func TestAccDatasourceMobileAppCatalogPackage_03_ByProductName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_by_product_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_product_name").Key("filter_type").HasValue("product_name"),
					check.That(dataSourceType+".by_product_name").Key("filter_value").HasValue("7-Zip"),
					check.That(dataSourceType+".by_product_name").Key("items.#").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)7-Zip`)),
					check.That(dataSourceType+".by_product_name").Key("items.0.mobile_app_catalog_package_id").Exists(),
				),
			},
		},
	})
}

// Test 04: Get by publisher name
func TestAccDatasourceMobileAppCatalogPackage_04_ByPublisherName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_by_publisher_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_publisher").Key("filter_type").HasValue("publisher_name"),
					check.That(dataSourceType+".by_publisher").Key("filter_value").HasValue("Microsoft"),
					check.That(dataSourceType+".by_publisher").Key("items.#").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
				),
			},
		},
	})
}

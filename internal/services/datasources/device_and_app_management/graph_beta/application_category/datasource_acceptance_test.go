package graphBetaApplicationCategory_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccApplicationCategoryDataSource_All tests fetching all application categories from live API
func TestAccApplicationCategoryDataSource_All(t *testing.T) {
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
				Config: testAccConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").Exists(),

					// Verify required fields for at least first item
					check.That(dataSourceType+".all").Key("items.0.id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".all").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.last_modified_date_time").Exists(),
				),
			},
		},
	})
}

// TestAccApplicationCategoryDataSource_ByDisplayName tests filtering categories by display name from live API
func TestAccApplicationCategoryDataSource_ByDisplayName(t *testing.T) {
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
				Config: testAccConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_display_name").Key("filter_type").HasValue("display_name"),
					check.That(dataSourceType+".by_display_name").Key("filter_value").HasValue("Business"),
					check.That(dataSourceType+".by_display_name").Key("items.#").Exists(),

					// Verify at least one item contains "Business" in display name
					check.That(dataSourceType+".by_display_name").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)Business`)),
					check.That(dataSourceType+".by_display_name").Key("items.0.last_modified_date_time").Exists(),
				),
			},
		},
	})
}

// TestAccApplicationCategoryDataSource_ODataFilter tests using OData filter queries from live API
func TestAccApplicationCategoryDataSource_ODataFilter(t *testing.T) {
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
				Config: testAccConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("startswith(displayName, 'Business')"),
					check.That(dataSourceType+".odata_filter").Key("items.#").Exists(),

					// Verify filtered results have required fields
					check.That(dataSourceType+".odata_filter").Key("items.0.id").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".odata_filter").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.last_modified_date_time").Exists(),

					// Verify display name starts with "Business"
					check.That(dataSourceType+".odata_filter").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`^Business`)),
				),
			},
		},
	})
}

// Helper functions to load acceptance test configurations
func testAccConfigAll() string {
	config := mocks.LoadTerraformConfigFile("01_all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigByDisplayName() string {
	config := mocks.LoadTerraformConfigFile("02_by_display_name.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigODataFilter() string {
	config := mocks.LoadTerraformConfigFile("03_odata_filter.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

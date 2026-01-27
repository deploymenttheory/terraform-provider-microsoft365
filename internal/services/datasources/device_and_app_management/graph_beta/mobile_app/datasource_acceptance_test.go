package graphBetaMobileApp_test

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

// Test 01: Get all mobile apps
func TestAccDatasourceMobileApp_01_All(t *testing.T) {
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

					// Verify required fields for at least first item
					check.That(dataSourceType+".all").Key("items.0.id").Exists(),
					check.That(dataSourceType+".all").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".all").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".all").Key("items.0.publisher").Exists(),

					// Verify core timestamp fields
					check.That(dataSourceType+".all").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".all").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".all").Key("items.0.publishing_state").Exists(),

					// Verify boolean fields (always present, may be true/false)
					check.That(dataSourceType+".all").Key("items.0.is_assigned").Exists(),
					check.That(dataSourceType+".all").Key("items.0.is_featured").Exists(),

					// Verify numeric fields (always present, may be 0)
					check.That(dataSourceType+".all").Key("items.0.upload_state").Exists(),
					check.That(dataSourceType+".all").Key("items.0.dependent_app_count").Exists(),
					check.That(dataSourceType+".all").Key("items.0.superseding_app_count").Exists(),
					check.That(dataSourceType+".all").Key("items.0.superseded_app_count").Exists(),

					// Verify lists exist (may be empty)
					check.That(dataSourceType+".all").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".all").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 02: Get by display name
func TestAccDatasourceMobileApp_02_ByDisplayName(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("02_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_display_name").Key("filter_type").HasValue("display_name"),
					check.That(dataSourceType+".by_display_name").Key("filter_value").HasValue("Microsoft"),
					check.That(dataSourceType+".by_display_name").Key("items.#").Exists(),

					// Verify at least one item contains "Microsoft" in display name
					check.That(dataSourceType+".by_display_name").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
					check.That(dataSourceType+".by_display_name").Key("items.0.publisher").Exists(),

					// Verify core app fields
					check.That(dataSourceType+".by_display_name").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.publishing_state").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.is_assigned").Exists(),

					// Verify lists exist (may be empty)
					check.That(dataSourceType+".by_display_name").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 03: Get by publisher name
func TestAccDatasourceMobileApp_03_ByPublisherName(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("03_by_publisher_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_publisher").Key("filter_type").HasValue("publisher_name"),
					check.That(dataSourceType+".by_publisher").Key("filter_value").HasValue("Microsoft"),
					check.That(dataSourceType+".by_publisher").Key("items.#").Exists(),

					// Verify at least one item has "Microsoft" in publisher
					check.That(dataSourceType+".by_publisher").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".by_publisher").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),

					// Verify core app fields
					check.That(dataSourceType+".by_publisher").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.publishing_state").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.is_assigned").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.is_featured").Exists(),

					// Verify numeric metadata fields (always present, may be 0)
					check.That(dataSourceType+".by_publisher").Key("items.0.upload_state").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.dependent_app_count").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.superseding_app_count").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.superseded_app_count").Exists(),

					// Verify lists exist (may be empty)
					check.That(dataSourceType+".by_publisher").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

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
					check.That(dataSourceType+".all").Key("list_all").HasValue("true"),
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

// Test 02: Get by app ID
func TestAccDatasourceMobileApp_02_ByAppId(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("02_by_app_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("app_id").Exists(),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".by_id").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".by_id").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.publisher").Exists(),

					check.That(dataSourceType+".by_id").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.publishing_state").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.is_assigned").Exists(),

					check.That(dataSourceType+".by_id").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 03: OData filter
func TestAccDatasourceMobileApp_03_ODataFilter(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("03_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".odata_filter").Key("odata_query").HasValue("startswith(publisher, 'Microsoft')"),
					check.That(dataSourceType+".odata_filter").Key("items.#").Exists(),

					check.That(dataSourceType+".odata_filter").Key("items.0.id").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),

					check.That(dataSourceType+".odata_filter").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.publishing_state").Exists(),

					check.That(dataSourceType+".odata_filter").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".odata_filter").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 04: With app type filter
func TestAccDatasourceMobileApp_04_WithAppTypeFilter(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("04_with_app_type_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".win32_apps").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".win32_apps").Key("app_type_filter").HasValue("win32LobApp"),
					check.That(dataSourceType+".win32_apps").Key("items.#").Exists(),

					check.That(dataSourceType+".win32_apps").Key("items.0.id").Exists(),
					check.That(dataSourceType+".win32_apps").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".win32_apps").Key("items.0.publisher").Exists(),

					check.That(dataSourceType+".win32_apps").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".win32_apps").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".win32_apps").Key("items.0.publishing_state").Exists(),

					check.That(dataSourceType+".win32_apps").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".win32_apps").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 05: Get by display name
func TestAccDatasourceMobileApp_05_ByDisplayName(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("05_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_display_name").Key("display_name").HasValue("Microsoft"),
					check.That(dataSourceType+".by_display_name").Key("items.#").Exists(),

					check.That(dataSourceType+".by_display_name").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),
					check.That(dataSourceType+".by_display_name").Key("items.0.publisher").Exists(),

					check.That(dataSourceType+".by_display_name").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.publishing_state").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.is_assigned").Exists(),

					check.That(dataSourceType+".by_display_name").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_display_name").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 06: Get by publisher name
func TestAccDatasourceMobileApp_06_ByPublisherName(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("06_by_publisher_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_publisher").Key("publisher").HasValue("Microsoft"),
					check.That(dataSourceType+".by_publisher").Key("items.#").Exists(),

					check.That(dataSourceType+".by_publisher").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".by_publisher").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.publisher").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),

					check.That(dataSourceType+".by_publisher").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.publishing_state").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.is_assigned").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.is_featured").Exists(),

					check.That(dataSourceType+".by_publisher").Key("items.0.upload_state").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.dependent_app_count").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.superseding_app_count").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.superseded_app_count").Exists(),

					check.That(dataSourceType+".by_publisher").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 07: Get by developer
func TestAccDatasourceMobileApp_07_ByDeveloper(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("07_by_developer.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_developer").Key("developer").HasValue("Microsoft"),
					check.That(dataSourceType+".by_developer").Key("items.#").Exists(),

					check.That(dataSourceType+".by_developer").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_developer").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_developer").Key("items.0.developer").MatchesRegex(regexp.MustCompile(`(?i)Microsoft`)),

					check.That(dataSourceType+".by_developer").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_developer").Key("items.0.last_modified_date_time").Exists(),
					check.That(dataSourceType+".by_developer").Key("items.0.publishing_state").Exists(),

					check.That(dataSourceType+".by_developer").Key("items.0.role_scope_tag_ids.#").Exists(),
					check.That(dataSourceType+".by_developer").Key("items.0.categories.#").Exists(),
				),
			},
		},
	})
}

// Test 08: Get by category
func TestAccDatasourceMobileApp_08_ByCategory(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("08_by_category.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_category").Key("category").HasValue("Productivity"),
					check.That(dataSourceType+".by_category").Key("items.#").Exists(),

					check.That(dataSourceType+".by_category").Key("items.0.id").Exists(),
					check.That(dataSourceType+".by_category").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".by_category").Key("items.0.categories.#").Exists(),

					check.That(dataSourceType+".by_category").Key("items.0.created_date_time").Exists(),
					check.That(dataSourceType+".by_category").Key("items.0.last_modified_date_time").Exists(),

					check.That(dataSourceType+".by_category").Key("items.0.role_scope_tag_ids.#").Exists(),
				),
			},
		},
	})
}

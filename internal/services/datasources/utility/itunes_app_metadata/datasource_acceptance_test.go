package itunes_app_metadata_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccItunesAppMetadataDataSource_Firefox(t *testing.T) {
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
				Config: testAccConfigFirefox(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".firefox").Key("id").HasValue("us_firefox"),
					check.That("data."+dataSourceType+".firefox").Key("search_term").HasValue("firefox"),
					check.That("data."+dataSourceType+".firefox").Key("country_code").HasValue("us"),
					check.That("data."+dataSourceType+".firefox").Key("results.#").IsSet(),
					// Verify we get at least one result
					check.That("data."+dataSourceType+".firefox").Key("results.0.track_id").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("results.0.track_name").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("results.0.bundle_id").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("results.0.seller_name").IsSet(),
					check.That("data."+dataSourceType+".firefox").Key("results.0.version").IsSet(),
				),
			},
		},
	})
}

func TestAccItunesAppMetadataDataSource_Office(t *testing.T) {
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
				Config: testAccConfigOffice(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".office").Key("id").HasValue("us_microsoft office"),
					check.That("data."+dataSourceType+".office").Key("search_term").HasValue("microsoft office"),
					check.That("data."+dataSourceType+".office").Key("country_code").HasValue("us"),
					check.That("data."+dataSourceType+".office").Key("results.#").IsSet(),
					// Verify we get at least one result
					check.That("data."+dataSourceType+".office").Key("results.0.track_id").IsSet(),
					check.That("data."+dataSourceType+".office").Key("results.0.track_name").IsSet(),
					check.That("data."+dataSourceType+".office").Key("results.0.bundle_id").IsSet(),
					check.That("data."+dataSourceType+".office").Key("results.0.seller_name").IsSet(),
				),
			},
		},
	})
}

func TestAccItunesAppMetadataDataSource_Teams(t *testing.T) {
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
				Config: testAccConfigTeams(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".teams").Key("id").HasValue("us_microsoft teams"),
					check.That("data."+dataSourceType+".teams").Key("search_term").HasValue("microsoft teams"),
					check.That("data."+dataSourceType+".teams").Key("country_code").HasValue("us"),
					check.That("data."+dataSourceType+".teams").Key("results.#").IsSet(),
					// Verify we get at least one result
					check.That("data."+dataSourceType+".teams").Key("results.0.track_id").IsSet(),
					check.That("data."+dataSourceType+".teams").Key("results.0.track_name").IsSet(),
					check.That("data."+dataSourceType+".teams").Key("results.0.bundle_id").IsSet(),
					check.That("data."+dataSourceType+".teams").Key("results.0.seller_name").IsSet(),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigFirefox() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/datasource_firefox.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigOffice() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/datasource_office.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigTeams() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/datasource_teams.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

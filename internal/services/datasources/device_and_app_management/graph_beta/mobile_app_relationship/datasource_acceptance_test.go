package graphBetaMobileAppRelationship_test

import (
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

// TestAccMobileAppRelationshipDataSource_All tests fetching all mobile app relationships from live API
func TestAccDatasourceMobileAppRelationship_01_All(t *testing.T) {
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
					check.That(dataSourceType + ".all").Key("filter_type").HasValue("all"),
					// Note: items.# may be 0 if no relationships exist in the tenant
				),
			},
		},
	})
}

// TestAccMobileAppRelationshipDataSource_BySourceId tests filtering relationships by source app ID from live API
func TestAccDatasourceMobileAppRelationship_02_BySourceId(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("02_by_source_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_source_id").Key("filter_type").HasValue("source_id"),
					check.That(dataSourceType+".by_source_id").Key("filter_value").HasValue("app-source-test-001"),
					// Note: items.# may be 0 if no relationships exist for this source_id
				),
			},
		},
	})
}

// TestAccMobileAppRelationshipDataSource_ODataFilter tests using OData filter queries from live API
func TestAccDatasourceMobileAppRelationship_03_ODataFilter(t *testing.T) {
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
					check.That(dataSourceType+".odata_filter").Key("filter_type").HasValue("odata"),
					check.That(dataSourceType+".odata_filter").Key("odata_filter").HasValue("targetType eq 'parent'"),
					// Note: items.# may be 0 if no relationships match the filter
				),
			},
		},
	})
}

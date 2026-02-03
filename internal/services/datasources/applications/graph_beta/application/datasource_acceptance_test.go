package graphBetaApplication_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceApplication_01_ByObjectId(t *testing.T) {
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
				Config: loadAccTestTerraform("01_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_object_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_object_id").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceApplication_02_ByDisplayName(t *testing.T) {
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
				Config: loadAccTestTerraform("02_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_display_name").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_display_name").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceApplication_03_ByAppId(t *testing.T) {
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
				Config: loadAccTestTerraform("03_by_app_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_app_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".by_app_id").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceApplication_04_ODataFilter(t *testing.T) {
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
				Config: loadAccTestTerraform("04_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_filter").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_filter").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceApplication_05_ODataAdvanced(t *testing.T) {
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
				Config: loadAccTestTerraform("05_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_advanced").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_advanced").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceApplication_06_ODataComprehensive(t *testing.T) {
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
				Config: loadAccTestTerraform("06_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_comprehensive").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("display_name").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("publisher_domain").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("tags.#").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("created_date_time").Exists(),
				),
			},
		},
	})
}

// Helper function to load acceptance test Terraform configs
func loadAccTestTerraform(filename string) string {
	accTestConfig, err := helpers.ParseHCLFile(fmt.Sprintf("tests/terraform/acceptance/%s", filename))
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

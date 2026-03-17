package graphBetaWindowsUpdateProduct_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateProduct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/product"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	dataSourceType = "data." + graphBetaWindowsUpdateProduct.DataSourceName
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccDatasourceWindowsUpdateProduct_01_ByCatalogId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(dataSourceType, "Searching for product by catalog ID")
				},
				Config: loadAcceptanceTestTerraform("01_by_catalog_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("search_type").HasValue("catalog_id"),
					check.That(dataSourceType+".test").Key("products.#").Exists(),
					check.That(dataSourceType+".test").Key("products.0.id").Exists(),
					check.That(dataSourceType+".test").Key("products.0.name").Exists(),
					check.That(dataSourceType+".test").Key("products.0.group_name").Exists(),
					check.That(dataSourceType+".test").Key("products.0.revisions.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsUpdateProduct_02_ByKbNumber(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(dataSourceType, "Searching for product by KB number")
				},
				Config: loadAcceptanceTestTerraform("02_by_kb_number.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("search_type").HasValue("kb_number"),
					check.That(dataSourceType+".test").Key("products.#").Exists(),
					check.That(dataSourceType+".test").Key("products.0.id").Exists(),
					check.That(dataSourceType+".test").Key("products.0.name").Exists(),
					check.That(dataSourceType+".test").Key("products.0.group_name").Exists(),
					check.That(dataSourceType+".test").Key("products.0.revisions.#").Exists(),
				),
			},
		},
	})
}

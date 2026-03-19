package graphBetaWindowsUpdateCatalog_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccDatasourceWindowsUpdateCatalog_01_AllEntries(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving all Windows Update catalog entries")
				},
				Config: loadAcceptanceTestTerraform("01_all_entries.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".test").Key("entries.#").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.id").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.display_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.release_date_time").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsUpdateCatalog_02_FeatureUpdatesOnly(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Filtering for feature updates only")
				},
				Config: loadAcceptanceTestTerraform("02_feature_updates_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("filter_type").HasValue("catalog_entry_type"),
					check.That(dataSourceType+".test").Key("filter_value").HasValue("featureUpdate"),
					check.That(dataSourceType+".test").Key("entries.#").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("featureUpdate"),
					check.That(dataSourceType+".test").Key("entries.0.version").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsUpdateCatalog_03_QualityUpdatesOnly(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Filtering for quality updates only")
				},
				Config: loadAcceptanceTestTerraform("03_quality_updates_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("filter_type").HasValue("catalog_entry_type"),
					check.That(dataSourceType+".test").Key("filter_value").HasValue("qualityUpdate"),
					check.That(dataSourceType+".test").Key("entries.#").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("qualityUpdate"),
					check.That(dataSourceType+".test").Key("entries.0.catalog_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.short_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.quality_update_classification").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsUpdateCatalog_04_ByDisplayName(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Searching by display name")
				},
				Config: loadAcceptanceTestTerraform("04_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("filter_type").HasValue("display_name"),
					check.That(dataSourceType+".test").Key("filter_value").HasValue("Windows 11"),
					check.That(dataSourceType+".test").Key("entries.#").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.display_name").Exists(),
				),
			},
		},
	})
}

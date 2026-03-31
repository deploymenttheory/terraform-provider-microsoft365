package graphBetaWindowsUpdatesApplicableContent_test

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

func TestAccDatasourceApplicableContent_01_Basic(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving all applicable content for deployment audience")
				},
				Config: loadAcceptanceTestTerraform("01_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").Exists(),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.#"),
				),
			},
		},
	})
}

func TestAccDatasourceApplicableContent_02_DriverUpdates(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving driver updates only")
				},
				Config: loadAcceptanceTestTerraform("02_driver_updates.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").Exists(),
					check.That(dataSourceType+".test").Key("catalog_entry_type").HasValue("driver"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.#"),
				),
			},
		},
	})
}

func TestAccDatasourceApplicableContent_03_DisplayDrivers(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving display driver updates only")
				},
				Config: loadAcceptanceTestTerraform("03_display_drivers.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").Exists(),
					check.That(dataSourceType+".test").Key("catalog_entry_type").HasValue("driver"),
					check.That(dataSourceType+".test").Key("driver_class").HasValue("Display"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.#"),
				),
			},
		},
	})
}

func TestAccDatasourceApplicableContent_04_WithMatchesOnly(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving content with matched devices only")
				},
				Config: loadAcceptanceTestTerraform("04_with_matches_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").Exists(),
					check.That(dataSourceType+".test").Key("include_no_matches").HasValue("false"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.#"),
				),
			},
		},
	})
}

func TestAccDatasourceApplicableContent_05_DeviceSpecific(t *testing.T) {
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
					testlog.StepAction(dataSourceType, "Retrieving applicable content for specific device")
				},
				Config: loadAcceptanceTestTerraform("05_device_specific.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("audience_id").Exists(),
					check.That(dataSourceType+".test").Key("device_id").Exists(),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "applicable_content.#"),
				),
			},
		},
	})
}

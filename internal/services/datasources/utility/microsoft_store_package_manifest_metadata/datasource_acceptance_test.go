package utilityMicrosoftStorePackageManifest_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMicrosoftStorePackageManifestDataSource_PackageIdentifier(t *testing.T) {
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
				Config: testAccConfigPackageIdentifier(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("package_identifier").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("id").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("manifests.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("manifests.0.package_identifier").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.type").IsSet(),
				),
			},
		},
	})
}

func TestAccMicrosoftStorePackageManifestDataSource_SearchTerm(t *testing.T) {
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
				Config: testAccConfigSearchTerm(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("search_term").HasValue("PC Manager"),
					check.That("data."+dataSourceType+".test").Key("id").HasValue("PC Manager"),
					check.That("data."+dataSourceType+".test").Key("manifests.#").IsSet(),
				),
			},
		},
	})
}

func TestAccMicrosoftStorePackageManifestDataSource_ValidateVersionStructure(t *testing.T) {
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
				Config: testAccConfigValidateStructure(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("package_identifier").HasValue("XP8M1ZJCZ99QJW"),
					check.That("data."+dataSourceType+".test").Key("id").HasValue("XP8M1ZJCZ99QJW"),
					check.That("data."+dataSourceType+".test").Key("manifests.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("manifests.0.package_identifier").HasValue("XP8M1ZJCZ99QJW"),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.0.package_version").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.0.default_locale.package_name").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.0.default_locale.publisher").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.0.installers.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.0.default_locale.tags.#").IsSet(),
				),
			},
		},
	})
}

func TestAccMicrosoftStorePackageManifestDataSource_MultipleResults(t *testing.T) {
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
				Config: testAccConfigMultipleResults(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("search_term").HasValue("Microsoft"),
					check.That("data."+dataSourceType+".test").Key("manifests.#").IsSet(),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigPackageIdentifier() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_package_identifier.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigSearchTerm() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_search_term.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigValidateStructure() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_validate_structure.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigMultipleResults() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_multiple_results.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

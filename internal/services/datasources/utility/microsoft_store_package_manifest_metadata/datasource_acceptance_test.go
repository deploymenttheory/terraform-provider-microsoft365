package utilityMicrosoftStorePackageManifest_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "package_identifier", "9PM860492SZD"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "id", "9PM860492SZD"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.package_identifier", "9PM860492SZD"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.type"),
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "search_term", "PC Manager"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "id", "PC Manager"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#"),
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "package_identifier", "XP8M1ZJCZ99QJW"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "id", "XP8M1ZJCZ99QJW"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.package_identifier", "XP8M1ZJCZ99QJW"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.0.package_version"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.0.default_locale.package_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.0.default_locale.publisher"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.0.installers.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.0.default_locale.tags.#"),
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
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "search_term", "Microsoft"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#"),
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


package utilityMicrosoftStorePackageManifest_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMicrosoftStorePackageManifestDataSource_PackageIdentifier(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigPackageIdentifier(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "package_identifier", "9PM860492SZD"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "id", "9PM860492SZD"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.package_identifier", "9PM860492SZD"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.0.versions.#"),
				),
			},
		},
	})
}

func TestMicrosoftStorePackageManifestDataSource_SearchTerm(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchTerm(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "search_term", "PC Manager"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_store_package_manifest_metadata.test", "manifests.#"),
				),
			},
		},
	})
}

func TestMicrosoftStorePackageManifestDataSource_BothParametersProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBothProvided(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination|conflicts with"),
			},
		},
	})
}

func TestMicrosoftStorePackageManifestDataSource_NeitherParameterProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigNeitherProvided(),
				ExpectError: regexp.MustCompile("Missing Input Parameter"),
			},
		},
	})
}

// Configuration functions
func testConfigPackageIdentifier() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_package_identifier.tf")
	if err != nil {
		panic("failed to load package_identifier config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSearchTerm() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_search_term.tf")
	if err != nil {
		panic("failed to load search_term config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigBothProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_both_provided.tf")
	if err != nil {
		panic("failed to load both_provided config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigNeitherProvided() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_neither_provided.tf")
	if err != nil {
		panic("failed to load neither_provided config: " + err.Error())
	}
	return unitTestConfig
}

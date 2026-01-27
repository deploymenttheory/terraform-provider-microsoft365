package utilityMicrosoftStorePackageManifest_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityMicrosoftStorePackageManifest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/microsoft_store_package_manifest_metadata"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityMicrosoftStorePackageManifest.DataSourceName
)

func TestUnitDatasourceMicrosoftStorePackageManifestMetadata_01_PackageIdentifier(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigPackageIdentifier(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("package_identifier").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("id").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("manifests.#").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.0.package_identifier").HasValue("9PM860492SZD"),
					check.That("data."+dataSourceType+".test").Key("manifests.0.versions.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceMicrosoftStorePackageManifestMetadata_02_SearchTerm(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSearchTerm(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("search_term").HasValue("PC Manager"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("manifests.#").IsSet(),
				),
			},
		},
	})
}

func TestUnitDatasourceMicrosoftStorePackageManifestMetadata_03_BothParametersProvided(t *testing.T) {
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

func TestUnitDatasourceMicrosoftStorePackageManifestMetadata_04_NeitherParameterProvided(t *testing.T) {
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

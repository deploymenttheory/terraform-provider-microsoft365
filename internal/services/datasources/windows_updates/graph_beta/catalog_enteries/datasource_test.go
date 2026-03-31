package graphBetaWindowsUpdateCatalog_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateCatalog "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/graph_beta/catalog_enteries"
	catalogMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/windows_updates/graph_beta/catalog_enteries/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + graphBetaWindowsUpdateCatalog.DataSourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *catalogMocks.WindowsUpdateCatalogMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catalogMock := &catalogMocks.WindowsUpdateCatalogMock{}
	catalogMock.RegisterMocks()
	return mockClient, catalogMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *catalogMocks.WindowsUpdateCatalogMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	catalogMock := &catalogMocks.WindowsUpdateCatalogMock{}
	catalogMock.RegisterErrorMocks()
	return mockClient, catalogMock
}

func TestUnitDatasourceWindowsUpdateCatalog_01_AllEntries(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_all_entries.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entries.#").HasValue("5"),
					check.That(dataSourceType+".test").Key("entries.0.id").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.display_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.release_date_time").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsUpdateCatalog_02_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entries.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("entries.0.id").HasValue("c1dec151-c151-c1de-51c1-dec151c1dec1"),
					check.That(dataSourceType+".test").Key("entries.0.display_name").HasValue("Windows 11, version 25H2"),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("featureUpdate"),
					check.That(dataSourceType+".test").Key("entries.0.version").HasValue("Windows 11, version 25H2"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsUpdateCatalog_03_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entries.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("entries.0.display_name").MatchesRegex(regexp.MustCompile("SecurityUpdate")),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("qualityUpdate"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsUpdateCatalog_04_FeatureUpdatesOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_feature_updates_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entries.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("featureUpdate"),
					check.That(dataSourceType+".test").Key("entries.0.version").Exists(),
					check.That(dataSourceType+".test").Key("entries.1.catalog_entry_type").HasValue("featureUpdate"),
					check.That(dataSourceType+".test").Key("entries.1.version").Exists(),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsUpdateCatalog_05_QualityUpdatesOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_quality_updates_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("entries.#").HasValue("3"),
					check.That(dataSourceType+".test").Key("entries.0.catalog_entry_type").HasValue("qualityUpdate"),
					check.That(dataSourceType+".test").Key("entries.0.catalog_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.short_name").Exists(),
					check.That(dataSourceType+".test").Key("entries.0.quality_update_classification").Exists(),

					// Verify CVE information for first entry (critical update with exploited CVEs)
					check.That(dataSourceType+".test").Key("entries.0.cve_severity_information.max_severity").HasValue("critical"),
					check.That(dataSourceType+".test").Key("entries.0.cve_severity_information.max_base_score").HasValue("9.8"),
					check.That(dataSourceType+".test").Key("entries.0.cve_severity_information.exploited_cves.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("entries.0.cve_severity_information.exploited_cves.0.number").HasValue("CVE-2026-12345"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsUpdateCatalog_06_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, catalogMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer catalogMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_all_entries.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|Insufficient privileges"),
			},
		},
	})
}

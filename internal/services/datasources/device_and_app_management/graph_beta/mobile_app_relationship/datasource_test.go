package graphBetaMobileAppRelationship_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMobileAppRelationship "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app_relationship"
	relMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app_relationship/mocks"
)

var dataSourceType = "data." + graphBetaMobileAppRelationship.DataSourceName

const (
	sourceID     = relMocks.SourceAppID
	dependencyID = sourceID + "_22222222-2222-2222-2222-222222222222"
	supersedesID = sourceID + "_33333333-3333-3333-3333-333333333333"
	parentID     = "44444444-4444-4444-4444-444444444444_" + sourceID
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment(t *testing.T, errorMocks bool) {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	relationshipMock := &relMocks.MobileAppRelationshipMock{}
	if errorMocks {
		relationshipMock.RegisterErrorMocks()
	} else {
		relationshipMock.RegisterMocks()
	}
	t.Cleanup(relationshipMock.CleanupMockState)
}

func runUnitTest(t *testing.T, filename string, checks ...resource.TestCheckFunc) {
	t.Helper()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: loadUnitTestTerraform(filename),
			Check:  resource.ComposeTestCheckFunc(checks...),
		}},
	})
}

func runUnitErrorTest(t *testing.T, filename string, expected *regexp.Regexp) {
	t.Helper()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      loadUnitTestTerraform(filename),
			ExpectError: expected,
		}},
	})
}

func TestUnitDatasourceMobileAppRelationship_01_ByRelationshipId(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "01_by_relationship_id.tf",
		d.Key("id").HasValue(dependencyID),
		d.Key("relationship_id").HasValue(dependencyID),
		d.Key("source_id").HasValue(sourceID),
		d.Key("source_display_name").HasValue("Contoso Line of Business App"),
		d.Key("source_display_version").HasValue("3.2.0"),
		d.Key("source_publisher_display_name").HasValue("Contoso"),
		d.Key("target_id").HasValue("22222222-2222-2222-2222-222222222222"),
		d.Key("target_display_name").HasValue("Microsoft Visual C++ Redistributable"),
		d.Key("target_display_version").HasValue("14.40.33810"),
		d.Key("target_publisher").HasValue("Microsoft Corporation"),
		d.Key("target_publisher_display_name").HasValue("Microsoft"),
		d.Key("target_type").HasValue("child"),
		d.Key("relationship_type").HasValue("dependency"),
		d.Key("dependency_type").HasValue("autoInstall"),
		d.Key("dependent_app_count").HasValue("1"),
		d.Key("depends_on_app_count").HasValue("2"),
		d.Key("supersedence_type").DoesNotExist(),
		d.Key("superseded_app_count").DoesNotExist(),
	)
}

func TestUnitDatasourceMobileAppRelationship_02_ByTargetId(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "02_by_target_id.tf",
		d.Key("id").HasValue(supersedesID),
		d.Key("target_id").HasValue("33333333-3333-3333-3333-333333333333"),
		d.Key("target_display_name").HasValue("Contoso Line of Business App (Legacy)"),
		d.Key("relationship_type").HasValue("supersedence"),
		d.Key("supersedence_type").HasValue("update"),
		d.Key("superseded_app_count").HasValue("1"),
		d.Key("superseding_app_count").HasValue("0"),
		d.Key("dependency_type").DoesNotExist(),
		d.Key("dependent_app_count").DoesNotExist(),
	)
}

func TestUnitDatasourceMobileAppRelationship_03_ByTargetDisplayName(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "03_by_target_display_name.tf",
		d.Key("id").HasValue(dependencyID),
		// The configured value is preserved even though the lookup is case-insensitive.
		d.Key("target_display_name").HasValue("microsoft visual c++ redistributable"),
		d.Key("relationship_type").HasValue("dependency"),
	)
}

func TestUnitDatasourceMobileAppRelationship_04_ByODataQuery(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "04_by_odata_query.tf",
		d.Key("odata_query").HasValue("targetType eq 'parent'"),
		d.Key("id").HasValue(parentID),
		d.Key("target_type").HasValue("parent"),
		d.Key("target_id").HasValue("44444444-4444-4444-4444-444444444444"),
		d.Key("supersedence_type").HasValue("replace"),
	)
}

func TestUnitDatasourceMobileAppRelationship_05_TargetTypeMismatch(t *testing.T) {
	setupMockEnvironment(t, false)
	runUnitErrorTest(t, "05_target_type_mismatch.tf",
		regexp.MustCompile(`no mobile app relationship found with target_id`))
}

func TestUnitDatasourceMobileAppRelationship_06_ODataQueryMultipleMatches(t *testing.T) {
	setupMockEnvironment(t, false)
	runUnitErrorTest(t, "06_odata_query_multiple_matches.tf",
		regexp.MustCompile(`must match exactly one mobile app relationship: found 2`))
}

func TestUnitDatasourceMobileAppRelationship_07_ODataQueryWithTargetType(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "07_odata_query_with_target_type.tf",
		d.Key("id").HasValue(supersedesID),
		d.Key("target_type").HasValue("child"),
	)
}

func TestUnitDatasourceMobileAppRelationship_08_ByOpaqueRelationshipId(t *testing.T) {
	setupMockEnvironment(t, false)
	d := check.That(dataSourceType + ".test")
	runUnitTest(t, "08_by_opaque_relationship_id.tf",
		d.Key("id").HasValue("legacy-relationship-001"),
		d.Key("source_id").HasValue(sourceID),
		d.Key("dependency_type").HasValue("detect"),
	)
}

func TestUnitDatasourceMobileAppRelationship_09_MissingSourceId(t *testing.T) {
	setupMockEnvironment(t, false)
	runUnitErrorTest(
		t,
		"09_missing_source_id.tf",
		regexp.MustCompile(
			`Attribute "source_id" must be specified when "target_id" is\s+specified`,
		),
	)
}

func TestUnitDatasourceMobileAppRelationship_10_ConflictingLookups(t *testing.T) {
	setupMockEnvironment(t, false)
	runUnitErrorTest(t, "10_conflicting_lookups.tf",
		regexp.MustCompile(`Invalid Attribute Combination`))
}

func TestUnitDatasourceMobileAppRelationship_11_NoLookup(t *testing.T) {
	setupMockEnvironment(t, false)
	runUnitErrorTest(t, "11_no_lookup.tf",
		regexp.MustCompile(`Invalid Attribute Combination`))
}

func TestUnitDatasourceMobileAppRelationship_12_GraphError(t *testing.T) {
	setupMockEnvironment(t, true)
	runUnitErrorTest(t, "02_by_target_id.tf",
		regexp.MustCompile(`(?i)forbidden|insufficient privileges`))
}

package graphBetaMobileAppRelationship_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	win32 "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/win32_app"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccDatasourceMobileAppRelationship_01_Lookups builds real relationships and
// resolves them by relationship_id, target_id, target_display_name and odata_query.
func TestAccDatasourceMobileAppRelationship_01_Lookups(t *testing.T) {
	const (
		dependency   = "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency.test"
		supersedence = "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			t.Setenv("TF_VAR_installer_zip", prepareAcceptanceInstaller(t))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: win32.ResourceName,
				TestResource: win32.Win32AppTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(dataSourceType, "01: Creating Win32 app relationships and looking them up")
				},
				Config: loadAcceptanceTestTerraform("01_lookups.tf"),
				Check: resource.ComposeTestCheckFunc(
					// relationship_id: the source app is derived from the {source}_{target} ID.
					check.That(dataSourceType+".by_relationship_id").Key("id").
						MatchesOtherKey(check.That(dependency).Key("id")),
					check.That(dataSourceType+".by_relationship_id").Key("source_id").
						MatchesOtherKey(check.That(win32.ResourceName+".main").Key("id")),
					check.That(dataSourceType+".by_relationship_id").Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+".runtime").Key("id")),
					check.That(dataSourceType+".by_relationship_id").Key("target_display_name").
						MatchesOtherKey(check.That(win32.ResourceName+".runtime").Key("display_name")),
					check.That(dataSourceType+".by_relationship_id").Key("source_display_name").
						MatchesOtherKey(check.That(win32.ResourceName+".main").Key("display_name")),
					check.That(dataSourceType+".by_relationship_id").Key("relationship_type").HasValue("dependency"),
					check.That(dataSourceType+".by_relationship_id").Key("dependency_type").HasValue("autoInstall"),
					check.That(dataSourceType+".by_relationship_id").Key("target_type").HasValue("child"),

					// target_id: from the main app's side the supersedence is incoming.
					check.That(dataSourceType+".by_target_id").Key("id").Exists(),
					check.That(dataSourceType+".by_target_id").Key("relationship_type").HasValue("supersedence"),
					check.That(dataSourceType+".by_target_id").Key("supersedence_type").HasValue("update"),
					check.That(dataSourceType+".by_target_id").Key("target_type").HasValue("parent"),

					// target_display_name
					check.That(dataSourceType+".by_target_display_name").Key("id").
						MatchesOtherKey(check.That(dependency).Key("id")),
					check.That(dataSourceType+".by_target_display_name").Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+".runtime").Key("id")),

					// odata_query: the main app has two relationships, so this proves Graph applied $filter.
					check.That(dataSourceType+".by_odata_query").Key("id").
						MatchesOtherKey(check.That(dependency).Key("id")),
					check.That(dataSourceType+".by_odata_query").Key("odata_query").
						HasValue("isof('microsoft.graph.mobileAppDependency')"),
					check.That(dataSourceType+".by_odata_query").Key("relationship_type").HasValue("dependency"),

					// odata_query from the dependent app's side returns the incoming relationship.
					check.That(dataSourceType+".parent_by_odata_query").Key("target_type").HasValue("parent"),
					check.That(dataSourceType+".parent_by_odata_query").Key("relationship_type").HasValue("dependency"),
					check.That(dataSourceType+".parent_by_odata_query").Key("id").Exists(),
				),
			},
		},
	})
}

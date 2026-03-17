package graphBetaWindowsUpdatesAutopatchContentApproval_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsAutopatchContentApproval "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_content_approval"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceType = graphBetaWindowsAutopatchContentApproval.ResourceName
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsUpdateContentApproval_01_FeatureUpdateApproval(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create content approval for feature update")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_approval.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("update_policy_id").Exists(),
					check.That(resourceType+".test").Key("catalog_entry_id").Exists(),
					check.That(resourceType+".test").Key("catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("is_revoked").HasValue("false"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("deployment_settings.schedule.start_date_time").HasValue("2026-04-01T00:00:00Z"),
					check.That(resourceType+".test").Key("deployment_settings.schedule.gradual_rollout.end_date_time").HasValue("2026-04-15T00:00:00Z"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import content approval")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceWindowsUpdateContentApproval_02_RequiresImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create initial content approval")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_approval.tf"),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Attempt to create duplicate content approval")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_approval.tf") + `
resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "import" {
  update_policy_id    = "45a01ef3-fb4b-8c1d-2428-1f060464033c"
  catalog_entry_id    = "c1dec151-c151-c1de-51c1-dec151c1dec1"
  catalog_entry_type  = "featureUpdate"
}
`,
				ExpectError: regexp.MustCompile(`already exists|AlreadyExists`),
			},
		},
	})
}

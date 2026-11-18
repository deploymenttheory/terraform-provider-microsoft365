package graphBetaGroup_test

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestAccResourceGroup_GeneratedImportCreationDefaults(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_import_defaults.tf")
	updated := strings.Replace(config, "Group import default preservation", "Group import default preservation updated", 1)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 20*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				func(_ *terraform.State) error {
					testlog.WaitForConsistency("group", 20*time.Second)
					time.Sleep(20 * time.Second)
					return nil
				},
				check.That(resourceType+".test").ExistsInGraph(testResource),
				check.That(resourceType+".test").Key("visibility").HasValue("Private"),
				check.That(resourceType+".test").Key("is_assignable_to_role").HasValue("false"),
			)},
			{
				ResourceName: resourceType + ".test",
				ImportState:  true, GenerateConfig: true, ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources[resourceType+".test"].Primary.ID + ":hard_delete=true", nil
				},
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".generated", plancheck.ResourceActionNoop),
					},
				},
			},
			{Config: config, PlanOnly: true},
			{Config: updated, Check: check.That(resourceType + ".test").Key("description").HasValue("Group import default preservation updated")},
			{Config: updated, PlanOnly: true},
		},
	})
}

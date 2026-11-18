package graphBetaSettingsCatalogConfigurationPolicyJson_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// The profile is assigned only to a newly created, empty security group.
func TestAccAssignmentGeneratedImport(t *testing.T) {
	config := loadAcceptanceTestTerraform("resource_assignment_import.tf")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName: resourceType + ".test",
				ImportState:  true, GenerateConfig: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".generated", plancheck.ResourceActionNoop),
					},
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

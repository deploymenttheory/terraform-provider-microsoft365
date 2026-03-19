package graphBetaWindowsUpdatesAutopatchContentApproval_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsAutopatchContentApproval "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_content_approval"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
}

// TestAccResourceWindowsUpdateContentApproval_01_FeatureUpdateApproval creates a content approval
// for a feature update using a dynamically resolved catalog entry, then verifies import.
func TestAccResourceWindowsUpdateContentApproval_01_FeatureUpdateApproval(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
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
				ResourceName: resourceType + ".test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test")
					}
					policyId := rs.Primary.Attributes["update_policy_id"]
					id := rs.Primary.ID
					return policyId + "/" + id, nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsUpdateContentApproval_02_UpdateDeploymentSettings creates a content approval
// then updates the deployment schedule, verifying the diff-based update.
func TestAccResourceWindowsUpdateContentApproval_02_UpdateDeploymentSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create content approval with initial schedule")
				},
				Config: loadAcceptanceTestTerraform("02_update_deployment_settings_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("deployment_settings.schedule.start_date_time").HasValue("2026-05-01T00:00:00Z"),
					check.That(resourceType+".test").Key("deployment_settings.schedule.gradual_rollout.end_date_time").HasValue("2026-05-15T00:00:00Z"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update content approval deployment schedule")
				},
				Config: loadAcceptanceTestTerraform("02_update_deployment_settings_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("deployment_settings.schedule.start_date_time").HasValue("2026-06-01T00:00:00Z"),
					check.That(resourceType+".test").Key("deployment_settings.schedule.gradual_rollout.end_date_time").HasValue("2026-06-30T00:00:00Z"),
				),
			},
		},
	})
}

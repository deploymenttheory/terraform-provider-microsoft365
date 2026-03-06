package graphBetaWindowsAutopatchDeployment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsAutopatchDeploymentResource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopatch_deployment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testResource = WindowsAutopatchDeploymentResource.WindowsUpdateDeploymentTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsUpdateDeployment_01_FeatureUpdateDeployment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment for feature update")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P7D"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("100"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.0.signal").HasValue("rollback"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.0.threshold").HasValue("5"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.0.action").HasValue("pauseDeployment"),
					check.That(resourceType+".test").Key("state.effective_value").Exists(),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import deployment")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_02_UpdateDeploymentState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Pause deployment")
					testlog.WaitForConsistency(fmt.Sprintf("%s (pause)", resourceType), 10*time.Second)
					time.Sleep(10 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("02_deployment_paused.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("state.requested_value").HasValue("paused"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

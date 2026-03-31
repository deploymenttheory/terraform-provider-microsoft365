package graphBetaWindowsUpdatesAutopatchDeploymentState_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsUpdatesAutopatchDeploymentStateResource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_deployment_state"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testResource = WindowsUpdatesAutopatchDeploymentStateResource.WindowsUpdatesAutopatchDeploymentStateTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsUpdateDeploymentState_01_PauseAndUnpause(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment and pause it")
				},
				Config: loadAcceptanceTestTerraform("01_pause_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("deployment_id").Exists(),
					check.That(resourceType+".test").Key("requested_value").HasValue("paused"),
					check.That(resourceType+".test").Key("effective_value").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Unpause deployment")
					testlog.WaitForConsistency(fmt.Sprintf("%s (unpause)", resourceType), 10*time.Second)
					time.Sleep(10 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("02_unpause_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("requested_value").HasValue("none"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

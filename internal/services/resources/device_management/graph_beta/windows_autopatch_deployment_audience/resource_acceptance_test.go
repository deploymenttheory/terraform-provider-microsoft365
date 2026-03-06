package graphBetaWindowsAutopatchDeploymentAudience_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdateDeploymentAudience "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_deployment_audience"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper function to load acceptance test configs
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

const resourceType = graphBetaWindowsUpdateDeploymentAudience.ResourceName

var testResource = graphBetaWindowsUpdateDeploymentAudience.WindowsUpdateDeploymentAudienceTestResource{}

// Test 001: Basic audience creation
func TestAccResourceWindowsUpdateDeploymentAudience_01_BasicAudience(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaWindowsUpdateDeploymentAudience.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_basic_audience.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

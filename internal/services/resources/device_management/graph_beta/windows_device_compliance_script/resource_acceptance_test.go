package graphBetaWindowsDeviceComplianceScript_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsDeviceComplianceScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_script"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaWindowsDeviceComplianceScript.ResourceName

var testResource = graphBetaWindowsDeviceComplianceScript.DeviceComplianceScriptTestResource{}

func TestAccResourceWindowsDeviceComplianceScript_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigLifecycleCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "display_name", "Acceptance - Windows Device Compliance Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "enforce_signature_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "run_as_32_bit", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "id"),
				),
			},
			{
				Config: testAccConfigLifecycleUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "display_name", "Acceptance - Windows Device Compliance Script - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle", "run_as_32_bit", "true"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_windows_device_compliance_script.lifecycle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccConfigLifecycleCreate() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_create.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigLifecycleUpdate() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_update.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

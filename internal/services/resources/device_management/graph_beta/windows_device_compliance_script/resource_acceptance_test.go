package graphBetaWindowsDeviceComplianceScript_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceWindowsDeviceComplianceScript_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceComplianceScriptDestroy,
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

func testAccCheckWindowsDeviceComplianceScriptDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_device_compliance_script" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			DeviceComplianceScripts().
			ByDeviceComplianceScriptId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows device compliance script %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows device compliance script %s still exists", rs.Primary.ID)
	}
	return nil
}

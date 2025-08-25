package graphBetaAppControlForBusinessManagedInstaller_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppControlForBusinessManagedInstallerResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckAppControlForBusinessManagedInstallerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "intune_management_extension_as_managed_installer", "Disabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "available_version", "1.93.102.0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "managed_installer_configured_date_time", ""),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled",
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

func TestAccAppControlForBusinessManagedInstallerResource_Enabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckAppControlForBusinessManagedInstallerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigEnabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "intune_management_extension_as_managed_installer", "Enabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "available_version", "1.93.102.0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "managed_installer_configured_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "id"),
				),
			},
		},
	})
}

func TestAccAppControlForBusinessManagedInstallerResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckAppControlForBusinessManagedInstallerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "intune_management_extension_as_managed_installer", "Disabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.disabled", "managed_installer_configured_date_time", ""),
				),
			},
			{
				Config: testAccConfigEnabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "intune_management_extension_as_managed_installer", "Enabled"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.enabled", "managed_installer_configured_date_time"),
				),
			},
		},
	})
}

func testAccConfigDisabled() string {
	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/managed_installer_disabled.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigEnabled() string {
	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/managed_installer_enabled.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccCheckAppControlForBusinessManagedInstallerDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_app_control_for_business_managed_installer" {
			continue
		}
		windowsManagementApp, err := graphClient.
			DeviceAppManagement().
			WindowsManagementApp().
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Windows Management App not found (404/NotFound) - this is expected\n")
				continue
			}
			return fmt.Errorf("error checking Windows Management App status: %v", err)
		}

		// Check if managed installer is disabled (which is our "destroyed" state)
		managedInstaller := windowsManagementApp.GetManagedInstaller()
		if managedInstaller != nil && managedInstaller.String() == "enabled" {
			return fmt.Errorf("Windows Management App managed installer is still enabled for resource %s", rs.Primary.ID)
		}

		fmt.Printf("DEBUG: Windows Management App managed installer is disabled (destroyed state) for resource %s\n", rs.Primary.ID)
	}
	return nil
}

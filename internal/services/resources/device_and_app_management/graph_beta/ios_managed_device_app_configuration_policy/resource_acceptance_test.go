package graphBetaIOSManagedDeviceAppConfigurationPolicy_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceIOSManagedDeviceAppConfigurationPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIOSManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			// Create with custom settings configuration
			{
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "display_name", regexp.MustCompile(`^acc-test-ios-managed-device-app-configuration-policy-custom-settings-[a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "settings.#", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to xml encoded configuration
			{
				Config: loadAcceptanceTestTerraform("resource_xml_config.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "display_name", regexp.MustCompile(`^acc-test-ios-managed-device-app-configuration-policy-xml-encoded-[a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "encoded_setting_xml"),
				),
			},
		},
	})
}

func TestAccResourceIOSManagedDeviceAppConfigurationPolicy_02_CustomSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIOSManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "display_name", regexp.MustCompile(`^acc-test-ios-managed-device-app-configuration-policy-custom-settings-[a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.custom_settings", "settings.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceIOSManagedDeviceAppConfigurationPolicy_03_XMLEncoded(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIOSManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_xml_config.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "id"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "display_name", regexp.MustCompile(`^acc-test-ios-managed-device-app-configuration-policy-xml-encoded-[a-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.xml_encoded", "encoded_setting_xml"),
				),
			},
		},
	})
}

func testAccCheckIOSManagedDeviceAppConfigurationPolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" {
			continue
		}

		_, err := graphClient.
			DeviceAppManagement().
			MobileAppConfigurations().
			ByManagedDeviceMobileAppConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if iOS managed device app configuration policy %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("iOS managed device app configuration policy %s still exists", rs.Primary.ID)
	}

	return nil
}

package graphBetaIOSMobileAppConfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIOSMobileAppConfigurationResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckIOSMobileAppConfigurationDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccIOSMobileAppConfigurationConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Test Acceptance iOS Mobile App Configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "version"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccIOSMobileAppConfigurationConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Test Acceptance iOS Mobile App Configuration - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "targeted_mobile_apps.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "settings.#", "2"),
				),
			},
		},
	})
}

func TestAccIOSMobileAppConfigurationResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckIOSMobileAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIOSMobileAppConfigurationConfig_description(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Test iOS Mobile App Configuration with Description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "description", "Test description for iOS mobile app configuration"),
				),
			},
		},
	})
}

func TestAccIOSMobileAppConfigurationResource_RoleScopeTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckIOSMobileAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIOSMobileAppConfigurationConfig_roleScopeTags(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Test iOS Mobile App Configuration with Role Scope Tags"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration.test", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

func testAccCheckIOSMobileAppConfigurationDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" {
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
			return fmt.Errorf("error checking if iOS mobile app configuration %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("iOS mobile app configuration %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccIOSMobileAppConfigurationConfig_minimal() string {
	return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test Acceptance iOS Mobile App Configuration"
}
`
}

func testAccIOSMobileAppConfigurationConfig_maximal() string {
	return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name        = "Test Acceptance iOS Mobile App Configuration - Updated"
  description         = "Updated description for acceptance testing"
  targeted_mobile_apps = ["12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321"]
  role_scope_tag_ids  = ["0", "1"]
  
  encoded_setting_xml = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+dGVzdEtleTwva2V5PgoJPHN0cmluZz50ZXN0VmFsdWU8L3N0cmluZz4KPC9kaWN0Pgo8L3BsaXN0Pgo="
  
  settings = [
    {
      app_config_key       = "testKey1"
      app_config_key_type  = "stringType"
      app_config_key_value = "testValue1"
    },
    {
      app_config_key       = "testKey2"
      app_config_key_type  = "integerType"
      app_config_key_value = "123"
    }
  ]
}
`
}

func testAccIOSMobileAppConfigurationConfig_description() string {
	return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Mobile App Configuration with Description"
  description  = "Test description for iOS mobile app configuration"
}
`
}

func testAccIOSMobileAppConfigurationConfig_roleScopeTags() string {
	return `
resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name       = "Test iOS Mobile App Configuration with Role Scope Tags"
  role_scope_tag_ids = ["0", "1", "2"]
}
`
}
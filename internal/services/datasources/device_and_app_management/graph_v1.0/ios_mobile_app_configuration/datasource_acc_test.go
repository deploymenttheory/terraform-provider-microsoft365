package graphIOSMobileAppConfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestAccIOSMobileAppConfigurationDataSource_Read(t *testing.T) {

	resourceType := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration"
	resourceName := "test"
	dataSourceName := "test_ds"
	resourceID := fmt.Sprintf("%s.%s", resourceType, resourceName)
	dataSourceID := fmt.Sprintf("data.%s.%s", resourceType, dataSourceName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create resource first
			{
				Config: testAccIOSMobileAppConfigurationResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "display_name", "Test iOS App Config"),
					resource.TestCheckResourceAttr(resourceID, "description", "Test iOS app configuration for acceptance testing"),
					resource.TestCheckResourceAttrSet(resourceID, "id"),
				),
			},
			// Test data source by ID
			{
				Config: testAccIOSMobileAppConfigurationDataSourceConfigByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceID, "id", resourceID, "id"),
					resource.TestCheckResourceAttrPair(dataSourceID, "display_name", resourceID, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceID, "description", resourceID, "description"),
					resource.TestCheckResourceAttrPair(dataSourceID, "targeted_mobile_apps.#", resourceID, "targeted_mobile_apps.#"),
					resource.TestCheckResourceAttrPair(dataSourceID, "settings.#", resourceID, "settings.#"),
					resource.TestCheckResourceAttrSet(dataSourceID, "created_date_time"),
					resource.TestCheckResourceAttrSet(dataSourceID, "last_modified_date_time"),
					resource.TestCheckResourceAttrSet(dataSourceID, "version"),
				),
			},
			// Test data source by display name
			{
				Config: testAccIOSMobileAppConfigurationDataSourceConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceID, "id", resourceID, "id"),
					resource.TestCheckResourceAttrPair(dataSourceID, "display_name", resourceID, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceID, "description", resourceID, "description"),
				),
			},
		},
	})
}

func TestAccIOSMobileAppConfigurationDataSource_ReadWithSettings(t *testing.T) {

	resourceType := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration"
	resourceName := "test"
	dataSourceName := "test_ds"
	resourceID := fmt.Sprintf("%s.%s", resourceType, resourceName)
	dataSourceID := fmt.Sprintf("data.%s.%s", resourceType, dataSourceName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create resource with settings
			{
				Config: testAccIOSMobileAppConfigurationResourceConfigWithSettings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "display_name", "Test iOS Config with Settings"),
					resource.TestCheckResourceAttr(resourceID, "settings.#", "2"),
					resource.TestCheckResourceAttrSet(resourceID, "id"),
				),
			},
			// Test data source
			{
				Config: testAccIOSMobileAppConfigurationDataSourceConfigWithSettings(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceID, "id", resourceID, "id"),
					resource.TestCheckResourceAttrPair(dataSourceID, "display_name", resourceID, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceID, "settings.#", resourceID, "settings.#"),
					resource.TestCheckResourceAttrPair(dataSourceID, "settings.0.app_config_key", resourceID, "settings.0.app_config_key"),
					resource.TestCheckResourceAttrPair(dataSourceID, "settings.0.app_config_key_type", resourceID, "settings.0.app_config_key_type"),
					resource.TestCheckResourceAttrPair(dataSourceID, "settings.0.app_config_key_value", resourceID, "settings.0.app_config_key_value"),
				),
			},
		},
	})
}

// Helper functions

func testAccIOSMobileAppConfigurationResourceConfig() string {
	return `
provider "microsoft365" {
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS App Config"
  description  = "Test iOS app configuration for acceptance testing"
  targeted_mobile_apps = ["com.example.testapp"]
}
`
}

func testAccIOSMobileAppConfigurationDataSourceConfigByID() string {
	return `
provider "microsoft365" {
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS App Config"
  description  = "Test iOS app configuration for acceptance testing"
  targeted_mobile_apps = ["com.example.testapp"]
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test_ds" {
  id = microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test.id
}
`
}

func testAccIOSMobileAppConfigurationDataSourceConfigByDisplayName() string {
	return `
provider "microsoft365" {
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS App Config"
  description  = "Test iOS app configuration for acceptance testing"
  targeted_mobile_apps = ["com.example.testapp"]
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test_ds" {
  display_name = microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test.display_name
}
`
}

func testAccIOSMobileAppConfigurationResourceConfigWithSettings() string {
	return `
provider "microsoft365" {
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Config with Settings"
  description  = "Test iOS app configuration with settings"
  targeted_mobile_apps = ["com.example.testapp"]
  
  settings {
    app_config_key       = "serverUrl"
    app_config_key_type  = "stringType"
    app_config_key_value = "https://api.example.com"
  }
  
  settings {
    app_config_key       = "syncInterval"
    app_config_key_type  = "integerType"
    app_config_key_value = "300"
  }
}
`
}

func testAccIOSMobileAppConfigurationDataSourceConfigWithSettings() string {
	return `
provider "microsoft365" {
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Test iOS Config with Settings"
  description  = "Test iOS app configuration with settings"
  targeted_mobile_apps = ["com.example.testapp"]
  
  settings {
    app_config_key       = "serverUrl"
    app_config_key_type  = "stringType"
    app_config_key_value = "https://api.example.com"
  }
  
  settings {
    app_config_key       = "syncInterval"
    app_config_key_type  = "integerType"
    app_config_key_value = "300"
  }
}

data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test_ds" {
  id = microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test.id
}
`
}
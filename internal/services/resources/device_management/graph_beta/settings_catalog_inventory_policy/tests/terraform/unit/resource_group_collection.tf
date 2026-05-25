resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "group_collection" {
  name               = "Test Inventory Policy Group Collection - Unit"
  description        = "Test inventory policy with full application inventory settings"
  platforms          = "windows10"
  technologies       = "extensibility"
  role_scope_tag_ids = ["0"]

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey"
          group_setting_collection_value = [
            {
              children = [
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_appname"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_appname_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_appversion"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_appversion_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_architectures"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_architectures_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_publisher"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_publisher_24"
                  }
                }
              ]
            }
          ]
        }
        id = "0"
      }
    ]
  }
}

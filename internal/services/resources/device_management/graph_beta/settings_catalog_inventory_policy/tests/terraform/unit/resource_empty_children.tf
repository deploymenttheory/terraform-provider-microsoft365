resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "empty_children" {
  name               = "Test Inventory Policy Empty Children - Unit"
  description        = "Test inventory policy for settings without children"
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
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_appversion"
                  choice_setting_collection_value = [
                    {
                      children = []
                      value    = "windows_applicationinventory_applicationproperties_applicationkey_appversion_24"
                    }
                  ]
                }
              ]
            },
            {
              children = []
            }
          ]
        }
        id = "0"
      },
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id = "windows_applicationinventory_collectioninterval"
          choice_setting_value = {
            children = []
            value    = "windows_applicationinventory_collectioninterval_24"
          }
        }
        id = "1"
      }
    ]
  }
}

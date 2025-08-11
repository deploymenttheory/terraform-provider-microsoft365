resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "complex_group_collection" {
  name         = "Test Complex Group Collection Setting - Unit"
  description  = "Testing complex group collection with nested simple collection from System Preferences"
  platforms    = "macOS"
  technologies = ["mdm"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "com.apple.systempreferences_com.apple.systempreferences"
          group_setting_collection_value = [
            {
              children = [
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
                  setting_definition_id = "com.apple.systempreferences_disabledpreferencepanes"
                  simple_setting_collection_value = [
                    {
                      odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                      value      = "com.apple.AirDrop-Handoff-Settings.extension"
                    },
                    {
                      odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                      value      = "com.apple.Family-Settings.extension"
                    },
                    {
                      odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                      value      = "com.apple.Game-Center-Settings.extension"
                    }
                  ]
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
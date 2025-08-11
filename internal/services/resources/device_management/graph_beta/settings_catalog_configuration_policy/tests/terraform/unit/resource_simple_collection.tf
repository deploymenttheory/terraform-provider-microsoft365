resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "simple_collection" {
  name         = "Test Simple Collection Setting - Unit"
  description  = "Testing simple collection setting type from Edge extensions"
  platforms    = "macOS"
  technologies = ["mdm"]
  
  template_reference = {
    template_id = ""
  }
  
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          setting_definition_id = "com.apple.managedclient.preferences_extensioninstallallowlist"
          simple_setting_collection_value = [
            {
              odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value      = "odfafepnkmbhccpbejgmiehpchacaeak"
            },
            {
              odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value      = "nkbndigcebkoaejohleckhekfmcecfja"
            }
          ]
        }
        id = "0"
      }
    ]
  }
}
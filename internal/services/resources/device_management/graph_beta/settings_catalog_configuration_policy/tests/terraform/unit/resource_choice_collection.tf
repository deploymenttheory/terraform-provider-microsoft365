resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "choice_collection" {
  name         = "Test Choice Collection Setting - Unit"
  description  = "Testing choice collection setting type"
  platforms    = "macOS"
  technologies = ["mdm"]
  
  template_reference = {
    template_id = ""
  }
  
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
          setting_definition_id = "test.choice.collection"
          choice_setting_collection_value = [
            {
              value    = "choice_value_option_1"
              children = []
            },
            {
              value    = "choice_value_option_2"
              children = []
            }
          ]
        }
        id = "0"
      }
    ]
  }
}
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "simple_collection" {
  name         = "Test Simple Collection Setting - JSON Unit"
  description  = "Test configuration for simple collection setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
          "settingDefinitionId" = "test.simple.collection.setting",
          "simpleSettingCollectionValue" = [
            {
              "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value"       = "collection_value_1"
            },
            {
              "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value"       = "collection_value_2"
            }
          ]
        }
      }
    ]
  })
}
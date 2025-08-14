resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "choice_collection" {
  name         = "Test Choice Collection Setting - JSON Unit"
  description  = "Test configuration for choice collection setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
          "settingDefinitionId" = "test.choice.collection.setting",
          "choiceSettingCollectionValue" = [
            {
              "children" = [],
              "value"    = "choice_option_1"
            },
            {
              "children" = [],
              "value"    = "choice_option_2"
            }
          ]
        }
      }
    ]
  })
}
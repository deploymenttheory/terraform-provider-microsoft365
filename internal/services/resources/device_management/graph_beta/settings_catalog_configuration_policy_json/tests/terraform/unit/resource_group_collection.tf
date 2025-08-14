resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "group_collection" {
  name         = "Test Group Collection Setting - JSON Unit"
  description  = "Test configuration for group collection setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" = "test.group.collection.setting",
          "groupSettingCollectionValue" = [
            {
              "children" = [
                {
                  "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId" = "test.group.child.string",
                  "simpleSettingValue" = {
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value"       = "child_string_value"
                  }
                },
                {
                  "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId" = "test.group.child.integer",
                  "simpleSettingValue" = {
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "value"       = 42
                  }
                },
                {
                  "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId" = "test.group.child.choice",
                  "choiceSettingValue" = {
                    "children" = [],
                    "value"    = "test_choice_value"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  })
}
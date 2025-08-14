resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "complex_group_collection" {
  name         = "Test Complex Group Collection Setting - JSON Unit"
  description  = "Test configuration for complex group collection setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" = "test.complex.group.collection.setting",
          "groupSettingCollectionValue" = [
            {
              "children" = [
                {
                  "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
                  "settingDefinitionId" = "test.nested.simple.collection",
                  "simpleSettingCollectionValue" = [
                    {
                      "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                      "value"       = "nested_value_1"
                    },
                    {
                      "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                      "value"       = "nested_value_2"
                    },
                    {
                      "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                      "value"       = "nested_value_3"
                    }
                  ]
                }
              ]
            }
          ]
        }
      }
    ]
  })
}
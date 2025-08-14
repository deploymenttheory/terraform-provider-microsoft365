resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "simple_string" {
  name         = "Test Simple String Setting - JSON Unit"
  description  = "Test configuration for simple string setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.simple.string.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "test_string_value"
          }
        }
      }
    ]
  })
}
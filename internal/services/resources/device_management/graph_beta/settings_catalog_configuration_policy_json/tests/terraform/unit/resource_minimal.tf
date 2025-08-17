resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "minimal" {
  name      = "Test Minimal Settings Catalog Policy - JSON Unit"
  platforms = "macOS"

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.minimal.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "minimal_value"
          }
        }
      }
    ]
  })
}
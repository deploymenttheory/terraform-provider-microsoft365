resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "choice" {
  name         = "Test Choice Setting - JSON Unit"
  description  = "Test configuration for choice setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" = "test.choice.setting",
          "choiceSettingValue" = {
            "children" = [],
            "value"    = "com.apple.managedclient.preferences_smartscreenenabled_true"
          }
        }
      }
    ]
  })
}
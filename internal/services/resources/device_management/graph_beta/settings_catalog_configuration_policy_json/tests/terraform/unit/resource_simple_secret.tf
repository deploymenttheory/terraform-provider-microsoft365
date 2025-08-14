resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "simple_secret" {
  name         = "Test Simple Secret Setting - JSON Unit"
  description  = "Test configuration for simple secret setting using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.simple.secret.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
            "value"       = "test_secret_value",
            "valueState"  = "notEncrypted"
          }
        }
      }
    ]
  })
}
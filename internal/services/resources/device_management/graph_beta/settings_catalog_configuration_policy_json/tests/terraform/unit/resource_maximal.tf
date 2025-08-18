resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "maximal" {
  name               = "Test Maximal Settings Catalog Policy - JSON Unit"
  description        = "Comprehensive test configuration with maximal settings using JSON format"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0", "1"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId"              = "test.maximal.group.setting",
          "settingInstanceTemplateReference" = null,
          "groupSettingCollectionValue" = [
            {
              "settingValueTemplateReference" = null,
              "children" = [
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "test.maximal.string",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = "maximal_string_value"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "test.maximal.integer",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = 999
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "test.maximal.choice",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "test.maximal.choice.option1"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "test.maximal.secret",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "settingValueTemplateReference" = null,
                    "valueState"                    = "notEncrypted",
                    "value"                         = "maximal_secret_value"
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
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name               = "Test Acceptance Settings Catalog Policy JSON - Updated"
  description        = "Updated description for acceptance testing - JSON"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId"              = "com.apple.mcx.filevault2_com.apple.mcx.filevault2",
          "settingInstanceTemplateReference" = null,
          "groupSettingCollectionValue" = [
            {
              "settingValueTemplateReference" = null,
              "children" = [
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_defer",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_defer_true"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_enable",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_enable_1"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_outputpath",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = "/updated/output/path"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  })

  timeouts = {
    create = "900s"
    read   = "900s"
    update = "900s"
    delete = "900s"
  }
}
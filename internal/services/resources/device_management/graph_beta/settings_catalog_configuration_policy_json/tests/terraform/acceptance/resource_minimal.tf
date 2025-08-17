resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "macos_mdm_filevault2_settings" {
  name               = "macos mdm filevault2 settings - JSON"
  description        = "Configure the FileVault payload to manage FileVault disk encryption settings on devices - JSON format."
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

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
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_deferdontaskatuserlogout",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_deferdontaskatuserlogout_false"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_deferforceatuserloginmaxbypassattempts",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = 0
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_enable",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_enable_0"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_forceenableinsetupassistant",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_forceenableinsetupassistant_false"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_outputpath",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = "/output/path"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_password",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "settingValueTemplateReference" = null,
                    "valueState"                    = "notEncrypted",
                    "value"                         = "3669d68b-ea40-4682-abc9-9445f3f6fc7e"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_recoverykeyrotationinmonths",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_recoverykeyrotationinmonths_10"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_showrecoverykey",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_showrecoverykey_true"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_usekeychain",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_usekeychain_false"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_userecoverykey",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_userecoverykey_true"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_userentersmissinginfo",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_userentersmissinginfo_false"
                  }
                },
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_username",
                  "settingInstanceTemplateReference" = null,
                  "simpleSettingValue" = {
                    "@odata.type"                   = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" = null,
                    "value"                         = "username"
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
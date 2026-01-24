resource "random_string" "office_macos_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "office_configuration_macos" {
  name               = "acc-test-14-office-macos-${random_string.office_macos_suffix.result}"
  description        = "Acceptance test policy for Office Configuration on macOS with nested group collections"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          settingDefinitionId              = "com.apple.servicemanagement_com.apple.servicemanagement"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          groupSettingCollectionValue = [
            {
              settingValueTemplateReference = null
              children = [
                {
                  "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
                  settingDefinitionId              = "com.apple.servicemanagement_rules"
                  settingInstanceTemplateReference = null
                  auditRuleInformation             = null
                  groupSettingCollectionValue = [
                    {
                      settingValueTemplateReference = null
                      children = [
                        {
                          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                          settingDefinitionId              = "com.apple.servicemanagement_rules_item_comment"
                          settingInstanceTemplateReference = null
                          auditRuleInformation             = null
                          simpleSettingValue = {
                            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                            value                         = "Office Licensing Helper"
                            settingValueTemplateReference = null
                          }
                        },
                        {
                          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                          settingDefinitionId              = "com.apple.servicemanagement_rules_item_ruletype"
                          settingInstanceTemplateReference = null
                          auditRuleInformation             = null
                          choiceSettingValue = {
                            value                         = "com.apple.servicemanagement_rules_item_ruletype_0"
                            settingValueTemplateReference = null
                            children                      = []
                          }
                        },
                        {
                          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                          settingDefinitionId              = "com.apple.servicemanagement_rules_item_rulevalue"
                          settingInstanceTemplateReference = null
                          auditRuleInformation             = null
                          simpleSettingValue = {
                            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                            value                         = "com.microsoft.office.licensingV2.helper"
                            settingValueTemplateReference = null
                          }
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "com.apple.managedclient.preferences_officeautosignin"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "com.apple.managedclient.preferences_officeautosignin_true"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId              = "com.apple.managedclient.preferences_officeactivationemailaddress"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          simpleSettingValue = {
            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value                         = "{{userprincipalname}}"
            settingValueTemplateReference = null
          }
        }
      }
    ]
  })
}

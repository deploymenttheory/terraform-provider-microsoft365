resource "random_string" "defender_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "defender_antivirus_baseline" {
  name               = "acc-test-15-defender-baseline-${random_string.defender_suffix.result}"
  description        = "Acceptance test policy for Defender Antivirus Security Baseline with complex nesting"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_disableblockatfirstseen"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_disableblockatfirstseen_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_realtimeprotection_disableioavprotection"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_realtimeprotection_disableioavprotection_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_disablerealtimemonitoring"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_disablerealtimemonitoring_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_realtimeprotection_disablebehaviormonitoring"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_realtimeprotection_disablebehaviormonitoring_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "4"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_scan_disableremovabledrivescanning"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_scan_disableremovabledrivescanning_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "5"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_defender_cloudblocklevel"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_defender_cloudblocklevel_2"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "6"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_defender_submitsamplesconsent"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_defender_submitsamplesconsent_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "7"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_spynetreporting"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_spynetreporting_1"
            settingValueTemplateReference    = null
            children = [
              {
                "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_spynetreporting_spynetreporting"
                settingInstanceTemplateReference     = null
                auditRuleInformation                 = null
                choiceSettingValue = {
                  value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_spynetreporting_spynetreporting_2"
                  settingValueTemplateReference    = null
                  children                         = []
                }
              }
            ]
          }
        }
      },
      {
        id = "8"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_1"
            settingValueTemplateReference    = null
            children = [
              {
                "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
                settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules"
                settingInstanceTemplateReference     = null
                auditRuleInformation                 = null
                groupSettingCollectionValue = [
                  {
                    settingValueTemplateReference = null
                    children = [
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_key"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "75668c1f-73b5-4cf0-bb93-3ecf5cb7cc84"
                          settingValueTemplateReference     = null
                        }
                      },
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_value"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "1"
                          settingValueTemplateReference     = null
                        }
                      }
                    ]
                  },
                  {
                    settingValueTemplateReference = null
                    children = [
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_key"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "3b576869-a4ec-4529-8536-b80a7769e899"
                          settingValueTemplateReference     = null
                        }
                      },
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_value"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "1"
                          settingValueTemplateReference     = null
                        }
                      }
                    ]
                  },
                  {
                    settingValueTemplateReference = null
                    children = [
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_key"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "d4f940ab-401b-4efc-aadc-ad5f3c50688a"
                          settingValueTemplateReference     = null
                        }
                      },
                      {
                        "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                        settingDefinitionId                  = "device_vendor_msft_policy_config_admx_microsoftdefenderantivirus_exploitguard_asr_rules_exploitguard_asr_rules_value"
                      settingInstanceTemplateReference     = null
                      auditRuleInformation                 = null
                      simpleSettingValue = {
                          "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                          value                             = "1"
                          settingValueTemplateReference     = null
                        }
                      }
                    ]
                  }
                ]
              }
            ]
          }
        }
      }
    ]
  })
}

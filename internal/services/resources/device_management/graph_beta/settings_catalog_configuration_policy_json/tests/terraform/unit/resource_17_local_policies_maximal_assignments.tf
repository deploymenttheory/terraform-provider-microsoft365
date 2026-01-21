resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "local_policies_maximal" {
  name               = "Test Local Policies Security Options Policy"
  description        = "Test policy for Local Policies Security Options with maximal assignments"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_blockmicrosoftaccounts"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_blockmicrosoftaccounts_3"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_enableadministratoraccountstatus"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_enableadministratoraccountstatus_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_renameadministratoraccount"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value                             = "DTAdmin"
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_renameguestaccount"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value                             = "DTGuest"
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "4"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_interactivelogon_machineinactivitylimit"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                             = 900
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "5"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_localpoliciessecurityoptions_interactivelogon_messagetextforusersattemptingtologon"
          settingInstanceTemplateReference     = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                             = "Unauthorized access is prohibited"
              settingValueTemplateReference     = null
            }
          ]
        }
      }
    ]
  })

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "22222222-2222-2222-2222-222222222222"
      filter_type = "none"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_type = "none"
    }
  ]
}

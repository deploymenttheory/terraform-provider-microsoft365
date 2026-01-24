resource "random_string" "local_policies_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_string" "local_policies_group1_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_string" "local_policies_group2_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create test groups for assignments
resource "microsoft365_graph_beta_groups_group" "local_policies_test_group1" {
  display_name     = "acc-test-17-group1-${random_string.local_policies_group1_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-17-group1-${random_string.local_policies_group1_suffix.result}"
  security_enabled = true
  description      = "Test group 1 for local policies assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "local_policies_test_group2" {
  display_name     = "acc-test-17-group2-${random_string.local_policies_group2_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-17-group2-${random_string.local_policies_group2_suffix.result}"
  security_enabled = true
  description      = "Test group 2 for local policies assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_30_seconds_local_policies" {
  create_duration = "30s"
  depends_on = [
    microsoft365_graph_beta_groups_group.local_policies_test_group1,
    microsoft365_graph_beta_groups_group.local_policies_test_group2
  ]
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "local_policies_maximal" {
  name               = "acc-test-17-local-policies-${random_string.local_policies_suffix.result}"
  description        = "Acceptance test policy for Local Policies Security Options with maximal assignments"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_blockmicrosoftaccounts"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_blockmicrosoftaccounts_3"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_enableadministratoraccountstatus"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_enableadministratoraccountstatus_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_renameadministratoraccount"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          simpleSettingValue = {
            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value                         = "DTAdmin"
            settingValueTemplateReference = null
          }
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_accounts_renameguestaccount"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          simpleSettingValue = {
            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value                         = "DTGuest"
            settingValueTemplateReference = null
          }
        }
      },
      {
        id = "4"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_interactivelogon_machineinactivitylimit"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          simpleSettingValue = {
            "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                         = 900
            settingValueTemplateReference = null
          }
        }
      },
      {
        id = "5"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_localpoliciessecurityoptions_interactivelogon_messagetextforusersattemptingtologon"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                         = "Unauthorized access is prohibited"
              settingValueTemplateReference = null
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
      group_id    = microsoft365_graph_beta_groups_group.local_policies_test_group1.id
      filter_type = "none"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.local_policies_test_group2.id
      filter_type = "none"
    }
  ]

  depends_on = [time_sleep.wait_30_seconds_local_policies]
}

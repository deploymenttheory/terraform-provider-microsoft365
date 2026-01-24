resource "random_string" "attachment_manager_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "attachment_manager" {
  name               = "acc-test-06-attachment-manager-${random_string.attachment_manager_suffix.result}"
  description        = "Acceptance test policy for attachment manager settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "user_vendor_msft_policy_config_attachmentmanager_notifyantivirus"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "user_vendor_msft_policy_config_attachmentmanager_notifyantivirus_0"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "user_vendor_msft_policy_config_attachmentmanager_hidezoneinfoonproperties"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "user_vendor_msft_policy_config_attachmentmanager_hidezoneinfoonproperties_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })
}

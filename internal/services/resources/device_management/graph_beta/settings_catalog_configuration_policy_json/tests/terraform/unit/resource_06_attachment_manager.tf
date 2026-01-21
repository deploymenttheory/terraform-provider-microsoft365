resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "attachment_manager" {
  name               = "Test Attachment Manager Policy"
  description        = "Test policy for attachment manager settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "user_vendor_msft_policy_config_attachmentmanager_notifyantivirus"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "user_vendor_msft_policy_config_attachmentmanager_notifyantivirus_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "user_vendor_msft_policy_config_attachmentmanager_hidezoneinfoonproperties"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "user_vendor_msft_policy_config_attachmentmanager_hidezoneinfoonproperties_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      }
    ]
  })
}

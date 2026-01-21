resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "credential_user_interface" {
  name               = "Test Credential User Interface Policy"
  description        = "Test policy for credential user interface settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_credentialui_enumerateadministrators"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_credentialui_enumerateadministrators_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "user_vendor_msft_policy_config_credentialsui_disablepasswordreveal"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "user_vendor_msft_policy_config_credentialsui_disablepasswordreveal_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      }
    ]
  })
}

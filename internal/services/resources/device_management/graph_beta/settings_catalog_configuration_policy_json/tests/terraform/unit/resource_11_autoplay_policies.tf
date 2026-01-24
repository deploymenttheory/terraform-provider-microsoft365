resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "autoplay" {
  name               = "Test AutoPlay Policies Policy"
  description        = "Test policy for AutoPlay with nested choice settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_autoplay_disallowautoplayfornonvolumedevices"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_autoplay_disallowautoplayfornonvolumedevices_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_autoplay_setdefaultautorunbehavior"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_autoplay_setdefaultautorunbehavior_1"
            settingValueTemplateReference = null
            children = [
              {
                "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId              = "device_vendor_msft_policy_config_autoplay_setdefaultautorunbehavior_noautorun_dropdown"
                settingInstanceTemplateReference = null
                choiceSettingValue = {
                  value                         = "device_vendor_msft_policy_config_autoplay_setdefaultautorunbehavior_noautorun_dropdown_1"
                  settingValueTemplateReference = null
                  children                      = []
                }
              }
            ]
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_autoplay_turnoffautoplay"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_autoplay_turnoffautoplay_1"
            settingValueTemplateReference = null
            children = [
              {
                "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId              = "device_vendor_msft_policy_config_autoplay_turnoffautoplay_autorun_box"
                settingInstanceTemplateReference = null
                choiceSettingValue = {
                  value                         = "device_vendor_msft_policy_config_autoplay_turnoffautoplay_autorun_box_255"
                  settingValueTemplateReference = null
                  children                      = []
                }
              }
            ]
          }
        }
      }
    ]
  })
}

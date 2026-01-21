resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "defender_smartscreen" {
  name               = "Test Defender Smartscreen Policy"
  description        = "Test policy for Defender Smartscreen with choice and collection children"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~smartscreen_smartscreenenabled"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~smartscreen_smartscreenenabled_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_microsoft_edgev80diff~policy~microsoft_edge~smartscreen_smartscreenpuaenabled"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_microsoft_edgev80diff~policy~microsoft_edge~smartscreen_smartscreenpuaenabled_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~smartscreen_smartscreenallowlistdomains"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~smartscreen_smartscreenallowlistdomains_1"
            settingValueTemplateReference    = null
            children = [
              {
                "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
                settingDefinitionId                  = "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~smartscreen_smartscreenallowlistdomains_smartscreenallowlistdomainsdesc"
                settingInstanceTemplateReference     = null
                simpleSettingCollectionValue = [
                  {
                    "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    value                             = "deploymenttheory.com"
                    settingValueTemplateReference     = null
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

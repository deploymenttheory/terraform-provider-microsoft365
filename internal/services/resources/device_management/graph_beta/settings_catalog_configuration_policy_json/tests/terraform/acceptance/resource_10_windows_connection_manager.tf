resource "random_string" "wcm_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "windows_connection_manager" {
  name               = "acc-test-10-wcm-${random_string.wcm_suffix.result}"
  description        = "Acceptance test policy for Windows Connection Manager with choice children"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_wcm_wcm_disablepowermanagement"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_wcm_wcm_disablepowermanagement_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_wcm_wcm_enablesoftdisconnect"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_wcm_wcm_enablesoftdisconnect_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_admx_wcm_wcm_minimizeconnections"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_admx_wcm_wcm_minimizeconnections_1"
            settingValueTemplateReference    = null
            children = [
              {
                "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId                  = "device_vendor_msft_policy_config_admx_wcm_wcm_minimizeconnections_wcm_minimizeconnections_options"
                settingInstanceTemplateReference     = null
                auditRuleInformation                 = null
                choiceSettingValue = {
                  value                            = "device_vendor_msft_policy_config_admx_wcm_wcm_minimizeconnections_wcm_minimizeconnections_options_3"
                  settingValueTemplateReference    = null
                  children                         = []
                }
              }
            ]
          }
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_windowsconnectionmanager_prohitconnectiontonondomainnetworkswhenconnectedtodomainauthenticatednetwork"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_windowsconnectionmanager_prohitconnectiontonondomainnetworkswhenconnectedtodomainauthenticatednetwork_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      }
    ]
  })
}

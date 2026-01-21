resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "storage_sense" {
  name               = "Test Storage Sense Policy"
  description        = "Test policy for storage sense configuration"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_allowdiskhealthmodelupdates"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_storage_allowdiskhealthmodelupdates_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_allowstoragesenseglobal"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_storage_allowstoragesenseglobal_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_allowstoragesensetemporaryfilescleanup"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_storage_allowstoragesensetemporaryfilescleanup_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_configstoragesensecloudcontentdehydrationthreshold"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                             = 30
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "4"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_configstoragesensedownloadscleanupthreshold"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                             = 90
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "5"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_configstoragesenseglobalcadence"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                             = 0
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "6"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_configstoragesenserecyclebincleanupthreshold"
          settingInstanceTemplateReference     = null
          simpleSettingValue = {
            "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value                             = 30
            settingValueTemplateReference     = null
          }
        }
      },
      {
        id = "7"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_storage_removablediskdenywriteaccess"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_storage_removablediskdenywriteaccess_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      }
    ]
  })
}

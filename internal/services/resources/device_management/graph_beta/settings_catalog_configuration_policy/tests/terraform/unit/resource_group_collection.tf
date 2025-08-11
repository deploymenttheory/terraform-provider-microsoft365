resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "group_collection" {
  name         = "Test Group Collection Setting - Unit"
  description  = "Testing group collection setting type from FileVault configuration"
  platforms    = "macOS"
  technologies = ["mdm"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "com.apple.mcx.filevault2_com.apple.mcx.filevault2"
          group_setting_collection_value = [
            {
              children = [
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "com.apple.mcx.filevault2_enable"
                  choice_setting_value = {
                    children = []
                    value    = "com.apple.mcx.filevault2_enable_0"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "com.apple.mcx.filevault2_forceenableinsetupassistant"
                  choice_setting_value = {
                    children = []
                    value    = "com.apple.mcx.filevault2_forceenableinsetupassistant_true"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "com.apple.mcx.filevault2_recoverykeyrotationinmonths"
                  choice_setting_value = {
                    children = []
                    value    = "com.apple.mcx.filevault2_recoverykeyrotationinmonths_6"
                  }
                }
              ]
            }
          ]
        }
        id = "0"
      }
    ]
  }
}
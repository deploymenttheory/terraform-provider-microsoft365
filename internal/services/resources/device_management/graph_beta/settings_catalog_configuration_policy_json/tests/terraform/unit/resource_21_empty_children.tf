resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "empty_children" {
  name               = "Test Empty Children Policy"
  description        = "Test policy for settings without children"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
          settingDefinitionId              = "device_vendor_msft_defender_configuration_passiveremediation"
          settingInstanceTemplateReference = null
          choiceSettingCollectionValue = [
            {
              value                         = "device_vendor_msft_defender_configuration_passiveremediation_1"
              settingValueTemplateReference = null
              children                      = []
            }
          ]
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_test_groupcollection"
          settingInstanceTemplateReference = null
          groupSettingCollectionValue = [
            {
              settingValueTemplateReference = null
              children                      = []
            }
          ]
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_camera_allowcamera"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_camera_allowcamera_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_type = "none"
    }
  ]
}

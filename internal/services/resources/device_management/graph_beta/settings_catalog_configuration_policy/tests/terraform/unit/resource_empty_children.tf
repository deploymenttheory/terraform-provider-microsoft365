resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "empty_children" {
  name         = "Test Empty Children Policy - Unit"
  description  = "Test policy for settings without children"
  platforms    = "windows10"
  technologies = ["mdm"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
          setting_definition_id = "device_vendor_msft_defender_configuration_passiveremediation"
          choice_setting_collection_value = [
            {
              value    = "device_vendor_msft_defender_configuration_passiveremediation_1"
              children = []
            }
          ]
        }
        id = "0"
      },
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "device_vendor_msft_policy_config_test_groupcollection"
          group_setting_collection_value = [
            {
              children = []
            }
          ]
        }
        id = "1"
      },
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id = "device_vendor_msft_policy_config_camera_allowcamera"
          choice_setting_value = {
            value    = "device_vendor_msft_policy_config_camera_allowcamera_1"
            children = []
          }
        }
        id = "2"
      },
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "device_vendor_msft_policy_config_test_nestedgroup"
          group_setting_collection_value = [
            {
              children = [
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "device_vendor_msft_policy_config_test_nestedgroup_choice"
                  choice_setting_value = {
                    value    = "device_vendor_msft_policy_config_test_nestedgroup_choice_1"
                    children = []
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
                  setting_definition_id = "device_vendor_msft_policy_config_test_nestedgroup_choicecollection"
                  choice_setting_collection_value = [
                    {
                      value    = "device_vendor_msft_policy_config_test_nestedgroup_choicecollection_1"
                      children = []
                    }
                  ]
                }
              ]
            }
          ]
        }
        id = "3"
      }
    ]
  }
}

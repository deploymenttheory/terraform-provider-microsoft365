resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all_devices_assignment" {
  name               = "Test All Devices Assignment Settings Catalog Policy - Unique"
  description        = ""
  platforms          = "macOS"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "test.setting"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "value"
          }
        }
        id = "0"
      }
    ]
  }

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}


resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "maximal" {
  name               = "Test Maximal Settings Catalog Policy - Unique"
  description        = "Maximal settings catalog policy for testing with all features"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0", "1"]

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "com.apple.test.sample_string"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "hello"
          }
        }
        id = "0"
      },
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
          setting_definition_id = "com.apple.test.sample_int"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            value      = "5"
          }
        }
        id = "1"
      }
    ]
  }

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "44444444-4444-4444-4444-444444444444"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "include"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "66666666-6666-6666-6666-666666666666"
      filter_type = "exclude"
    },
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


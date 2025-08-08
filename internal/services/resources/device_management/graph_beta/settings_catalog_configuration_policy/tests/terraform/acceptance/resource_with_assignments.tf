resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test_assignments" {
  name               = "Test Settings Catalog Policy with Assignments"
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
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    }
  ]
}


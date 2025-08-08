resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all_assignment_types" {
  name               = "Test All Assignment Types Settings Catalog Policy - Unique"
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
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    }
  ]
}


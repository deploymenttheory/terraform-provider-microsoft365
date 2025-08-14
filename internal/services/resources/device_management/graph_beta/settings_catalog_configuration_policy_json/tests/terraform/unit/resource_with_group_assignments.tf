resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "group_assignments" {
  name         = "Test Group Assignments - JSON Unit"
  description  = "Test configuration with group assignments using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.group.assignment.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "group_assignment_test_value"
          }
        }
      }
    ]
  })

  assignments = [
    {
      type      = "groupAssignmentTarget"
      group_id  = "test-group-1-id"
      device_id = ""
    },
    {
      type      = "groupAssignmentTarget"
      group_id  = "test-group-2-id"
      device_id = ""
    }
  ]
}
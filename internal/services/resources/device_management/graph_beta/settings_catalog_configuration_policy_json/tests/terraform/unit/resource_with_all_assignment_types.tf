resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "all_assignment_types" {
  name         = "Test All Assignment Types - JSON Unit"
  description  = "Test configuration with all assignment types using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.assignment.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "assignment_test_value"
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
      type      = "allLicensedUsersAssignmentTarget"
      group_id  = ""
      device_id = ""
    },
    {
      type      = "allDevicesAssignmentTarget"
      group_id  = ""
      device_id = ""
    },
    {
      type      = "exclusionGroupAssignmentTarget"
      group_id  = "test-group-exclusion-id"
      device_id = ""
    }
  ]
}
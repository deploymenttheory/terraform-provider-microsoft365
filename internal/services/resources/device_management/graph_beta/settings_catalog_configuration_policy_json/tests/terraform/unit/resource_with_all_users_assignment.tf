resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "all_users_assignment" {
  name         = "Test All Users Assignment - JSON Unit"
  description  = "Test configuration with all users assignment using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.all.users.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "all_users_test_value"
          }
        }
      }
    ]
  })

  assignments = [
    {
      type      = "allLicensedUsersAssignmentTarget"
      group_id  = ""
      device_id = ""
    }
  ]
}
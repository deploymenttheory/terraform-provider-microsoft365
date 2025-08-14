resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "exclusion_assignment" {
  name         = "Test Exclusion Assignment - JSON Unit"
  description  = "Test configuration with exclusion assignment using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.exclusion.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "exclusion_test_value"
          }
        }
      }
    ]
  })

  assignments = [
    {
      type      = "exclusionGroupAssignmentTarget"
      group_id  = "test-exclusion-group-id"
      device_id = ""
    }
  ]
}
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "all_devices_assignment" {
  name         = "Test All Devices Assignment - JSON Unit"
  description  = "Test configuration with all devices assignment using JSON format"
  platforms    = "macOS"
  technologies = ["mdm"]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"         = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" = "test.all.devices.setting",
          "simpleSettingValue" = {
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value"       = "all_devices_test_value"
          }
        }
      }
    ]
  })

  assignments = [
    {
      type      = "allDevicesAssignmentTarget"
      group_id  = ""
      device_id = ""
    }
  ]
}
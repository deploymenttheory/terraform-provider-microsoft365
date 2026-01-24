resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "task_manager" {
  name               = "Test Task Manager Policy"
  description        = "Test policy for task manager settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_taskmanager_allowendtask"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_taskmanager_allowendtask_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })
}

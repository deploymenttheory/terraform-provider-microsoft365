resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "file_explorer_minimal" {
  name               = "Test File Explorer Policy"
  description        = "Test policy for File Explorer with minimal assignments"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_fileexplorer_turnoffdataexecutionpreventionforexplorer"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_fileexplorer_turnoffdataexecutionpreventionforexplorer_0"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_fileexplorer_turnoffheapterminationoncorruption"
          settingInstanceTemplateReference = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_fileexplorer_turnoffheapterminationoncorruption_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_type = "none"
    }
  ]
}

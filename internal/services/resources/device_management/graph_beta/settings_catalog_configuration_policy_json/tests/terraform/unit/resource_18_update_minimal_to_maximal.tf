resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "update_test" {
  name               = "Test Update Policy"
  description        = "Test policy for update from minimal to maximal"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_fileexplorer_turnoffdataexecutionpreventionforexplorer"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_fileexplorer_turnoffdataexecutionpreventionforexplorer_0"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_fileexplorer_turnoffheapterminationoncorruption"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_fileexplorer_turnoffheapterminationoncorruption_1"
            settingValueTemplateReference    = null
            children                         = []
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

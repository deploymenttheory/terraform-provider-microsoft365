resource "random_string" "assignment_import" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "assignment_import" {
  display_name     = "acc-test-4873-catalog-group-${random_string.assignment_import.result}"
  mail_nickname    = "acc-test-4873-catalog-${random_string.assignment_import.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name               = "acc-test-4873-catalog-${random_string.assignment_import.result}"
  description        = "Test policy for File Explorer with minimal assignments"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          auditRuleInformation             = null
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
          auditRuleInformation             = null
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
      group_id    = microsoft365_graph_beta_groups_group.assignment_import.id
      filter_type = "none"
    }
  ]
}

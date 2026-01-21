resource "random_string" "file_explorer_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_string" "file_explorer_group_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create test group for assignment
resource "microsoft365_graph_beta_groups_group" "file_explorer_test_group" {
  display_name     = "acc-test-16-group-${random_string.file_explorer_group_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-16-group-${random_string.file_explorer_group_suffix.result}"
  security_enabled = true
  description      = "Test group for file explorer policy assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_30_seconds_file_explorer" {
  create_duration = "30s"
  depends_on      = [microsoft365_graph_beta_groups_group.file_explorer_test_group]
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "file_explorer_minimal" {
  name               = "acc-test-16-file-explorer-${random_string.file_explorer_suffix.result}"
  description        = "Acceptance test policy for File Explorer with minimal assignments"
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
          auditRuleInformation                 = null
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
          auditRuleInformation                 = null
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
      group_id    = microsoft365_graph_beta_groups_group.file_explorer_test_group.id
      filter_type = "none"
    }
  ]

  depends_on = [time_sleep.wait_30_seconds_file_explorer]
}

resource "random_string" "passive_remediation_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create test groups for assignment
resource "microsoft365_graph_beta_groups_group" "passive_remediation_group_1" {
  display_name     = "acc-test-18-group-1-${random_string.passive_remediation_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-18-group-1-${random_string.passive_remediation_suffix.result}"
  security_enabled = true
  description      = "Test group for passive remediation policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "passive_remediation_group_2" {
  display_name     = "acc-test-18-group-2-${random_string.passive_remediation_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-18-group-2-${random_string.passive_remediation_suffix.result}"
  security_enabled = true
  description      = "Test group for passive remediation policy assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_30_seconds_passive_remediation" {
  create_duration = "30s"
  depends_on = [
    microsoft365_graph_beta_groups_group.passive_remediation_group_1,
    microsoft365_graph_beta_groups_group.passive_remediation_group_2,
  ]
}

# Passive Remediation is a choice setting collection without children (#4956)
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "passive_remediation" {
  name               = "acc-test-18-passive-remediation-${random_string.passive_remediation_suffix.result}"
  description        = "Acceptance test policy for Defender passive remediation - updated"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
          settingDefinitionId              = "device_vendor_msft_defender_configuration_passiveremediation"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingCollectionValue = [
            {
              value                         = "device_vendor_msft_defender_configuration_passiveremediation_1"
              settingValueTemplateReference = null
              children                      = []
            }
          ]
        }
      }
    ]
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.passive_remediation_group_2.id
      filter_type = "none"
    }
  ]

  depends_on = [time_sleep.wait_30_seconds_passive_remediation]
}

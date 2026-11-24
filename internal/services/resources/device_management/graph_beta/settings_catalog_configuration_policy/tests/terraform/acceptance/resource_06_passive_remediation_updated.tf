resource "random_string" "passive_remediation_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "passive_remediation_group_1" {
  display_name     = "acc-test-06-group-1-${random_string.passive_remediation_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-06-group-1-${random_string.passive_remediation_suffix.result}"
  security_enabled = true
  description      = "Test group for passive remediation policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "passive_remediation_group_2" {
  display_name     = "acc-test-06-group-2-${random_string.passive_remediation_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-06-group-2-${random_string.passive_remediation_suffix.result}"
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

# Passive Remediation is a choice setting collection without children and Allow Archive Scanning is a choice setting
# without children (#4956). Both must be sent as "children": [].
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "passive_remediation" {
  name               = "acc-test-06-passive-remediation-${random_string.passive_remediation_suffix.result}"
  description        = "Acceptance test policy for Defender passive remediation - updated"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        id = "0"
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
          setting_definition_id = "device_vendor_msft_defender_configuration_passiveremediation"
          choice_setting_collection_value = [
            {
              value    = "device_vendor_msft_defender_configuration_passiveremediation_1"
              children = []
            }
          ]
        }
      },
      {
        id = "1"
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id = "device_vendor_msft_policy_config_defender_allowarchivescanning"
          choice_setting_value = {
            value    = "device_vendor_msft_policy_config_defender_allowarchivescanning_1"
            children = []
          }
        }
      }
    ]
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.passive_remediation_group_2.id
    }
  ]

  depends_on = [time_sleep.wait_30_seconds_passive_remediation]
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "test" {
  display_name               = "acceptance-test-macos-3974-${random_string.suffix.result}"
  depends_on                 = [time_sleep.group_replication]
  storage_require_encryption = true
  os_minimum_version         = "26.0"

  assignments = [{
    type        = "groupAssignmentTarget"
    group_id    = microsoft365_graph_beta_groups_group.test.id
    filter_type = "none"
  }]

  scheduled_actions_for_rule = [{
    scheduled_action_configurations = [{
      action_type        = "block"
      grace_period_hours = 72
    }]
  }]
}

# Keep the empty test group throughout the lifecycle, including assignment removal.
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-macos-3974-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-macos-3974-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "group_replication" {
  create_duration = "30s"
  depends_on      = [microsoft365_graph_beta_groups_group.test]
}

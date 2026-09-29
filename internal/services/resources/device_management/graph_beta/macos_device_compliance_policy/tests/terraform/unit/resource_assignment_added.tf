resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "test" {
  display_name               = "unit-test-macos-3974"
  storage_require_encryption = true
  os_minimum_version         = "26.0"

  assignments = [{
    type        = "groupAssignmentTarget"
    group_id    = "22222222-2222-2222-2222-222222222222"
    filter_type = "none"
  }]

  scheduled_actions_for_rule = [{
    scheduled_action_configurations = [{
      action_type        = "block"
      grace_period_hours = 72
    }]
  }]
}

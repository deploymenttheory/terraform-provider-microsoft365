resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "test" {
  display_name               = "acceptance-test-macos-3974-${random_string.suffix.result}"
  storage_require_encryption = true
  os_minimum_version         = "26.0"

  scheduled_actions_for_rule = [{
    rule_name = "PasswordRequired"
    scheduled_action_configurations = [{
      action_type        = "block"
      grace_period_hours = 72
    }]
  }]
}

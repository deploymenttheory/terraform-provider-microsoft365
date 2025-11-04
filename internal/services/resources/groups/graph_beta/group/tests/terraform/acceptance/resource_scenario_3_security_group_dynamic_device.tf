resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_3" {
  display_name                     = "acc-security-group-dynamic-device-${random_string.group_suffix.result}"
  mail_enabled                     = false
  mail_nickname                    = "accsg3${random_string.group_suffix.result}"
  security_enabled                 = true
  description                      = "Acceptance test - Security group with dynamic device membership"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(device.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}


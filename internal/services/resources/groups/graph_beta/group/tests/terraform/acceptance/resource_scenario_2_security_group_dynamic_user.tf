resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_2" {
  display_name                     = "acc-security-group-dynamic-user-${random_string.group_suffix.result}"
  mail_enabled                     = false
  mail_nickname                    = "accsg2${random_string.group_suffix.result}"
  security_enabled                 = true
  description                      = "Acceptance test - Security group with dynamic user membership"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}


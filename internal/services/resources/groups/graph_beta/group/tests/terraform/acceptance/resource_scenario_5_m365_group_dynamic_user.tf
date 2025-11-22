resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_5" {
  display_name                     = "acc-m365-group-dynamic-user-${random_string.group_suffix.result}"
  description                      = "Acceptance test - M365 group with dynamic user membership"
  mail_enabled                     = true
  mail_nickname                    = "accm365g5${random_string.group_suffix.result}"
  security_enabled                 = true
  group_types                      = ["Unified", "DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  visibility                       = "Private"
}


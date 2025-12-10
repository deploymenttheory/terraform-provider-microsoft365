resource "microsoft365_graph_beta_groups_group" "scenario_2" {
  display_name                     = "acc-security-group-with-dynamic-user-membership-type"
  mail_enabled                     = false
  mail_nickname                    = "f9a72987-7"
  security_enabled                 = true
  description                      = "acc-security-group-with-dynamic-user-membership-type"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  hard_delete                      = true
}


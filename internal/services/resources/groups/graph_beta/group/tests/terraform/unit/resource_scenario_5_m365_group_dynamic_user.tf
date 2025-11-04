resource "microsoft365_graph_beta_groups_group" "scenario_5" {
  display_name                    = "acc-m365-group-with-dynamic-user-membership-type"
  mail_enabled                    = true
  mail_nickname                   = "some-string"
  security_enabled                = true
  group_types                     = ["Unified", "DynamicMembership"]
  membership_rule                 = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  visibility                      = "Private"
}


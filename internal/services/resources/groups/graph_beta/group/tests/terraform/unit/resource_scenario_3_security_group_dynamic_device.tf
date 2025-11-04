resource "microsoft365_graph_beta_groups_group" "scenario_3" {
  display_name                     = "acc-security-group-with-dynamic-device-membership-type"
  mail_enabled                     = false
  mail_nickname                    = "17bf0e02-0"
  security_enabled                 = true
  description                      = "acc-security-group-with-dynamic-device-membership-type"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(device.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}


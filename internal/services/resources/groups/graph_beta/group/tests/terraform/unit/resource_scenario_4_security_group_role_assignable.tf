resource "microsoft365_graph_beta_groups_group" "scenario_4" {
  display_name          = "acc-security-group-with-entra-role-assignment"
  mail_enabled          = false
  mail_nickname         = "dec34327-9"
  security_enabled      = true
  description           = "acc-security-group-with-entra-role-assignment"
  is_assignable_to_role = true
}


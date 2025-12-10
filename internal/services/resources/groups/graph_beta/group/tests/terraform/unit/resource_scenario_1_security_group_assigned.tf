resource "microsoft365_graph_beta_groups_group" "scenario_1" {
  display_name     = "acc-security-group-with-assigned-membership-type"
  mail_enabled     = false
  mail_nickname    = "c660a1b4-5"
  security_enabled = true
  description      = "test"
  hard_delete      = true
}


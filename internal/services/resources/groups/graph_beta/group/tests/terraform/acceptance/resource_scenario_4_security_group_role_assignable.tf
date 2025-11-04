resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_4" {
  display_name          = "acc-security-group-role-assignable-${random_string.group_suffix.result}"
  mail_enabled          = false
  mail_nickname         = "accsg4${random_string.group_suffix.result}"
  security_enabled      = true
  description           = "Acceptance test - Security group with Entra role assignment capability"
  is_assignable_to_role = true
  visibility            = "Private"
}


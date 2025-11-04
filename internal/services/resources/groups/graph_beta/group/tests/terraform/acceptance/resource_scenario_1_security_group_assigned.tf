resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_1" {
  display_name     = "acc-security-group-assigned-${random_string.group_suffix.result}"
  mail_enabled     = false
  mail_nickname    = "accsg1${random_string.group_suffix.result}"
  security_enabled = true
  description      = "Acceptance test - Security group with assigned membership"
}


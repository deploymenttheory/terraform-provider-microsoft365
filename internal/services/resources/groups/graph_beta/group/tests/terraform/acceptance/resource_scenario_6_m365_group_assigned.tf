resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "scenario_6" {
  display_name          = "acc-m365-group-assigned-${random_string.group_suffix.result}"
  mail_enabled          = true
  security_enabled      = true
  group_types           = ["Unified"]
  description           = "Acceptance test - M365 group with assigned membership"
  mail_nickname         = "accm365g6${random_string.group_suffix.result}"
  is_assignable_to_role = true
  visibility            = "Private"
  hard_delete           = true
}


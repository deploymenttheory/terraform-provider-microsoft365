resource "microsoft365_graph_beta_groups_group" "scenario_6" {
  display_name        = "acc-m365-group-with-assigned-membership-type"
  mail_enabled        = true
  security_enabled    = true
  group_types         = ["Unified"]
  description         = "something"
  mail_nickname       = "acc-m365-group-with-assigned-membership-type"
  is_assignable_to_role = true
  visibility          = "Private"
}


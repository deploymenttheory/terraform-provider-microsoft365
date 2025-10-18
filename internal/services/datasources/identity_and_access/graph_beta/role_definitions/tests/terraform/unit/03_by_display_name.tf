data "microsoft365_graph_beta_identity_and_access_role_definitions" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Global Administrator"
}


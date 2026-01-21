data "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "IT Security Team"
  security_enabled = true
}

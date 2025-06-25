resource "microsoft365_graph_beta_groups_group" "minimal" {
  display_name     = "Minimal Group"
  mail_nickname    = "minimal.group"
  mail_enabled     = false
  security_enabled = true
} 
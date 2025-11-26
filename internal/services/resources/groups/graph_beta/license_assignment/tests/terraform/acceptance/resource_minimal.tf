resource "random_string" "minimal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_license_assignment" "minimal" {
  group_id = microsoft365_graph_beta_groups_group.minimal.id
  sku_id   = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE
}

resource "microsoft365_graph_beta_groups_group" "minimal" {
  display_name     = "License Assignment Test Minimal Group ${random_string.minimal_suffix.result}"
  mail_nickname    = "lictest${random_string.minimal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
}

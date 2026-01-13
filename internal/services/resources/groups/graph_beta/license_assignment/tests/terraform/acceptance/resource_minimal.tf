resource "random_string" "minimal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "minimal" {
  display_name     = "License Assignment Test Minimal Group ${random_string.minimal_suffix.result}"
  mail_nickname    = "lictest${random_string.minimal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_license_assignment" "minimal" {
  group_id = microsoft365_graph_beta_groups_group.minimal.id
  sku_id   = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  depends_on = [
    microsoft365_graph_beta_groups_group.minimal
  ]
}

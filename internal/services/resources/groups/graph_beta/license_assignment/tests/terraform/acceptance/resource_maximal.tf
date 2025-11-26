resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "dependancy" {
  display_name     = "License Assignment Test Maximal Group ${random_string.maximal_suffix.result}"
  mail_nickname    = "lictestmax${random_string.maximal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_license_assignment" "dependancy" {
  group_id = microsoft365_graph_beta_groups_group.dependancy.id
  sku_id   = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE

  depends_on = [microsoft365_graph_beta_groups_group.dependancy]
}

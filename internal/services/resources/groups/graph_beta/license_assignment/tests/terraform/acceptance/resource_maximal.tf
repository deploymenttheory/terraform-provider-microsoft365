resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "maximal" {
  display_name     = "License Assignment Test Maximal Group ${random_string.maximal_suffix.result}"
  mail_nickname    = "lictestmax${random_string.maximal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

// Wait for group creation to complete
resource "time_sleep" "wait_for_group_creation_maximal" {
  depends_on      = [microsoft365_graph_beta_groups_group.maximal]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_license_assignment" "maximal" {
  group_id = microsoft365_graph_beta_groups_group.maximal.id
  sku_id   = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  disabled_plans = [
    "c948ea65-2053-4a5a-8a62-9eaaaf11b522"  # PURVIEW_DISCOVERY
  ]

  depends_on = [
    time_sleep.wait_for_group_creation_maximal
  ]
}

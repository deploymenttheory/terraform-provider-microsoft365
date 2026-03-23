resource "random_string" "lifecycle_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "lifecycle" {
  display_name     = "License Assignment Lifecycle Test Group ${random_string.lifecycle_suffix.result}"
  mail_nickname    = "lictestlc${random_string.lifecycle_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_group_creation_lifecycle" {
  depends_on      = [microsoft365_graph_beta_groups_group.lifecycle]
  create_duration = "30s"
}

// Step 1: assign license WITH one disabled plan
resource "microsoft365_graph_beta_groups_license_assignment" "lifecycle" {
  group_id = microsoft365_graph_beta_groups_group.lifecycle.id
  sku_id   = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  disabled_plans = [
    "c948ea65-2053-4a5a-8a62-9eaaaf11b522" # PURVIEW_DISCOVERY
  ]

  depends_on = [
    time_sleep.wait_for_group_creation_lifecycle
  ]
}

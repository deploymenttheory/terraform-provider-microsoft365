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
  depends_on       = [microsoft365_graph_beta_groups_group.lifecycle]
  create_duration  = "30s"
  destroy_duration = "30s"
}

// Step 2: remove disabled_plans entirely.
// The provider must send disabledPlans=[] to the API so all plans become enabled.
// If the bug were present, the API would silently retain the plan from Step 1 and
// the next terraform plan would show unexpected drift.
resource "microsoft365_graph_beta_groups_license_assignment" "lifecycle" {
  group_id = microsoft365_graph_beta_groups_group.lifecycle.id
  sku_id   = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  depends_on = [
    time_sleep.wait_for_group_creation_lifecycle
  ]
}

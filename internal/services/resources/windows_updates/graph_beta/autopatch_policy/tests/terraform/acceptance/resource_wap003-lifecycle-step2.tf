resource "random_string" "wap003" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap003_lifecycle" {
  display_name = "acc-test-wap003-lifecycle-${random_string.wap003.result}"
  description  = "Acceptance test - lifecycle step 2: with approval rules"

  approval_rules = [
    {
      deferral_in_days = 0
      classification   = "security"
      cadence          = "monthly"
    },
    {
      deferral_in_days = 14
      classification   = "nonSecurity"
      cadence          = "monthly"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

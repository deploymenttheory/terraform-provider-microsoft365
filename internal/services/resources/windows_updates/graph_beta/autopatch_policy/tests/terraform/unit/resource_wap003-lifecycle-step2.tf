resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap003_lifecycle" {
  display_name = "WAP003: Windows Autopatch Policy Lifecycle-v1.0"
  description  = "Lifecycle step 2: with approval rules added"

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
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap002_approval_rules" {
  display_name = "WAP002: Windows Autopatch Policy With Approval Rules-v1.0"
  description  = "Policy with all approval rule types"

  approval_rules = [
    {
      deferral_in_days = 0
      classification   = "security"
      cadence          = "monthly"
    },
    {
      deferral_in_days = 7
      classification   = "nonSecurity"
      cadence          = "monthly"
    },
    {
      deferral_in_days = 0
      classification   = "security"
      cadence          = "outOfBand"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

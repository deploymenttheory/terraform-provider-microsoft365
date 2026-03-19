resource "random_string" "wap002" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap002_approval_rules" {
  display_name = "acc-test-wap002-approval-rules-${random_string.wap002.result}"
  description  = "Acceptance test - policy with approval rules"

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
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

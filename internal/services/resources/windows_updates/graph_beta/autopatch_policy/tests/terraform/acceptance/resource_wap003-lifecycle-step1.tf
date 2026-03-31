resource "random_string" "wap003" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap003_lifecycle" {
  display_name = "acc-test-wap003-lifecycle-${random_string.wap003.result}"
  description  = "Acceptance test - lifecycle step 1: no approval rules"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

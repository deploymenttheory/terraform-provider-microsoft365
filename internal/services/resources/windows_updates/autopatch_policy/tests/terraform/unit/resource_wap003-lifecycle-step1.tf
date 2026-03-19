resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap003_lifecycle" {
  display_name = "WAP003: Windows Autopatch Policy Lifecycle-v1.0"
  description  = "Lifecycle step 1: no approval rules"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

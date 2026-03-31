resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap001_minimal" {
  display_name = "WAP001: Minimal Windows Autopatch Policy-v1.0"
  description  = "Minimal policy with no approval rules"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

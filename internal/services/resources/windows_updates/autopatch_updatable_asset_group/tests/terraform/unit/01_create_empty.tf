resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "test" {
  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

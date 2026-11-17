resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "test" {
  entra_device_object_ids = []

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

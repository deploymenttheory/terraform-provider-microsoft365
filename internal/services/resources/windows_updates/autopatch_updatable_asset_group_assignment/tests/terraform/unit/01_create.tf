resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment" "test" {
  updatable_asset_group_id   = "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"
  entra_device_ids = ["aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"]
}

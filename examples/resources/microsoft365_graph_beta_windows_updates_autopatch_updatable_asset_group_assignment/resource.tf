resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment" "example" {
  updatable_asset_group_id = microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group.example.id

  entra_device_object_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
  ]
}

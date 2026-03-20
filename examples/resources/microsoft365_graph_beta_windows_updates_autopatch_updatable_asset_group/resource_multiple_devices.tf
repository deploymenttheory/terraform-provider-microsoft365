# Multiple devices — assigns multiple devices to an updatable asset group using
# known Entra device object IDs. Devices can be added or removed in-place without
# replacing the group.

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "example" {
  entra_device_object_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    "cccccccc-cccc-cccc-cccc-cccccccccccc",
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

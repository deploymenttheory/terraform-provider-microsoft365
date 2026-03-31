data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "test" {
  entra_device_object_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_active_directory_device_id
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

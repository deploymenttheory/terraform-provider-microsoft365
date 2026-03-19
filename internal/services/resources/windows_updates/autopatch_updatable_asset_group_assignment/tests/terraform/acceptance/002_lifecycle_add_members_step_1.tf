resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "test_002" {
  timeouts = {
    create = "60s"
    read   = "30s"
    delete = "60s"
  }
}

data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  filter_type = "all"

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment" "test_002" {
  updatable_asset_group_id = microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group.test_002.id

  entra_device_object_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_ad_device_id
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

# Minimal example — assigns a single device to an updatable asset group.
# Uses a data source to fetch managed devices and assigns the first device.

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "example" {
  timeouts = {
    create = "60s"
    read   = "30s"
    delete = "60s"
  }
}

data "microsoft365_graph_beta_device_management_managed_device" "devices" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment" "example" {
  updatable_asset_group_id = microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group.example.id

  entra_device_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.devices.items[0].azure_active_directory_device_id
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

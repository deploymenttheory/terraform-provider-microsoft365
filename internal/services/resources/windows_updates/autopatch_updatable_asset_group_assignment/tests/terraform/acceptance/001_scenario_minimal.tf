resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "test_001" {
  timeouts = {
    create = "60s"
    read   = "30s"
    delete = "60s"
  }
}

# Get managed devices from Intune which includes their Entra ID device IDs
data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group_assignment" "test_001" {
  updatable_asset_group_id = microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group.test_001.id

  entra_device_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_active_directory_device_id
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

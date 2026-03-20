# Example with device members — uses a managed device data source to enrol a single device
# into the updatable asset group via its Entra device object ID.

data "microsoft365_graph_beta_device_management_managed_device" "devices" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "example" {
  entra_device_object_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.devices.items[0].azure_active_directory_device_id
  ]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

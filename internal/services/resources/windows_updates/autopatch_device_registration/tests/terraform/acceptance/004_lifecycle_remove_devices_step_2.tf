
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get Entra ID devices
data "microsoft365_graph_beta_identity_and_access_device" "test_devices" {
  list_all = true
  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_004" {
  update_category = "feature"
  # Use the id field which contains the Entra ID device object ID
  entra_device_object_ids = [
    data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].id
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get managed devices from Intune which includes their Entra ID device IDs
data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_autopatch_device_registration" "test_001" {
  update_category = "feature"
  # Use the azure_ad_device_id field which contains the Entra ID device object ID
  device_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.managed_devices[0].azure_ad_device_id
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
